package utils

import (
	"fcpa/envars"
	ex "fcpa/exchange"
	"fmt"
	"math/big"
	"strconv"

	"github.com/fatih/color"
)

func Pln(message string) {
	if envars.IsProduction() {
		return
	}

	fmt.Println(message)
}

func P(colorName string, message string) {
	if envars.IsProduction() {
		return
	}

	colorOption := color.FgGreen
	option := color.Bold

	if colorName == "cyan" {
		colorOption = color.FgCyan
	}

	if colorName == "magenta" {
		colorOption = color.FgMagenta
	}

	if colorName == "yellow" {
		colorOption = color.FgYellow
	}

	if colorName == "red" {
		colorOption = color.FgRed
	}

	printer := color.New(colorOption, option)

	if colorName == "higreen" {
		color.HiGreen(message)
	} else {
		printer.Println(message)
	}

	// .Add(color.Underline)
}

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
