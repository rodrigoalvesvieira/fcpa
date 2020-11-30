package btcturk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"github.com/fatih/color"
	ex "github.com/rodrigoalvesvieira/operatordesk/exchange"
	"github.com/rodrigoalvesvieira/operatordesk/formatcurrency"
	"github.com/rodrigoalvesvieira/operatordesk/utils"
)

type OrderBookResp struct {
	Bid [][]string `json:"bids"`
	Ask [][]string `json:"asks"`
}

type BalanceResponse struct {
	ResponseData OrderBookResp `json:"data"`
	StatusCode   int           `json:"code"`
}

func AdaptOrderbook(orderbookAsks [][]string, orderbookBids [][]string) *ex.Orderbook {
	itemArrayAsks := []ex.OrderItem{}
	itemArrayBids := []ex.OrderItem{}

	for i, _ := range orderbookAsks {
		// Parsing ASK order @ i
		v := orderbookAsks[i]
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
		itemArrayAsks = append(itemArrayAsks, item)

		// Parsing BID order @ i

		v = orderbookBids[i]
		rawPrice = v[0]
		rawVol = v[1]

		orderPrice = big.NewFloat(0.0)
		orderVol = big.NewFloat(0.0)

		_, err = fmt.Sscan(rawPrice, orderPrice)
		if err != nil {
			log.Println("error scanning price:", err)
		}

		_, err = fmt.Sscan(rawVol, orderVol)
		if err != nil {
			log.Println("error scanning volume:", err)
		}

		item = ex.OrderItem{Price: orderPrice, Volume: orderVol}
		itemArrayBids = append(itemArrayBids, item)
	}

	return &ex.Orderbook{Asks: itemArrayAsks, Bids: itemArrayBids}
}

func GetOrderBook() (*ex.Orderbook, error) {
	btcTurkBaseAPIUrl := "https://api.btcturk.com/api/v2/orderbook?pairSymbol=BTC_TRY&limit=300"

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
				rawBook := AdaptOrderbook(jsonData.ResponseData.Ask, jsonData.ResponseData.Bid)
				return rawBook, nil
			} else {

			}
		} else {

		}
	}

	return &ex.Orderbook{}, nil
}

// Takes targetVolume ()
func GetCryptoPayout(targetVolume *big.Float) *big.Float {
	takerFee := big.NewFloat(0.002)
	zero := big.NewFloat(0.0)
	aggregate := big.NewFloat(0.0)

	// Total fiat
	aggregatePayout := big.NewFloat(0.0)

	raw_book, _ := GetOrderBook()
	book := raw_book.Asks

	i := 0
	for i = 0; i < len(book); i++ {
		currentOrder := book[i]
		fmt.Println(currentOrder)

		loc := big.NewFloat(0.0).Mul(currentOrder.Price, currentOrder.Volume)

		if aggregate.Cmp(targetVolume) < 0 {
			aheadAggregate := zero.Add(aggregate, loc)

			if aheadAggregate.Cmp(targetVolume) > 0 {
				rem := big.NewFloat(0.0).Copy(targetVolume)
				rem.Sub(rem, aggregate)

				aggregate.Add(aggregate, rem)

				// We need less than the volume of the current entry to satisfy our target
				// Calculating remaining volume to satisfy target
				remVol := big.NewFloat(0.0).Copy(rem) // 0.4
				remVol.Mul(remVol, currentOrder.Price)

				payout := utils.SubtractFee(remVol, takerFee)
				aggregatePayout.Add(aggregatePayout, payout)
			} else {
				aggregate.Add(aggregate, loc)

				payout := utils.SubtractFee(currentOrder.Price, takerFee)
				aggregatePayout.Add(aggregatePayout, payout)
			}
		} else {
			break
		}
	}

	return aggregatePayout
}

func GetFiatPayout(targetVolume *big.Float) *big.Float {
	printer := color.New(color.FgYellow, color.Bold)
	printer.Println("GETTING ORDER BOOK - BTCTURK")

	// Assumes 1-5M TRY 30-day volume
	// https://pro.btcturk.com/en/help/fees
	takerFee := big.NewFloat(0.0016)
	aggregate := big.NewFloat(0.0)
	aggregateVol := big.NewFloat(0.0)

	raw_book, _ := GetOrderBook()
	book := raw_book.Bids
	i := 0

	for i = 0; i < len(book); i++ {
		currentOrder := book[i]

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

	fmt.Println("BTCTURK TOTAL TRAVERSED  " + strconv.Itoa(i+1))
	fmt.Println("BTCTURK TOTAL FIAT       " + formatcurrency.PrintTRY(aggregate))
	fmt.Println("BTCTURK TOTAL VOLUME     " + aggregateVol.String())

	return aggregate
}

func GetFiatPayoutForBook(book []ex.OrderItem, targetVolume *big.Float) (*big.Float, *big.Float, int, *big.Float, *big.Float, int) {
	// Assumes 1-5M TRY 30-day volume
	// https://pro.btcturk.com/en/help/fees
	takerFee := big.NewFloat(0.0016)
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
			} else {
				fiatPayout := big.NewFloat(0.0).Copy(currentOrder.Price)
				fiatPayout.Mul(fiatPayout, currentOrder.Volume)

				realFiat := utils.SubtractFee(fiatPayout, takerFee)
				aggregate.Add(aggregate, realFiat)

				aggregateVol.Add(aggregateVol, currentOrder.Volume)
			}
		}
	}

	return aggregateVol, aggregate, orderCount, totalCryptoVolume, totalFiatVolume, i
}
