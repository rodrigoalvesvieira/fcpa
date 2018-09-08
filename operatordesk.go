package main

import (
	"log"
	"fmt"
	"github.com/fatih/color"
	"github.com/beldur/kraken-go-api-client"
	"github.com/ajph/bitstamp-go"
)

func main() {
	green := color.New(color.FgGreen, color.Bold)
	green.Println("OPERATOR DESK")
	green.Println("-------------------------------------------\n\n")

	// Color setups
	cyan := color.New(color.FgCyan).Add(color.Underline)
	yellow := color.New(color.FgYellow, color.Bold)

	api := krakenapi.New("KEY", "SECRET")
	bitstamp.SetAuth("123456", "key", "secret")

	balances, err := bitstamp.Ticker("btceur")
	if err != nil {
		log.Fatal("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	// GET TICKER FOR BISTAMP
	cyan.Println("BITSTAMP")
	yellow.Print("BTC ")
	yellow.Print(balances.Ask)
	yellow.Println()

	fmt.Println()

	// GET TICKER FOR KRAKEN
	ticker, err := api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Fatal(err)
	}

	cyan.Println("KRAKEN")
	yellow.Print("BTC ")
	yellow.Print(ticker.XXBTZEUR.Ask)
	yellow.Println()
}
