package btcturk

import (
	"log"
	"fmt"
	"github.com/fatih/color"
	"math/big"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"github.com/rodrigoalvesvieira/operatordesk/formatcurrency"
	"github.com/rodrigoalvesvieira/operatordesk/utils"
	ex "github.com/rodrigoalvesvieira/operatordesk/exchange"
)

type OrderBookResp struct {
	Bid [][]string `json:"bids"`
	Ask [][]string `json:"asks"`
}

type BalanceResponse struct {
	ResponseData OrderBookResp  `json:"data"`
	StatusCode   int            `json:"code"`
}

func adaptOrderbook(orderbook [][]string) []ex.OrderItem {
	itemArray := []ex.OrderItem{}

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

		item := ex.OrderItem{Price: orderPrice, Volume: orderVol}
		itemArray = append(itemArray, item)
	}

	return itemArray
}

func GetBTCTurkCostForVolume(targetVolume *big.Float) {
	printer := color.New(color.FgYellow, color.Bold)
	printer.Println("GETTING ORDER BOOK - BTCTURK")

	btcTurkBaseAPIUrl := "https://api.btcturk.com/api/v2/orderbook?pairSymbol=BTC_TRY&limit=300"

	// Assumes 1-5M TRY 30-day volume
	// https://pro.btcturk.com/en/help/fees
	takerFee := big.NewFloat(0.0016)

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

				bookBids := adaptOrderbook(jsonData.ResponseData.Bid)
				fmt.Println("BTCTURK BOOK LENGTH " + strconv.Itoa(len(bookBids)));
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

							realFiat := utils.SubtractFee(remFiat, takerFee)
							aggregate.Add(aggregate, realFiat)

							aggregateVol.Add(aggregateVol, rem)
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

				fmt.Println("BTCTURK TOTAL TRAVERSED  " + strconv.Itoa(i + 1))
				fmt.Println("BTCTURK TOTAL FIAT       " + formatcurrency.PrintTRY(aggregate))
				fmt.Println("BTCTURK TOTAL VOLUME     " + aggregateVol.String())
			}
		}
	}
}
