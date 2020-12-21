package exchange

import (
	"math/big"
)

type Coin string
type Fiat string

type OrderItem struct {
	Price  *big.Float // unit price
	Volume *big.Float // volume at this price
}

// Orderbook Represents the parsed order complete book with
// asks, bids, the exchange-issued timestamp and the timestamp
// on receipt
type Orderbook struct {
	CoinName     Coin
	FiatName     Fiat
	Asks         []OrderItem
	Bids         []OrderItem
	Time         int64
	ExchangeTime int64
}
