package bitstamp

import (
	"log"
	"fmt"
	"github.com/fatih/color"
	"github.com/beldur/kraken-go-api-client"
	"github.com/ajph/bitstamp-go"
	"github.com/rodrigoalvesvieira/operatordesk/envars"
	"github.com/joho/godotenv"
	"math/big"
	"encoding/json"
	"io/ioutil"
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
