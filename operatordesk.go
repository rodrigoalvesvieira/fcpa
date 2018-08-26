package main

import (
	"github.com/rodrigoalvesvieira/omnitrade"
	"github.com/rodrigoalvesvieira/omnitrade/core"
	"github.com/fatih/color"
	"log"

	"github.com/rodrigoalvesvieira/bitcointrade"
	"github.com/rodrigoalvesvieira/negociecoins"
)

func main() {
	green := color.New(color.FgGreen, color.Bold)
	green.Println("OPERATOR DESK")
	green.Println("-------------------------------------------\n\n\n")

	yellow := color.New(color.FgYellow, color.Bold)

	btcTickerOmnitrade, err := omnitrade.GetTicker(core.BITCOIN)
	if err != nil {
		log.Print(err)
	}

	btcTickerBitcoinTrade, err := bitcointrade.GetTicker()
	if err != nil {
		log.Print(err)
	}

	btcTickerNegocieCoins, err := negociecoins.GetTicker()
	if err != nil {
		log.Print(err)
	}

	yellow.Printf("%#v\r\n\n", btcTickerOmnitrade)
	yellow.Printf("%#v\r\n\n", btcTickerBitcoinTrade)
	yellow.Printf("%#v\r\n\n", btcTickerNegocieCoins)
}
