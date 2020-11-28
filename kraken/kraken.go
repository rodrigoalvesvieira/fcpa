package kraken

import (
	"fmt"
	"log"
	"math/big"
	"strconv"

	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fatih/color"
	//"github.com/rodrigoalvesvieira/operatordesk/btcturk"
	"github.com/rodrigoalvesvieira/operatordesk/envars"
	ex "github.com/rodrigoalvesvieira/operatordesk/exchange"
	"github.com/rodrigoalvesvieira/operatordesk/formatcurrency"
	"github.com/rodrigoalvesvieira/operatordesk/utils"
)

func AdaptOrderbook(orderbook []krakenapi.OrderBookItem) []ex.OrderItem {
	itemArray := []ex.OrderItem{}
	for _, v := range orderbook {
		item := ex.OrderItem{Price: big.NewFloat(v.Price), Volume: big.NewFloat(v.Amount)}
		itemArray = append(itemArray, item)
	}
	return itemArray
}

func printOrder(i int, volume *big.Float, price *big.Float) {
	fmt.Println("Order #" + strconv.Itoa(i+1) + " Volume: " + volume.String() + " Price: " + price.String())
}

// Takes a target amount of crypto and returns the amount of fiat currency that you can sell it for
func GetKrakenFiatPayout(targetVolume *big.Float) *big.Float {
	printer := color.New(color.FgBlue, color.Bold)
	printer.Println("GETTING ORDER BOOK TO SELL CRYPTO - KRAKEN")

	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())

	krakenUSDBalance, err := api.Depth(krakenapi.XXBTZUSD, 200)
	if err != nil {
		log.Fatal(err)
	}

	book := AdaptOrderbook(krakenUSDBalance.Asks)

	// Assumes 30-day $250,001 - $500,000 volume
	// https://www.kraken.com/en-us/features/fee-schedule
	takerFee := big.NewFloat(0.002)
	aggregate := big.NewFloat(0.0)
	aggregateVol := big.NewFloat(0.0)
	i := 0

	for i = 0; i < len(book); i++ {
		currentOrder := book[i]
		// printOrder(i, currentOrder.Volume, currentOrder.Price)

		if aggregateVol.Cmp(targetVolume) < 0 {
			aheadAggregate := big.NewFloat(0.0).Add(aggregateVol, currentOrder.Volume)

			if aheadAggregate.Cmp(targetVolume) > 0 {
				rem := big.NewFloat(0.0).Copy(targetVolume)
				rem.Sub(rem, aggregateVol)

				// We need less than the volume of the current entry to satisfy our target
				// Calculating remaining volume to satisfy target
				remFiat := big.NewFloat(0.0).Copy(rem)
				remFiat.Mul(remFiat, currentOrder.Volume)

				realFiat := utils.SubtractFee(remFiat, takerFee)
				aggregate.Add(aggregate, realFiat)
				aggregateVol.Add(aggregateVol, rem)

				break
			} else {
				fiatPayout := big.NewFloat(0.0).Copy(currentOrder.Price)
				fiatPayout.Mul(fiatPayout, currentOrder.Volume)

				realFiat := utils.SubtractFee(fiatPayout, takerFee)
				aggregate.Add(aggregate, realFiat)

				aggregateVol.Add(aggregateVol, currentOrder.Volume)
			}
		} else {
			break
		}
	}

	fmt.Println("KRAKEN TOTAL TRAVERSED   " + strconv.Itoa(i+1))
	fmt.Println("KRAKEN TOTAL FIAT        " + formatcurrency.PrintUSD(aggregate))
	fmt.Println("KRAKEN TOTAL VOLUME      " + aggregateVol.String())

	return aggregate
}

func getOrderCost(order ex.OrderItem) *big.Float {
	base := big.NewFloat(0.0).Copy(order.Price)
	return base.Mul(base, order.Volume)
}

func calculateCryptoPayout(book []ex.OrderItem, total *big.Float) (*big.Float, *big.Float, int, *big.Float, *big.Float, int) {
	// Assumes 30-day $250,001 - $500,000 volume
	// https://www.kraken.com/en-us/features/fee-schedule
	takerFee := big.NewFloat(0.002)
	zero := big.NewFloat(0.0)
	aggregate := big.NewFloat(0.0)
	aggregateVol := big.NewFloat(0.0)

	totalCryptoVolume := big.NewFloat(0.0)
	totalFiatVolume := big.NewFloat(0.0)

	orderCount := 0
	i := 0

	for i = 0; i < len(book); i++ {
		currentOrder := book[i]

		totalCryptoVolume.Add(totalCryptoVolume, currentOrder.Volume)
		totalFiatVolume.Add(totalFiatVolume, getOrderCost(currentOrder))

		orVol := currentOrder.Volume
		loc := big.NewFloat(0.0).Mul(currentOrder.Price, orVol)

		if aggregate.Cmp(total) < 0 {
			orderCount += 1
			// printOrder(i, currentOrder.Volume, currentOrder.Price)

			aheadAggregate := zero.Add(loc, aggregate)

			if aheadAggregate.Cmp(total) > 0 {
				rem := big.NewFloat(0.0).Copy(total)
				rem.Sub(rem, aggregate)

				aggregate.Add(aggregate, rem)

				// We need less than the volume of the current entry to satisfy our target
				// Calculating remaining volume to satisfy target
				remVol := big.NewFloat(0.0).Copy(rem)
				remVol.Mul(remVol, orVol)
				remVol.Quo(remVol, loc)

				realVol := utils.SubtractFee(remVol, takerFee)
				aggregateVol.Add(aggregateVol, realVol)
			} else {
				aggregate.Add(aggregate, loc)

				realVol := utils.SubtractFee(orVol, takerFee)
				aggregateVol.Add(aggregateVol, realVol)
			}
		}
	}

	return aggregateVol, aggregate, orderCount, totalCryptoVolume, totalFiatVolume, i

}

func GetCryptoPayoutForBook(book []ex.OrderItem, total *big.Float) (*big.Float, *big.Float, int, *big.Float, *big.Float, int) {
	return calculateCryptoPayout(book, total)
}

func GetKrakenBook() *krakenapi.OrderBook {
	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())

	krakenBook, err := api.Depth(krakenapi.XXBTZUSD, 200)
	if err != nil {
		log.Fatal(err)
	}

	return krakenBook
}

// Takes a target amount of fiat currency and returns the amount of crypto that it can buy
func GetCryptoPayout(total *big.Float) (*big.Float, *big.Float, int, *big.Float, *big.Float, int) {
	printer := color.New(color.FgGreen, color.Bold)
	printer.Println("GETTING ORDER BOOK - KRAKEN")

	book := AdaptOrderbook(GetKrakenBook().Asks)
	return calculateCryptoPayout(book, total)
}
