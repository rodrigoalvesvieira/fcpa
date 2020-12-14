package kraken

import (
	"log"
	"math/big"
	krakenapi "github.com/beldur/kraken-go-api-client"
	"github.com/fatih/color"
	ex "operatordesk/exchange"
	"operatordesk/utils"
	"operatordesk/envars"
)

func adaptOrderbook(orderbook []krakenapi.OrderBookItem) []ex.OrderItem {
	itemArray := []ex.OrderItem{}
	for _, v := range orderbook {
		item := ex.OrderItem{Price: big.NewFloat(v.Price), Volume: big.NewFloat(v.Amount)}
		itemArray = append(itemArray, item)
	}
	return itemArray
}

// Takes a target amount of crypto and returns the amount of fiat currency that you can sell it for
func GetFiatPayoutForBook(book []ex.OrderItem, targetVolume *big.Float) (*big.Float, *big.Float, int, *big.Float, *big.Float, int) {
	printer := color.New(color.FgBlue, color.Bold)
	printer.Println("GETTING ORDER BOOK TO SELL CRYPTO - KRAKEN")

	// Assumes 30-day $250,001 - $500,000 volume
	// https://www.kraken.com/en-us/features/fee-schedule
	takerFee := big.NewFloat(0.002)
	aggregate := big.NewFloat(0.0)
	aggregateVol := big.NewFloat(0.0)

	totalCryptoVolume := big.NewFloat(0.0)
	totalFiatVolume := big.NewFloat(0.0)

	orderCount := 0
	i := 0

	for i = 0; i < len(book); i++ {
		currentOrder := book[i]

		totalCryptoVolume.Add(totalCryptoVolume, currentOrder.Volume)
		totalFiatVolume.Add(totalFiatVolume, utils.GetOrderCost(currentOrder))

		if aggregateVol.Cmp(targetVolume) < 0 {
			orderCount += 1

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

	return aggregateVol, aggregate, orderCount, totalCryptoVolume, totalFiatVolume, i
}

func CalculateCryptoPayout(book []ex.OrderItem, total *big.Float) (*big.Float, *big.Float, int, *big.Float, *big.Float, int) {
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
		totalFiatVolume.Add(totalFiatVolume, utils.GetOrderCost(currentOrder))

		orVol := currentOrder.Volume
		loc := big.NewFloat(0.0).Mul(currentOrder.Price, orVol)

		if aggregate.Cmp(total) < 0 {
			orderCount += 1
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
	return CalculateCryptoPayout(book, total)
}

func GetOrderBook() (*ex.Orderbook, error) {
	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())

	krakenBook, err := api.Depth(krakenapi.XXBTZUSD, 200)
	if err != nil {
		log.Fatal(err)
	}

	itemArrayAsks := []ex.OrderItem{}
	itemArrayBids := []ex.OrderItem{}

	for i, _ := range krakenBook.Asks {
		// Parsing ASK order @ i
		v := krakenBook.Asks[i]
		item := ex.OrderItem{Price: big.NewFloat(v.Price), Volume: big.NewFloat(v.Amount)}

		itemArrayAsks = append(itemArrayAsks, item)

		// Parsing BID order @ i
		v = krakenBook.Bids[i]
		item = ex.OrderItem{Price: big.NewFloat(v.Price), Volume: big.NewFloat(v.Amount)}

		itemArrayBids = append(itemArrayBids, item)
	}

	return &ex.Orderbook{Asks: itemArrayAsks, Bids: itemArrayBids}, nil
}

// Takes a target amount of fiat currency and returns the amount of crypto that it can buy
func GetCryptoPayout(total *big.Float) (*big.Float, *big.Float, int, *big.Float, *big.Float, int) {
	printer := color.New(color.FgGreen, color.Bold)
	printer.Println("GETTING ORDER BOOK - KRAKEN")

	raw_book, _ := GetOrderBook()
	return CalculateCryptoPayout(raw_book.Asks, total)
}
