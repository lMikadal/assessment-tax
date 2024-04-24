package tax

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (t Tax) TaxHandler(c echo.Context) error {
	var req ReqTax
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request"})
	}

	if req.TotalIncome <= 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "totalIncome must be greater than 0"})
	}

	if req.Wht < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Wht must be greater than 0"})
	}

	if req.Wht > req.TotalIncome {
		return c.JSON(http.StatusBadRequest, Err{Message: "Wht must be less than totalIncome"})
	}

	len_allowances := len(req.Allowances)
	if len_allowances > 2 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Allowances must be less than or equal to 2"})
	} else if len_allowances > 0 {
		allowance_type := []string{"donation", "k-receipt"}
		have_type := []string{}
		for _, v := range req.Allowances {
			allowance_type_low := strings.ToLower(v.AllowanceType)
			if ok := slices.Contains(allowance_type, allowance_type_low); !ok {
				return c.JSON(http.StatusBadRequest, Err{Message: "Not found allowanceType"})
			}
			if v.Amount < 0 {
				return c.JSON(http.StatusBadRequest, Err{Message: "Amount must be greater than 0"})
			}
			if ok := slices.Contains(have_type, allowance_type_low); ok {
				return c.JSON(http.StatusBadRequest, Err{Message: "Duplicate allowanceType"})
			}
			have_type = append(have_type, allowance_type_low)
		}
	}

	personal, err := t.info.GetTaxDeducationByType("Personal")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprintf("failed to get personal deduction: %v", err)})
	}

	req.TotalIncome -= personal.Amount
	if len_allowances > 0 {
		for _, v := range req.Allowances {
			deduction, err := t.info.GetTaxDeducationByType(cases.Title(language.English, cases.Compact).String(strings.ToLower(v.AllowanceType)))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprintf("failed to get deduction: %v", err)})
			}
			if v.Amount > deduction.Amount {
				req.TotalIncome -= deduction.Amount
			} else {
				req.TotalIncome -= v.Amount
			}
		}
	}

	tax_rate, err := t.info.GetTax()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprintf("failed to get tax rate: %v", err)})
	}
	var res ResTaxLevel
	var rang_now float64
	var cal float64
	for _, v := range tax_rate {
		rang_now = v.Maximum_salary - v.Minimum_salary
		if v.Rate != 0 {
			rang_now += 1
		}

		if req.TotalIncome <= 0 {
			cal = 0
		} else {
			if rang_now > req.TotalIncome || v.Maximum_salary == 0 {
				cal = t.calculateTax(req.TotalIncome, v)
			} else {
				cal = t.calculateTax(rang_now, v)
			}
			res.Tax += cal
			req.TotalIncome -= rang_now
		}

		t.addTaxLevel(&res.TaxLevel, v, cal)
	}

	res.Tax -= req.Wht
	if res.Tax < 0 {
		res.TaxRefund = res.Tax * -1
		res.Tax = 0
	}

	return c.JSON(http.StatusOK, res)
}

func (t Tax) TaxDeducateHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Err{Message: "For test only"})
}
