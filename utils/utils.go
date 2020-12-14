package utils

import (
	ex "fcpa/exchange"
	"fmt"
	"math/big"
	"strconv"
)

func SubtractFee(volume *big.Float, fee *big.Float) *big.Float {
	feeCost := big.NewFloat(0.0).Copy(volume)
	feeCost.Mul(feeCost, fee)

	vol := big.NewFloat(0.0).Copy(volume)
	vol.Sub(vol, feeCost)

	return vol
}

func PrintOrder(i int, volume *big.Float, price *big.Float) {
	fmt.Println("Order #" + strconv.Itoa(i+1) + " Volume: " + volume.String() + " Price: " + price.String())
}

func GetOrderCost(order ex.OrderItem) *big.Float {
	base := big.NewFloat(0.0).Copy(order.Price)
	return base.Mul(base, order.Volume)
}
