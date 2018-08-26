package main

import (
	"github.com/rodrigoalvesvieira/omnitrade"
	"github.com/rodrigoalvesvieira/omnitrade/orderbook"
	"github.com/rodrigoalvesvieira/omnitrade/core"
	"github.com/fatih/color"
	"log"
)

func main() {
	green := color.New(color.FgGreen, color.Bold)
	green.Println("OPERATOR DESK")
	green.Println("-------------------------------------------\n\n\n")

	yellow := color.New(color.FgYellow, color.Bold)

	btcTicker, err := omnitrade.GetTicker(core.BITCOIN)
	if err != nil {
		log.Print(err)
	}

	yellow.Printf("%#v\r\n\n", btcTicker)

	orderbook.GetOrderBook(core.BITCOIN)
}
