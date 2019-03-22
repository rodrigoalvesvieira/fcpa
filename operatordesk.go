package main

import (
	"log"
	"fmt"
	"github.com/fatih/color"
	"github.com/beldur/kraken-go-api-client"
	"github.com/ajph/bitstamp-go"
	"github.com/rodrigoalvesvieira/operatordesk/envars"
)

func bitstampQuery(cl *color.Color) {
	bitstamp.SetAuth("912534", envars.GetBitstampAPIKey(), envars.GetBitstampAPISecret())

	bitstampEURBalance, err := bitstamp.AccountBalance()
	if err != nil {
		fmt.Printf("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	bitstampBTCTicker, err := bitstamp.Ticker("btceur")
	if err != nil {
		log.Fatal("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	bitstampLTCTicker, err := bitstamp.Ticker("ltceur")
	if err != nil {
		log.Fatal("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	bitstampETHTicker, err := bitstamp.Ticker("etheur")
	if err != nil {
		log.Fatal("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	bitstampBCHTicker, err := bitstamp.Ticker("bcheur")
	if err != nil {
		log.Fatal("Can't get balance using bitstamp API: %s\n", err)
		return
	}

	row := fmt.Sprintf("Bitstamp      | %f    | %f  | %f    | %f    | %f    ", bitstampEURBalance.EurBalance, bitstampBTCTicker.Ask, bitstampLTCTicker.Ask, bitstampETHTicker.Ask, bitstampBCHTicker.Ask)
	cl.Println(row)
}

func krakenQuery(cl *color.Color) {
	api := krakenapi.New(envars.GetKrakenAPIKey(), envars.GetKrakenPrivateKey())

	// KRAKEN API INTEGRATION
	krakenEURBalance, err := api.Balance()
	if err != nil {
		log.Fatal(err)
	}

	krakenBTCTicker, err := api.Ticker(krakenapi.XXBTZEUR)
	if err != nil {
		log.Fatal(err)
	}

	krakenLTCTicker, err := api.Ticker(krakenapi.XLTCZEUR)
	if err != nil {
		log.Fatal(err)
	}

	krakenETHTicker, err := api.Ticker(krakenapi.XETHZEUR)
	if err != nil {
		log.Fatal(err)
	}

	krakenBCHTicker, err := api.Ticker(krakenapi.BCHEUR)
	if err != nil {
		log.Fatal(err)
	}

	row := fmt.Sprintf("Kraken        | %f   | %s   | %s     | %s     | %s    ", krakenEURBalance.ZEUR, krakenBTCTicker.XXBTZEUR.Ask[0], krakenLTCTicker.XLTCZEUR.Ask[0], krakenETHTicker.XETHZEUR.Ask[0], krakenBCHTicker.BCHEUR.Ask[0])

	cl.Println(row)
}

func main() {
	green := color.New(color.FgGreen, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)

	green.Println("OPERATOR DESK")
	green.Println("-------------------------------------------\n\n")

	// Color setups
	cyan := color.New(color.FgCyan).Add(color.Underline)
	yellow := color.New(color.FgYellow, color.Bold)

	row := "EXCHANGE      | EUR BALANCE | BTC          | LTC          | ETH           | BCH          "
	magenta.Println(row)

	go bitstampQuery(yellow)
	krakenQuery(yellow)

	fmt.Println()
	cyan.Println("END")
}
