package main

import (
	"log"
	"fmt"
	"github.com/fatih/color"
	"github.com/beldur/kraken-go-api-client"
	"github.com/ajph/bitstamp-go"
	"github.com/rodrigoalvesvieira/operatordesk/envars"
)

func main() {
	green := color.New(color.FgGreen, color.Bold)
	green.Println("OPERATOR DESK")
	green.Println("-------------------------------------------\n\n")

	// Color setups
	cyan := color.New(color.FgCyan).Add(color.Underline)
	yellow := color.New(color.FgYellow, color.Bold)

	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())
	bitstamp.SetAuth("123456", "key", "secret")

	balances, err := bitstamp.Ticker("btceur")
	if err != nil {
		log.Fatal("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	// GET TICKER FOR BISTAMP
	cyan.Println("BITSTAMP")
	yellow.Print("Current EUR balance: €")
	yellow.Println(0)
	yellow.Print("BTC ")
	yellow.Print(balances.Ask)
	yellow.Println()

	fmt.Println()

	// GET TICKER FOR KRAKEN
	ticker, err := api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Fatal(err)
	}

	balance, err := api.Balance()
	if err != nil {
		log.Fatal(err)
	}

	cyan.Println("KRAKEN")
	yellow.Print("Current EUR balance: €")
	yellow.Println(balance.ZEUR)
	yellow.Print("BTC ")
	yellow.Print(ticker.XXBTZEUR.Ask)
	yellow.Println()
}
