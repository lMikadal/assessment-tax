package tax

import (
	"slices"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
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

func (t Tax) validateCsv(head []string) (map[string]int, map[string]float64, Err) {
	simple := []string{"totalIncome", "wht", "donation", "k-receipt"}
	position := make(map[string]int)
	deducate := make(map[string]float64)

	for i, v := range head {
		if _, ok := position[v]; !ok && slices.Contains(simple, v) {
			position[v] = i
		} else {
			return make(map[string]int), make(map[string]float64), Err{Message: "invalid csv"}
		}

		if v == "donation" || v == "k-receipt" {
			d, err := t.info.GetTaxDeducationByType(cases.Title(language.English, cases.Compact).String(strings.ToLower(v)))
			if err != nil {
				return make(map[string]int), make(map[string]float64), Err{Message: "failed to get deduction"}
			}
			deducate[v] = d.Amount
		}
	}
	if _, ok := position["totalIncome"]; !ok {
		return make(map[string]int), make(map[string]float64), Err{Message: "invalid csv have not totalIncome"}
	}

	personal, err := t.info.GetTaxDeducationByType("Personal")
	if err != nil {
		return make(map[string]int), make(map[string]float64), Err{Message: "failed to get deduction"}
	}
	deducate["personal"] = personal.Amount

	return position, deducate, Err{}
}

func (t Tax) calDeducation(p map[string]int, name_p string, str []string, deducate map[string]float64) (float64, Err) {
	if _, ok := p[name_p]; ok {
		dedu, err := strconv.ParseFloat(str[p[name_p]], 64)
		if err != nil {
			return 0.0, Err{Message: "invalid donation"}
		}
		if dedu > deducate[name_p] {
			return deducate[name_p], Err{}
		} else {
			return dedu, Err{}
		}
	}

	return 0.0, Err{}
}
