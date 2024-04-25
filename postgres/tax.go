package postgres

import "github.com/lMikadal/assessment-tax/tax"

func (p *Postgres) GetTax() ([]tax.DB, error) {
	rows, err := p.Db.Query("SELECT * FROM tax_rates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tax_rates []tax.DB
	for rows.Next() {
		var tax_rate tax.DB
		err := rows.Scan(&tax_rate.ID, &tax_rate.Minimum_salary, &tax_rate.Maximum_salary, &tax_rate.Rate, &tax_rate.Created_at)
		if err != nil {
			return nil, err
		}
		tax_rates = append(tax_rates, tax_rate)
	}

	return tax_rates, nil
}

func (p *Postgres) GetTaxDeducationByType(deducation_type string) (tax.DbDeduction, error) {
	rows, err := p.Db.Query("SELECT * FROM tax_deductions WHERE type = $1", deducation_type)
	if err != nil {
		return tax.DbDeduction{}, err
	}
	defer rows.Close()

	var tax_deduction tax.DbDeduction
	for rows.Next() {
		err := rows.Scan(&tax_deduction.ID, &tax_deduction.Type, &tax_deduction.Minimum_amount, &tax_deduction.Maximum_amount, &tax_deduction.Amount, &tax_deduction.Created_at, &tax_deduction.Updated_at)
		if err != nil {
			return tax.DbDeduction{}, err
		}
	}

	return tax_deduction, nil
}

func (p *Postgres) SetTaxDeducationByType(deducation_type string, amount float64) error {
	_, err := p.Db.Exec("UPDATE tax_deductions SET amount = $1 WHERE type = $2", amount, deducation_type)
	if err != nil {
		return err
	}

	return nil
}
