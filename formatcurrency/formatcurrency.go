package formatcurrency

import (
	"math/big"
	"golang.org/x/text/message"
	"golang.org/x/text/language"
	"golang.org/x/text/currency"
)

func PrintUSD(total *big.Float) string {
	p := message.NewPrinter(language.English)
	res := p.Sprintf("%d", currency.NarrowSymbol(currency.USD.Amount(total)))

	return res
}


func PrintTRY(total *big.Float) string {
	p := message.NewPrinter(language.English)
	res := p.Sprintf("%d", currency.NarrowSymbol(currency.TRY.Amount(total)))

	return res
}
