package main

import (
	"log"
	"fmt"
	"github.com/fatih/color"
	"github.com/beldur/kraken-go-api-client"
	"github.com/rodrigoalvesvieira/operatordesk/envars"
	"github.com/joho/godotenv"
	"math/big"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type Coin string
type Fiat string

type OrderItem struct {
	Price  *big.Float // unit price
	Volume *big.Float // volume at this price
}

type Orderbook struct {
	CoinName Coin
	FiatName Fiat
	Asks     []OrderItem
	Bids     []OrderItem
}

func adaptOrderbook(orderbook []krakenapi.OrderBookItem) []OrderItem {
	itemArray := []OrderItem{}
	for _, v := range orderbook {
		item := OrderItem{Price: big.NewFloat(v.Price), Volume: big.NewFloat(v.Amount)}
		itemArray = append(itemArray, item)
	}
	return itemArray
}

func adaptOrderbookTurk(orderbook [][]string) []OrderItem {
	itemArray := []OrderItem{}

	for _, v := range orderbook {

		rawPrice := v[0]
		rawVol := v[1]

		orderPrice := big.NewFloat(0.0)
		orderVol := big.NewFloat(0.0)

		_, err := fmt.Sscan(rawPrice, orderPrice)
		if err != nil {
			log.Println("error scanning price:", err)
		}

		_, err = fmt.Sscan(rawVol, orderVol)
		if err != nil {
			log.Println("error scanning volume:", err)
		}

		item := OrderItem{Price: orderPrice, Volume: orderVol}
		itemArray = append(itemArray, item)
	}

	return itemArray
}

type TurkOrderBookResp struct {
	Bid [][]string `json:"bids"`
	Ask [][]string `json:"asks"`
}

type BalanceResponse struct {
	ResponseData TurkOrderBookResp   `json:"data"`
	StatusCode   int                 `json:"code"`
}

func getBTCTurkCostForVolume(targetVolume *big.Float) {
	printer := color.New(color.FgYellow, color.Bold)
	printer.Println("GETTING ORDER BOOK - BTCTURK")

	btcTurkBaseAPIUrl := "https://api.btcturk.com/api/v2/orderbook?pairSymbol=BTC_TRY"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", btcTurkBaseAPIUrl, nil)
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			var jsonData *BalanceResponse
			errJSON := json.Unmarshal(body, &jsonData)

			if errJSON == nil {
				aggregate := big.NewFloat(0.0)
				aggregateVol := big.NewFloat(0.0)

				bookBids := adaptOrderbookTurk(jsonData.ResponseData.Bid)
				i := 0

				for i = 0; i < len(bookBids); i++ {
					currentOrder := bookBids[i]

					if (aggregateVol.Cmp(targetVolume) < 0) {
						aheadAggregate := big.NewFloat(0.0).Add(aggregateVol, currentOrder.Volume)

						if (aheadAggregate.Cmp(targetVolume) > 0) {
							rem := big.NewFloat(0.0).Copy(targetVolume)
							rem.Sub(rem, aggregateVol)

							// We need less than the volume of the current entry to satisfy our target
							// Calculating remaining volume to satisfy target
							remFiat := big.NewFloat(0.0).Copy(rem)
							remFiat.Mul(remFiat, currentOrder.Volume)

							aggregate.Add(aggregate, remFiat)
							aggregateVol.Add(aggregateVol, rem)
						} else {
							fiatPayout := big.NewFloat(0.0).Copy(currentOrder.Price)
							fiatPayout.Mul(fiatPayout, currentOrder.Volume)

							aggregate.Add(aggregate, fiatPayout)
							aggregateVol.Add(aggregateVol, currentOrder.Volume)
						}
					} else {
						break
					}
				}

				fmt.Println("BTCTURK TOTAL TRAVERSED  " + strconv.Itoa(i + 1))
				fmt.Println("BTCTURK TOTAL FIAT       " + aggregate.String())
				fmt.Println("BTCTURK TOTAL VOLUME     " + aggregateVol.String())
			}
		}
	}
}

func printOrder(i int, volume *big.Float, price *big.Float) {
	fmt.Println("Order #" + strconv.Itoa(i + 1) + " Volume: " + volume.String() + " Price: " + price.String())
}

func getKrakenCostForVolume(total *big.Float) {
	printer := color.New(color.FgGreen, color.Bold)
	printer.Println("GETTING ORDER BOOK - KRAKEN")

	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())

	krakenUSDBalance, err := api.Depth(krakenapi.XXBTZUSD, 200)
	if err != nil {
		log.Fatal(err)
	}

	bookA := adaptOrderbook(krakenUSDBalance.Asks)

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

				aggregateVol.Add(aggregateVol, remVol)
			} else {
				aggregate.Add(aggregate, loc)
				aggregateVol.Add(aggregateVol, orVol)
			}
		} else {
			break
		}
	}

	fmt.Println("--------------")

	getBTCTurkCostForVolume(aggregateVol)

	fmt.Println("KRAKEN TOTAL TRAVERSED " + strconv.Itoa(i + 1))
	fmt.Println("KRAKEN TOTAL FIAT      " + aggregate.String())
	fmt.Println("KRAKEN TOTAL VOLUME    " + aggregateVol.String())
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	green := color.New(color.FgGreen, color.Bold)
	green.Println("BRENT's FIAT-CRYPTO PARITY ARBITRAGE DESK")
	green.Println("-------------------------------------------\n")

	target := big.NewFloat(100000.0)
	getKrakenCostForVolume(target)

	magenta := color.New(color.FgMagenta, color.Bold).Add(color.Underline)
	magenta.Println("-----------------------------------------")
}
