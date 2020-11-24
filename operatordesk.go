package main

import (
	"log"
	"fmt"
	"github.com/fatih/color"
	"github.com/beldur/kraken-go-api-client"
	"github.com/rodrigoalvesvieira/operatordesk/envars"
	"github.com/joho/godotenv"
	"math/big"
	"strconv"
	"github.com/rodrigoalvesvieira/operatordesk/formatcurrency"
	"github.com/rodrigoalvesvieira/operatordesk/btcturk"
	"github.com/rodrigoalvesvieira/operatordesk/utils"
	ex "github.com/rodrigoalvesvieira/operatordesk/exchange"
)

func adaptOrderbook(orderbook []krakenapi.OrderBookItem) []ex.OrderItem {
	itemArray := []ex.OrderItem{}
	for _, v := range orderbook {
		item := ex.OrderItem{Price: big.NewFloat(v.Price), Volume: big.NewFloat(v.Amount)}
		itemArray = append(itemArray, item)
	}
	return itemArray
}

func printOrder(i int, volume *big.Float, price *big.Float) {
	fmt.Println("Order #" + strconv.Itoa(i + 1) + " Volume: " + volume.String() + " Price: " + price.String())
}

// Takes a target amount of crypto and returns the amount of fiat currency that you can sell it for
func getKrakenFiatPayout(targetVolume *big.Float) *big.Float {
	printer := color.New(color.FgBlue, color.Bold)
	printer.Println("GETTING ORDER BOOK TO SELL CRYPTO - KRAKEN")

	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())

	krakenUSDBalance, err := api.Depth(krakenapi.XXBTZUSD, 200)
	if err != nil {
		log.Fatal(err)
	}

	book := adaptOrderbook(krakenUSDBalance.Asks)

	// Assumes 30-day $250,001 - $500,000 volume
	// https://www.kraken.com/en-us/features/fee-schedule
	takerFee := big.NewFloat(0.002)
	aggregate := big.NewFloat(0.0)
	aggregateVol := big.NewFloat(0.0)
	i := 0

	for i = 0; i < len(book); i++ {
		currentOrder := book[i]
		printOrder(i, currentOrder.Volume, currentOrder.Price)

		if (aggregateVol.Cmp(targetVolume) < 0) {
			aheadAggregate := big.NewFloat(0.0).Add(aggregateVol, currentOrder.Volume)

			if (aheadAggregate.Cmp(targetVolume) > 0) {
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

	fmt.Println("KRAKEN TOTAL TRAVERSED   " + strconv.Itoa(i + 1))
	fmt.Println("KRAKEN TOTAL FIAT        " + formatcurrency.PrintUSD(aggregate))
	fmt.Println("KRAKEN TOTAL VOLUME      " + aggregateVol.String())

	return aggregate
}

// Takes a target amount of fiat currency and returns the amount of crypto that it can buy
func getKrakenCryptoPayout(total *big.Float) *big.Float {
	printer := color.New(color.FgGreen, color.Bold)
	printer.Println("GETTING ORDER BOOK - KRAKEN")

	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())

	krakenUSDBalance, err := api.Depth(krakenapi.XXBTZUSD, 200)
	if err != nil {
		log.Fatal(err)
	}

	bookA := adaptOrderbook(krakenUSDBalance.Asks)

	// Assumes 30-day $250,001 - $500,000 volume
	// https://www.kraken.com/en-us/features/fee-schedule
	takerFee := big.NewFloat(0.002)
	zero := big.NewFloat(0.0)
	aggregate := big.NewFloat(0.0)
	aggregateVol := big.NewFloat(0.0)

	i := 0

	for i = 0; i < len(bookA); i++ {
		orVol := bookA[i].Volume
		loc := big.NewFloat(0.0)
		loc = loc.Mul(bookA[i].Price, orVol)

		if aggregate.Cmp(total) < 0 {
			aheadAggregate := zero.Add(loc, aggregate)

			if (aheadAggregate.Cmp(total) > 0) {
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

				break
			} else {
				aggregate.Add(aggregate, loc)

				realVol := utils.SubtractFee(orVol, takerFee)
				aggregateVol.Add(aggregateVol, realVol)
			}
		} else {
			break
		}
	}

	fmt.Println("KRAKEN TOTAL TRAVERSED   " + strconv.Itoa(i + 1))
	fmt.Println("KRAKEN TOTAL FIAT        " + formatcurrency.PrintUSD(aggregate))
	fmt.Println("KRAKEN TOTAL VOLUME      " + aggregateVol.String())
	return aggregateVol
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	green := color.New(color.FgGreen, color.Bold)
	green.Println("BRENT's FIAT-CRYPTO PARITY ARBITRAGE DESK")
	green.Println("-------------------------------------------\n")

	usdIn := big.NewFloat(40000.0)
	cryptoOut := getKrakenCryptoPayout(usdIn)

	btcturk.GetBTCTurkCostForVolume(cryptoOut)

	magenta := color.New(color.FgMagenta, color.Bold).Add(color.Underline)
	//magenta.Println("-----------------------------------------")

	cryptoIn := big.NewFloat(2.0)
	usdOut := getKrakenFiatPayout(cryptoIn)
	magenta.Println("KRAKEN Input BTC: " + cryptoIn.String() + " - " + usdOut.String())

	red := color.New(color.FgRed).Add(color.Bold)
	red.Println("-----------------------------------------")
}
