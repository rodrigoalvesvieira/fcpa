package formatcurrency

import (
	"math/big"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// PrintUSD formats a decimal monetary value in United States Dollar format
func PrintUSD(total *big.Float) string {
	p := message.NewPrinter(language.English)
	res := p.Sprintf("%d", currency.NarrowSymbol(currency.USD.Amount(total)))

	return res
}

// PrintTRY formats a decimal monetary value in Turkish Lira format
func PrintTRY(total *big.Float) string {
	p := message.NewPrinter(language.English)
	res := p.Sprintf("%d", currency.NarrowSymbol(currency.TRY.Amount(total)))

	return res
}
