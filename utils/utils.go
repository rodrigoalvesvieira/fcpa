package utils

import (
	"math/big"
)

func SubtractFee(volume *big.Float, fee *big.Float) *big.Float {
	feeCost := big.NewFloat(0.0).Copy(volume)
	feeCost.Mul(feeCost, fee)

	vol := big.NewFloat(0.0).Copy(volume)
	vol.Sub(vol, feeCost)

	return vol
}
