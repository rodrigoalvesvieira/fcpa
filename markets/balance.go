package balance

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
