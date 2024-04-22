package postgres

import "github.com/lMikadal/assessment-tax/tax"

func (p *Postgres) TaxByIncome(income uint) (tax.DB, error) {
	var taxDb tax.DB

	return taxDb, nil
}
