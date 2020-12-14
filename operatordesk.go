package main

import (
	"bytes"
	"runtime/debug"
	"sync"
	"fmt"
	"time"
	"strconv"
	"github.com/fatih/color"
	"operatordesk/formatcurrency"
	"operatordesk/btcturk"
	"operatordesk/kraken"
	ts "operatordesk/timeseries"
	"operatordesk/envars"
	ex "operatordesk/exchange"
	clock "operatordesk/clockwork"
)

// Alias for a function that takes no args and returns an error
type GenericFunction func() error

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

		Pln("RESPONSE TIME (THREADS): " + strconv.FormatInt(nsec, 10))

		res.Time = nsec

		*book = *res
		return nil
	}
}

func processArbitrageBuyBTCTurkSellKraken(books []ex.Orderbook, tiers ts.Result) {
	P("green", "INPUT TRY OUTPUT USD")

	krakenBook := books[0]
	btcTurkBook := books[1]

	for i := 0; i < len(tiers.Values); i++ {
		tier := tiers.Values[i]
		//id := tiers.Ids[i]

		P("red", "QUERYING " + formatcurrency.PrintTRY(tier))

		cryptoOut, fiatIn, orderCount, _, _, _ := btcturk.GetCryptoPayoutForBook(btcTurkBook.Asks, tier)

		P("magenta", "BTCTURK TOTAL TRAVERSED    " + strconv.Itoa(orderCount))
		P("magenta", "BTCTURK TOTAL FIAT         " + formatcurrency.PrintTRY(fiatIn))
		P("magenta", "BTCTURK TOTAL VOLUME       " + cryptoOut.String())

		cryptoIn, fiatOut, krakenOrderCount, _, _, _ := kraken.GetFiatPayoutForBook(krakenBook.Bids, cryptoOut)
		P("yellow", "KRAKEN TOTAL TRAVERSED   " + strconv.Itoa(krakenOrderCount))
		P("yellow", "KRAKEN TOTAL FIAT        " + formatcurrency.PrintUSD(fiatOut))
		P("yellow", "KRAKEN TOTAL VOLUME      " + cryptoIn.String())
	}
}

func Pln(message string) {
	if envars.IsProduction() {
		return
	}

	fmt.Println(message)
}

func P(colorName string, message string) {
	if envars.IsProduction() {
		return
	}

	colorOption := color.FgGreen
	option := color.Bold

	if colorName == "cyan" {
		colorOption = color.FgCyan
	}

	if colorName == "magenta" {
		colorOption = color.FgMagenta
	}

	if colorName == "yellow" {
		colorOption = color.FgYellow
	}

	if colorName == "red" {
		colorOption = color.FgRed
	}

	printer := color.New(colorOption, option)
	printer.Println(message)

	// .Add(color.Underline)
}

func processArbitrageBuyKrakenSellBTCTurk(books []ex.Orderbook, tiers ts.Result) {
	P("green", "INPUT USD OUTPUT TRY")

	krakenBook := books[0]
	btcTurkBook := books[1]

	for i := 0; i < len(tiers.Values); i++ {
		tier := tiers.Values[i]
		id := tiers.Ids[i]

		P("red", "QUERYING " + formatcurrency.PrintUSD(tier))

		cryptoOut, fiatIn, orderCount, bookCrypto, bookFiat, bookCount := kraken.GetCryptoPayoutForBook(krakenBook.Asks, tier)
		P("magenta", "KRAKEN TOTAL TRAVERSED    " + strconv.Itoa(orderCount))
		P("magenta", "KRAKEN TOTAL FIAT         " + formatcurrency.PrintUSD(fiatIn))
		P("magenta", "KRAKEN TOTAL VOLUME       " + cryptoOut.String())

		cryptoIn, fiatOut, btcTurkOrderCount, btcTurkBookCrypto, btcTurkBookFiat, btcTurkBookCount := btcturk.GetFiatPayoutForBook(btcTurkBook.Bids, cryptoOut)
		P("yellow", "QUERYING " + formatcurrency.PrintUSD(tier) + " " + cryptoOut.String())
		P("yellow", "BTCTURK TOTAL TRAVERSED   " + strconv.Itoa(btcTurkOrderCount))
		P("yellow", "BTCTURK TOTAL FIAT        " + formatcurrency.PrintTRY(fiatOut))
		P("yellow", "BTCTURK TOTAL VOLUME      " + cryptoIn.String())

		ts.InsertEntry(1, id, 1, 3, krakenBook.Time, orderCount, fiatIn, cryptoOut, bookCount, bookCrypto, bookFiat, btcTurkBookCount, btcTurkBookCrypto, btcTurkBookFiat, btcTurkOrderCount, fiatOut)

		ts.InsertEntry(1, id, 2, 4, btcTurkBook.Time, btcTurkOrderCount, cryptoIn, fiatOut, bookCount, bookCrypto, bookFiat, btcTurkBookCount, btcTurkBookCrypto, btcTurkBookFiat, btcTurkOrderCount, fiatOut)

		P("cyan", "KRAKEN TOTAL BOOK VOLUME  " + bookCrypto.String())
		P("cyan", "KRAKEN TOTAL BOOK LENGTH  " + strconv.Itoa(bookCount))
		P("cyan", "KRAKEN TOTAL BOOK FIAT    " + formatcurrency.PrintUSD(bookFiat))

		P("cyan", "BTCTURK TOTAL BOOK VOLUME " + btcTurkBookCrypto.String())
		P("cyan", "BTCTURK TOTAL BOOK LENGTH " + strconv.Itoa(btcTurkBookCount))
		P("cyan", "BTCTURK TOTAL BOOK FIAT   " + formatcurrency.PrintUSD(btcTurkBookFiat))
	}
}

func fetchResults(_ ts.Result, trytiers ts.Result) {
	P("green", "FIAT-CRYPTO PARITY ARBITRAGE DESK")
	P("green", "-------------------------------------------\n")

	var book1 ex.Orderbook
	var book2 ex.Orderbook
	funcs := []GenericFunction{
		makeRequest(kraken.GetOrderBook,  &book1),
		makeRequest(btcturk.GetOrderBook, &book2),
	}
    	RunAsyncAllowErrors(funcs...)
	//processArbitrageBuyKrakenSellBTCTurk([]ex.Orderbook{book1, book2}, usdtiers)
	processArbitrageBuyBTCTurkSellKraken([]ex.Orderbook{book1, book2}, trytiers)
}

func main() {
	usdtiers, _ := ts.GetTiers(3)
	trytiers, _ := ts.GetTiers(4)

	clock.Start(fetchResults, usdtiers, trytiers)
}
