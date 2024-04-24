package tax

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (t Tax) calculateTax(income float64, rate DB) float64 {
	cal := (income * rate.Rate) / 100

	return cal
}

func (t Tax) addTaxLevel(level *[]TaxLevel, rate DB, cal float64) {
	var format string
	newP := message.NewPrinter(language.English)
	if rate.Maximum_salary != 0 {
		format = newP.Sprintf("%d-%d", int(rate.Minimum_salary), int(rate.Maximum_salary))
	} else {
		format = newP.Sprintf("%d ขึ้นไป", int(rate.Minimum_salary))

	}

	*level = append(*level, TaxLevel{
		Level: format,
		Tax:   cal,
	})
}
