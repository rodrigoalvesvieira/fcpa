package main

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	//"github.com/rodrigoalvesvieira/operatordesk/btcturk"
	"github.com/rodrigoalvesvieira/operatordesk/formatcurrency"
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
	green.Println("FIAT-CRYPTO PARITY ARBITRAGE DESK")
	green.Println("INPUT USD OUTPUT TRY")

	green.Println("-------------------------------------------\n")

	// USD ENTRY POINTS
	usdInTier1 := big.NewFloat(10000.0)
	usdInTier2 := big.NewFloat(100000.0)
	usdInTier3 := big.NewFloat(500000.0)

	//btcturk.GetFiatPayout(cryptoOut)

	// GETTING ORDER BOOK
	krakenBook := kraken.GetKrakenBook()
	krakenAskBook := kraken.AdaptOrderbook(krakenBook.Asks)

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier1))
	cryptoOut, fiatIn, orderCount, bookCrypto, bookFiat, bookCount := kraken.GetCryptoPayoutForBook(krakenAskBook, usdInTier1)

	magenta.Println("KRAKEN TOTAL TRAVERSED   " + strconv.Itoa(orderCount))
	magenta.Println("KRAKEN TOTAL FIAT        " + formatcurrency.PrintUSD(fiatIn))
	magenta.Println("KRAKEN TOTAL VOLUME      " + cryptoOut.String())

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier2))
	cryptoOut, fiatIn, orderCount, bookCrypto, bookFiat, bookCount = kraken.GetCryptoPayoutForBook(krakenAskBook, usdInTier2)

	magenta.Println("KRAKEN TOTAL TRAVERSED   " + strconv.Itoa(orderCount))
	magenta.Println("KRAKEN TOTAL FIAT        " + formatcurrency.PrintUSD(fiatIn))
	magenta.Println("KRAKEN TOTAL VOLUME      " + cryptoOut.String())

	fmt.Println("QUERYING " + formatcurrency.PrintUSD(usdInTier3))
	cryptoOut, fiatIn, orderCount, bookCrypto, bookFiat, bookCount = kraken.GetCryptoPayoutForBook(krakenAskBook, usdInTier3)

	magenta.Println("KRAKEN TOTAL TRAVERSED   " + strconv.Itoa(orderCount))
	magenta.Println("KRAKEN TOTAL FIAT        " + formatcurrency.PrintUSD(fiatIn))
	magenta.Println("KRAKEN TOTAL VOLUME      " + cryptoOut.String())

	cyan.Println("KRAKEN TOTAL BOOK VOLUME " + bookCrypto.String())
	cyan.Println("KRAKEN TOTAL BOOK LENGTH " + strconv.Itoa(bookCount))
	cyan.Println("KRAKEN TOTAL BOOK FIAT   " + formatcurrency.PrintUSD(bookFiat))

	//cryptoIn := big.NewFloat(5.0)
	//usdOut := getKrakenFiatPayout(cryptoIn)
	//magenta.Println("KRAKEN SELL BTC: TOTAL USD " + usdOut.String() + " FOR " + cryptoIn.String() + " BTC")

	//red := color.New(color.FgRed).Add(color.Bold)
	//red.Println("-----------------------------------------")
	//magenta.Println("-----------------------------------------")

	//green.Println("INPUT TRY OUTPUT USD")

	//cryptoIn = big.NewFloat(5.0)
	//tryOut := btcturk.GetCryptoPayout(cryptoIn)

	//green.Println("BTCTURK SELL BTC: TOTAL TRY " + tryOut.String() + " FOR " + cryptoIn.String() + " BTC")
}
