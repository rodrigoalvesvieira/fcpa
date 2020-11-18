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

type TurkOrderBookResp struct {
	Bid [][]string `json:"bids"`
	Ask [][]string `json:"asks"`
}

type BalanceResponse struct {
	ResponseData TurkOrderBookResp   `json:"data"`
	StatusCode   int                 `json:"code"`
}

func getBTCTurkCostForVolume(volume *big.Float) {
	printer := color.New(color.FgYellow, color.Bold)
	printer.Println("GETTING ORDER BOOK - BTC TURX")

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
				fmt.Println("Parsing successful")

				// aggregate := big.NewFloat(0.0)

				for i := 0; i < len(jsonData.ResponseData.Bid); i++ {
					currentOrder := jsonData.ResponseData.Bid[i]

					// orderPrice := currentOrder[0]
					rawVol := currentOrder[1]
					orderVol := big.NewFloat(0.0)

					_, err := fmt.Sscan(rawVol, orderVol)
					if err != nil {
						log.Println("error scanning value:", err)
					} else {


						fmt.Println(orderVol)
					}

					fmt.Println(orderVol)
				}

			}
		}
	}
}

func printOrder(i int, volume *big.Float, price *big.Float) {
	fmt.Println("Order #" + strconv.Itoa(i + 1) + " Volume: " + volume.String() + " Price: " + price.String())
}

func getKrakenCostForVolume() {
	printer := color.New(color.FgGreen, color.Bold)
	printer.Println("GETTING ORDER BOOK - KRAKEN")

	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())

	krakenUSDBalance, err := api.Depth(krakenapi.XXBTZUSD, 200)
	if err != nil {
		log.Fatal(err)
	}

	bookA := adaptOrderbook(krakenUSDBalance.Asks)

	zero := big.NewFloat(0.0)
	total := big.NewFloat(100000.0)
	aggregate := big.NewFloat(0.0)
	aggregateVol := big.NewFloat(0.0)

	// ordersSum := 0

	for i := 0; i < len(bookA); i++ {
		orVol := bookA[i].Volume
		orPrice := bookA[i].Price

		printOrder(i, orVol, orPrice)

		loc := big.NewFloat(0.0)
		loc = loc.Mul(bookA[i].Price, orVol)

		if aggregate.Cmp(total) < 0 {
			aheadAggregate := zero.Add(loc, aggregate)

			if (aheadAggregate.Cmp(total) > 0) {
				rem := zero

				rem.Sub(total, aggregate)

				// fmt.Println("FALTA " + rem.String())

				fmt.Println(loc)
			} else {
				aggregateVol.Add(aggregateVol, orVol)
				aggregate.Add(aggregate, loc)
				fmt.Print("=>")
				fmt.Println(i)
			}

			//fmt.Println("AGG       " + aggregate.String())
			//fmt.Println("AHEAD AGG " + aheadAggregate.String())

		} else {
			break
		}
	}

	fmt.Println("--------------")


	getBTCTurkCostForVolume(aggregateVol)

	// fmt.Println("KRAKEN TOTAL TRAVERSED " + strconv.Itoa(ordersSum))

	fmt.Println("Total price Kraken " + aggregate.String())
	fmt.Println("Total volume Kraken " + aggregateVol.String())
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	green := color.New(color.FgGreen, color.Bold)
	green.Println("BRENT's FIAT-CRYPTO PARITY ARBITRAGE DESK")
	green.Println("-------------------------------------------\n\n")

	getKrakenCostForVolume()

	magenta := color.New(color.FgMagenta, color.Bold).Add(color.Underline)
	magenta.Println("BRENT's FIAT-CRYPTO PARITY ARBITRAGE DESK")
}
