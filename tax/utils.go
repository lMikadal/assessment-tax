package tax

import (
	"slices"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (t Tax) validateReq(req ReqTax) (bool, Err) {
	if req.TotalIncome <= 0 {
		return false, Err{Message: "totalIncome must be greater than 0"}
	}

	if req.Wht < 0 {
		return false, Err{Message: "Wht must be greater than 0"}
	}

	if req.Wht > req.TotalIncome {
		return false, Err{Message: "Wht must be less than totalIncome"}
	}

	len_allowances := len(req.Allowances)
	if len_allowances > 2 {
		return false, Err{Message: "Allowances must be less than or equal to 2"}
	} else if len_allowances > 0 {
		allowance_type := []string{"donation", "k-receipt"}
		have_type := []string{}
		for _, v := range req.Allowances {
			allowance_type_low := strings.ToLower(v.AllowanceType)
			if ok := slices.Contains(allowance_type, allowance_type_low); !ok {
				return false, Err{Message: "Not found allowanceType"}
			}
			if v.Amount < 0 {
				return false, Err{Message: "Amount must be greater than 0"}
			}
			if ok := slices.Contains(have_type, allowance_type_low); ok {
				return false, Err{Message: "Duplicate allowanceType"}
			}
			have_type = append(have_type, allowance_type_low)
		}
	}

	return true, Err{}
}

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
