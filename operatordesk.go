package main

import (
	"bytes"
	"runtime/debug"
	"sync"
	"fmt"
	"math/big"
	"time"
	"strconv"
	"github.com/fatih/color"
	"operatordesk/formatcurrency"
	"operatordesk/btcturk"
	"operatordesk/kraken"
	"operatordesk/timeseries"
	ex "operatordesk/exchange"
	//_ "github.com/rodrigoalvesvieira/operatordesk/startup"
)

// export GOPRIVATE=github.com/rodrigoalvesvieira

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
		fmt.Print("RESPONSE TIME (THREADS): ")
		fmt.Println(nsec)

		res.Time = nsec

		*book = *res
		return nil
	}
}

func processArbitrage(books []ex.Orderbook) {
	krakenBook := books[0]
	btcTurkBook := books[1]

	fmt.Println(len(krakenBook.Asks))
	fmt.Println(len(btcTurkBook.Asks))

	cyan := color.New(color.FgCyan).Add(color.Underline)
	magenta := color.New(color.FgMagenta, color.Bold).Add(color.Underline)
	yellow := color.New(color.FgYellow, color.Bold)

	// USD ENTRY POINTS
	usdInTier1 := big.NewFloat(10000.0)
	usdInTier2 := big.NewFloat(100000.0)
	usdInTier3 := big.NewFloat(500000.0)

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier1))
	cryptoOutT1, fiatInT1, orderCountT1, bookCrypto, bookFiat, bookCount := kraken.GetCryptoPayoutForBook(krakenBook.Bids, usdInTier1)

	magenta.Println("KRAKEN TOTAL TRAVERSED    " + strconv.Itoa(orderCountT1))
	magenta.Println("KRAKEN TOTAL FIAT         " + formatcurrency.PrintUSD(fiatInT1))
	magenta.Println("KRAKEN TOTAL VOLUME       " + cryptoOutT1.String())

	cryptoIn, fiatOut, btcTurkOrderCount, btcTurkBookCrypto, btcTurkBookFiat, btcTurkBookCount := btcturk.GetFiatPayoutForBook(btcTurkBook.Bids, cryptoOutT1)
	yellow.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier1) + " " + cryptoOutT1.String())
	yellow.Println("BTCTURK TOTAL TRAVERSED   " + strconv.Itoa(btcTurkOrderCount))
	yellow.Println("BTCTURK TOTAL FIAT        " + formatcurrency.PrintTRY(fiatOut))
	yellow.Println("BTCTURK TOTAL VOLUME      " + cryptoIn.String())

	timeseries.InsertEntry(1, 1, 1, 3, krakenBook.Time, orderCountT1, fiatInT1, cryptoOutT1, bookCount, bookCrypto, bookFiat, btcTurkBookCount, btcTurkBookCrypto, btcTurkBookFiat)
	timeseries.InsertEntry(1, 3, 2, 4, btcTurkBook.Time, btcTurkOrderCount, cryptoIn, fiatOut, bookCount, bookCrypto, bookFiat, btcTurkBookCount, btcTurkBookCrypto, btcTurkBookFiat)

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier2))
	cryptoOutT2, fiatInT2, orderCountT2, _, _, _ := kraken.GetCryptoPayoutForBook(krakenBook.Asks, usdInTier2)

	magenta.Println("KRAKEN TOTAL TRAVERSED    " + strconv.Itoa(orderCountT2))
	magenta.Println("KRAKEN TOTAL FIAT         " + formatcurrency.PrintUSD(fiatInT2))
	magenta.Println("KRAKEN TOTAL VOLUME       " + cryptoOutT2.String())

	cryptoIn, fiatOut, btcTurkOrderCount, btcTurkBookCrypto, _, btcTurkBookCount = btcturk.GetFiatPayoutForBook(btcTurkBook.Bids, cryptoOutT2)
	yellow.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier2) + " " + cryptoOutT2.String())
	yellow.Println("BTCTURK TOTAL TRAVERSED   " + strconv.Itoa(btcTurkOrderCount))
	yellow.Println("BTCTURK TOTAL FIAT        " + formatcurrency.PrintTRY(fiatOut))
	yellow.Println("BTCTURK TOTAL VOLUME      " + cryptoIn.String())

	timeseries.InsertEntry(1, 2, 1, 3, krakenBook.Time, orderCountT2, fiatInT2, cryptoOutT2, bookCount, bookCrypto, bookFiat, btcTurkBookCount, btcTurkBookCrypto, btcTurkBookFiat)
	timeseries.InsertEntry(1, 2, 2, 4, btcTurkBook.Time, btcTurkOrderCount, cryptoIn, fiatOut, bookCount, bookCrypto, bookFiat, btcTurkBookCount, btcTurkBookCrypto, btcTurkBookFiat)

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier3))
	cryptoOutT3, fiatInT3, orderCountT3, _, _, _ := kraken.GetCryptoPayoutForBook(krakenBook.Asks, usdInTier3)

	magenta.Println("KRAKEN TOTAL TRAVERSED    " + strconv.Itoa(orderCountT3))
	magenta.Println("KRAKEN TOTAL FIAT         " + formatcurrency.PrintUSD(fiatInT3))
	magenta.Println("KRAKEN TOTAL VOLUME       " + cryptoOutT3.String())

	cryptoIn, fiatOut, btcTurkOrderCount, btcTurkBookCrypto, _, btcTurkBookCount = btcturk.GetFiatPayoutForBook(btcTurkBook.Bids, cryptoOutT3)
	yellow.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier3) + " " + cryptoOutT3.String())
	yellow.Println("BTCTURK TOTAL TRAVERSED   " + strconv.Itoa(btcTurkOrderCount))
	yellow.Println("BTCTURK TOTAL FIAT        " + formatcurrency.PrintTRY(fiatOut))
	yellow.Println("BTCTURK TOTAL VOLUME      " + cryptoIn.String())

	cyan.Println("KRAKEN TOTAL BOOK VOLUME  " + bookCrypto.String())
	cyan.Println("KRAKEN TOTAL BOOK LENGTH  " + strconv.Itoa(bookCount))
	cyan.Println("KRAKEN TOTAL BOOK FIAT    " + formatcurrency.PrintUSD(bookFiat))

	cyan.Println("BTCTURK TOTAL BOOK VOLUME " + btcTurkBookCrypto.String())
	cyan.Println("BTCTURK TOTAL BOOK LENGTH " + strconv.Itoa(btcTurkBookCount))
	cyan.Println("BTCTURK TOTAL BOOK FIAT   " + formatcurrency.PrintUSD(btcTurkBookFiat))

	timeseries.InsertEntry(1, 3, 1, 3, krakenBook.Time, orderCountT3, fiatInT3, cryptoOutT3, bookCount, bookCrypto, bookFiat, btcTurkBookCount, btcTurkBookCrypto, btcTurkBookFiat)
	timeseries.InsertEntry(1, 3, 2, 4, btcTurkBook.Time, btcTurkOrderCount, cryptoIn, fiatOut, bookCount, bookCrypto, bookFiat, btcTurkBookCount, btcTurkBookCrypto, btcTurkBookFiat)
}

func main() {
	green := color.New(color.FgGreen, color.Bold)

	green.Println("FIAT-CRYPTO PARITY ARBITRAGE DESK")
	green.Println("INPUT USD OUTPUT TRY")

	green.Println("-------------------------------------------\n")

	var book1 ex.Orderbook
	var book2 ex.Orderbook
	funcs := []GenericFunction{
		makeRequest(kraken.GetOrderBook,  &book1),
		makeRequest(btcturk.GetOrderBook, &book2),
	}
    	RunAsyncAllowErrors(funcs...)
	processArbitrage([]ex.Orderbook{book1, book2})
}
