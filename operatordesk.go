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
	bitstamp.SetAuth("912534", envars.GetBitstampAPIKey(), envars.GetBitstampAPISecret())

	bitstampBalance, err := bitstamp.AccountBalance()
	if err != nil {
		fmt.Printf("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	bitstampTicker, err := bitstamp.Ticker("btceur")
	if err != nil {
		log.Fatal("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	// GET TICKER FOR BISTAMP
	cyan.Println("BITSTAMP")
	yellow.Print("Current EUR balance: €")
	yellow.Println(bitstampBalance.EurBalance)
	yellow.Print("BTC ")
	yellow.Print(bitstampTicker.Ask)
	yellow.Println()

	fmt.Println()

	// GET TICKER FOR KRAKEN
	krakenTicker, err := api.Ticker(krakenapi.XXBTZEUR)
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
	yellow.Print(krakenTicker.XXBTZEUR.Ask)
	yellow.Println()
}
