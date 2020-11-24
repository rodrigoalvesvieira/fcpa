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

type Orderbook struct {
	CoinName Coin
	FiatName Fiat
	Asks     []OrderItem
	Bids     []OrderItem
}
