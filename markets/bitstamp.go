package bitstamp

import (
	"fcpa/envars"
	"fmt"
	"log"

	"github.com/ajph/bitstamp-go"
	"github.com/fatih/color"
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
