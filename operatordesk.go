package main

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/rodrigoalvesvieira/operatordesk/formatcurrency"
	"github.com/rodrigoalvesvieira/operatordesk/btcturk"
	"github.com/rodrigoalvesvieira/operatordesk/kraken"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cyan := color.New(color.FgCyan).Add(color.Underline)
	magenta := color.New(color.FgMagenta, color.Bold).Add(color.Underline)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)

	green.Println("FIAT-CRYPTO PARITY ARBITRAGE DESK")
	green.Println("INPUT USD OUTPUT TRY")

	green.Println("-------------------------------------------\n")

	// USD ENTRY POINTS
	usdInTier1 := big.NewFloat(10000.0)
	usdInTier2 := big.NewFloat(100000.0)
	usdInTier3 := big.NewFloat(500000.0)

	//btcturk.GetFiatPayout(cryptoOut)

	// KRAKEN - GETTING ORDER BOOK
	krakenBook := kraken.GetKrakenBook()
	krakenAskBook := kraken.AdaptOrderbook(krakenBook.Asks)

	// BTCTURK - GETTING ORDER BOOK
	btcTurkBook := btcturk.GetOrderBook()
	btcTurkBidBook := btcturk.AdaptOrderbook(btcTurkBook.Bid)

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier1))
	cryptoOut, fiatIn, orderCount, bookCrypto, bookFiat, bookCount := kraken.GetCryptoPayoutForBook(krakenAskBook, usdInTier1)

	magenta.Println("KRAKEN TOTAL TRAVERSED    " + strconv.Itoa(orderCount))
	magenta.Println("KRAKEN TOTAL FIAT         " + formatcurrency.PrintUSD(fiatIn))
	magenta.Println("KRAKEN TOTAL VOLUME       " + cryptoOut.String())

	cryptoIn, fiatOut, btcTurkOrderCount, btcTurkBookCrypto, btcTurkBookFiat,  btcTurkBookCount := btcturk.GetFiatPayoutForBook(btcTurkBidBook, cryptoOut)
	yellow.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier1) + " " + cryptoOut.String())
	yellow.Println("BTCTURK TOTAL TRAVERSED   " + strconv.Itoa(btcTurkOrderCount))
	yellow.Println("BTCTURK TOTAL FIAT        " + formatcurrency.PrintTRY(fiatOut))
	yellow.Println("BTCTURK TOTAL VOLUME      " + cryptoIn.String())

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier2))
	cryptoOut, fiatIn, orderCount, bookCrypto, bookFiat, bookCount = kraken.GetCryptoPayoutForBook(krakenAskBook, usdInTier2)

	magenta.Println("KRAKEN TOTAL TRAVERSED    " + strconv.Itoa(orderCount))
	magenta.Println("KRAKEN TOTAL FIAT         " + formatcurrency.PrintUSD(fiatIn))
	magenta.Println("KRAKEN TOTAL VOLUME       " + cryptoOut.String())

	cryptoIn, fiatOut, btcTurkOrderCount, btcTurkBookCrypto, _, btcTurkBookCount = btcturk.GetFiatPayoutForBook(btcTurkBidBook, cryptoOut)
	yellow.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier2) + " " + cryptoOut.String())
	yellow.Println("BTCTURK TOTAL TRAVERSED   " + strconv.Itoa(btcTurkOrderCount))
	yellow.Println("BTCTURK TOTAL FIAT        " + formatcurrency.PrintTRY(fiatOut))
	yellow.Println("BTCTURK TOTAL VOLUME      " + cryptoIn.String())

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier3))
	cryptoOut, fiatIn, orderCount, bookCrypto, bookFiat, bookCount = kraken.GetCryptoPayoutForBook(krakenAskBook, usdInTier3)

	magenta.Println("KRAKEN TOTAL TRAVERSED    " + strconv.Itoa(orderCount))
	magenta.Println("KRAKEN TOTAL FIAT         " + formatcurrency.PrintUSD(fiatIn))
	magenta.Println("KRAKEN TOTAL VOLUME       " + cryptoOut.String())

	cryptoIn, fiatOut, btcTurkOrderCount, btcTurkBookCrypto, _, btcTurkBookCount = btcturk.GetFiatPayoutForBook(btcTurkBidBook, cryptoOut)
	yellow.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier3) + " " + cryptoOut.String())
	yellow.Println("BTCTURK TOTAL TRAVERSED   " + strconv.Itoa(btcTurkOrderCount))
	yellow.Println("BTCTURK TOTAL FIAT        " + formatcurrency.PrintTRY(fiatOut))
	yellow.Println("BTCTURK TOTAL VOLUME      " + cryptoIn.String())

	cyan.Println("KRAKEN TOTAL BOOK VOLUME  " + bookCrypto.String())
	cyan.Println("KRAKEN TOTAL BOOK LENGTH  " + strconv.Itoa(bookCount))
	cyan.Println("KRAKEN TOTAL BOOK FIAT    " + formatcurrency.PrintUSD(bookFiat))

	cyan.Println("BTCTURK TOTAL BOOK VOLUME " + btcTurkBookCrypto.String())
	cyan.Println("BTCTURK TOTAL BOOK LENGTH " + strconv.Itoa(btcTurkBookCount))
	cyan.Println("BTCTURK TOTAL BOOK FIAT   " + formatcurrency.PrintUSD(btcTurkBookFiat))
}
