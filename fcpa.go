package main

import (
	"bytes"
	clock "fcpa/clockwork"
	ex "fcpa/exchange"
	"fcpa/formatcurrency"
	"fcpa/markets/btcturk"
	"fcpa/markets/kraken"
	ts "fcpa/timeseries"
	"fcpa/utils"
	"fmt"
	"math/big"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

// GenericFunction Alias for a function that takes no args and returns an error
type GenericFunction func() error

// RunAsyncAllowErrors allows one or more function to run asynchronously
func RunAsyncAllowErrors(functions ...GenericFunction) []error {
	// Channel for communicating when execution has finished.
	finished := make(chan struct{}, 1)

	// Create a wait group for the functions to be executed.
	// Note: wg.Add sets the number of goroutines to wait on.
	//       wg.Done decrements the number of goroutines left to wait on by one.
	var wg sync.WaitGroup
	wg.Add(len(functions))

	// Wait group in go routine. Will not close until all functions have completed
	go func() {
		wg.Wait()
		close(finished)
	}()

	// Spawn a goroutine for each function
	errors := make([]error, len(functions))
	for j := range functions {
		go func(j int) {
			defer wg.Done()
			// Clause for handling panic errors
			defer func() {
				if r := recover(); r != nil {
					// Skip 4 stack frames:
					// 1) debug.Stack()
					// 2) formatStack()
					// 3) this anonymous func
					// 4) runtime/panic
					err := fmt.Errorf(
						"panic in async function: %v\n%s",
						r, formatStack(4))
					errors[j] = err
				}
			}()
			// Execute function and log error at appropriate
			// position if applicable.
			if err := functions[j](); err != nil {
				errors[j] = err
				return
			}
		}(j)
	}

	// The code below this will not be reached until the wait group has finished
	<-finished
	return errors
}

// Return formatted stack trace, skipping "skip" leading stack frames
func formatStack(skip int) string {
	lines := bytes.Split(bytes.TrimSpace(debug.Stack()), []byte("\n"))
	formatted := bytes.Join(lines[1+2*skip:], []byte("\n"))
	return string(formatted)
}

func makeRequest(fn func() (*ex.Orderbook, error), book *ex.Orderbook) GenericFunction {
	return func() error {
		res, err := fn()

		if err != nil {
			return err
		}

		nsec := time.Now().UnixNano()

		// Pln("RESPONSE TIME (THREADS): " + strconv.FormatInt(nsec, 10))

		res.Time = nsec
		*book = *res
		return nil
	}
}

func processArbitrageBuyBTCTurkSellKraken(books []ex.Orderbook, tiers ts.Result) {
	utils.P("green", "INPUT TRY OUTPUT USD")

	krakenBook := books[0]
	btcTurkBook := books[1]

	for i := 0; i < len(tiers.Values); i++ {
		tier := tiers.Values[i]
		id := tiers.Ids[i]

		utils.P("red", "QUERYING "+formatcurrency.PrintTRY(tier))

		cryptoOut, fiatIn, buyCount, buyExCryptoVol, buyExFiatVol, buyExBookCount := btcturk.GetCryptoPayoutForBook(btcTurkBook.Asks, tier)

		utils.P("magenta", "BTCTURK TOTAL TRAVERSED    "+strconv.Itoa(buyCount))
		utils.P("magenta", "BTCTURK TOTAL FIAT         "+formatcurrency.PrintTRY(fiatIn))
		utils.P("magenta", "BTCTURK TOTAL VOLUME       "+cryptoOut.String())

		cryptoIn, fiatOut, sellCount, sellExCryptoVol, sellExFiatVol, sellExBookLength := kraken.GetFiatPayoutForBook(krakenBook.Bids, cryptoOut)
		utils.P("yellow", "KRAKEN TOTAL TRAVERSED   "+strconv.Itoa(sellCount))
		utils.P("yellow", "KRAKEN TOTAL FIAT        "+formatcurrency.PrintUSD(fiatOut))
		utils.P("yellow", "KRAKEN TOTAL VOLUME      "+cryptoIn.String())

		exchangerate := big.NewFloat(0.0).Copy(fiatOut)
		exchangerate.Quo(exchangerate, fiatIn)

		utils.P("higreen", "EXCHANGE RATE "+formatcurrency.PrintUSD(exchangerate))

		entry := ts.PriceQuery{
			Coin:                   1,
			Tier:                   id,
			Exchange_buy:           2,
			Exchange_sell:          1,
			Exchange_buy_coin:      4,
			Exchange_sell_coin:     3,
			Buy_query_time:         btcTurkBook.Time,
			Sell_query_time:        krakenBook.Time,
			Buy_input:              tier,
			Buy_count:              buyCount,
			Buy_vol:                cryptoOut,
			Sell_input:             cryptoOut,
			Sell_count:             sellCount,
			Sell_vol:               fiatOut,
			Buy_total_book_length:  buyExBookCount,
			Buy_total_book_crypto:  buyExCryptoVol,
			Buy_total_book_fiat:    buyExFiatVol,
			Sell_total_book_length: sellExBookLength,
			Sell_total_book_crypto: sellExCryptoVol,
			Sell_total_book_fiat:   sellExFiatVol,
			Exchange_rate:          exchangerate,
		}

		ts.InsertPriceQuery(&entry)
	}
}

func processArbitrageBuyKrakenSellBTCTurk(books []ex.Orderbook, tiers ts.Result) {
	utils.P("green", "INPUT USD OUTPUT TRY")

	krakenBook := books[0]
	btcTurkBook := books[1]

	for i := 0; i < len(tiers.Values); i++ {
		tier := tiers.Values[i]
		id := tiers.Ids[i]

		utils.P("red", "QUERYING "+formatcurrency.PrintUSD(tier))

		cryptoOut, fiatIn, buyCount, buyExCryptoVol, buyExFiatVol, buyExBookCount := kraken.GetCryptoPayoutForBook(krakenBook.Asks, tier)
		utils.P("magenta", "KRAKEN TOTAL TRAVERSED    "+strconv.Itoa(buyCount))
		utils.P("magenta", "KRAKEN TOTAL FIAT         "+formatcurrency.PrintUSD(fiatIn))
		utils.P("magenta", "KRAKEN TOTAL VOLUME       "+cryptoOut.String())

		cryptoIn, fiatOut, sellCount, sellExCryptoVol, sellExFiatVol, sellExBookLength := btcturk.GetFiatPayoutForBook(btcTurkBook.Bids, cryptoOut)
		utils.P("yellow", "QUERYING "+formatcurrency.PrintUSD(tier)+" "+cryptoOut.String())
		utils.P("yellow", "BTCTURK TOTAL TRAVERSED   "+strconv.Itoa(sellCount))
		utils.P("yellow", "BTCTURK TOTAL FIAT        "+formatcurrency.PrintTRY(fiatOut))
		utils.P("yellow", "BTCTURK TOTAL VOLUME      "+cryptoIn.String())

		exchangerate := big.NewFloat(0.0).Copy(fiatOut)
		exchangerate.Quo(exchangerate, fiatIn)

		utils.P("higreen", "EXCHANGE RATE "+formatcurrency.PrintTRY(exchangerate))

		utils.P("cyan", "KRAKEN TOTAL BOOK VOLUME  "+buyExCryptoVol.String())
		utils.P("cyan", "KRAKEN TOTAL BOOK LENGTH  "+strconv.Itoa(buyExBookCount))
		utils.P("cyan", "KRAKEN TOTAL BOOK FIAT    "+formatcurrency.PrintUSD(buyExFiatVol))

		utils.P("cyan", "BTCTURK TOTAL BOOK VOLUME "+sellExCryptoVol.String())
		utils.P("cyan", "BTCTURK TOTAL BOOK LENGTH "+strconv.Itoa(sellExBookLength))
		utils.P("cyan", "BTCTURK TOTAL BOOK FIAT   "+formatcurrency.PrintUSD(sellExFiatVol))

		entry := ts.PriceQuery{
			Coin:                   1,
			Tier:                   id,
			Exchange_buy:           1,
			Exchange_sell:          2,
			Exchange_buy_coin:      3,
			Exchange_sell_coin:     4,
			Buy_query_time:         krakenBook.Time,
			Sell_query_time:        btcTurkBook.Time,
			Buy_input:              tier,
			Buy_count:              buyCount,
			Buy_vol:                cryptoOut,
			Sell_input:             cryptoOut,
			Sell_count:             sellCount,
			Sell_vol:               fiatOut,
			Buy_total_book_length:  buyExBookCount,
			Buy_total_book_crypto:  buyExCryptoVol,
			Buy_total_book_fiat:    buyExFiatVol,
			Sell_total_book_length: sellExBookLength,
			Sell_total_book_crypto: sellExCryptoVol,
			Sell_total_book_fiat:   sellExFiatVol,
			Exchange_rate:          exchangerate,
		}

		ts.InsertPriceQuery(&entry)
	}
}

func fetchResults(usdtiers ts.Result, trytiers ts.Result) {
	utils.P("green", "FIAT-CRYPTO PARITY ARBITRAGE DESK")
	utils.P("green", "-------------------------------------------\n")

	var book1 ex.Orderbook
	var book2 ex.Orderbook
	funcs := []GenericFunction{
		makeRequest(kraken.GetOrderBook, &book1),
		makeRequest(btcturk.GetOrderBook, &book2),
	}
	RunAsyncAllowErrors(funcs...)
	processArbitrageBuyKrakenSellBTCTurk([]ex.Orderbook{book1, book2}, usdtiers)
	processArbitrageBuyBTCTurkSellKraken([]ex.Orderbook{book1, book2}, trytiers)
}

func main() {
	usdtiers, _ := ts.GetTiers(3)
	trytiers, _ := ts.GetTiers(4)

	clock.Start(fetchResults, usdtiers, trytiers)
}
