package tax

import (
	"fmt"
	"net/http"
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

	if ok, err := t.validateReq(req); !ok {
		return c.JSON(http.StatusBadRequest, err)
	}

	personal, err := t.info.GetTaxDeducationByType("Personal")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprintf("failed to get personal deduction: %v", err)})
	}

	req.TotalIncome -= personal.Amount
	if len(req.Allowances) > 0 {
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
	var req ReqAmount
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request"})
	}

	personal, err := t.info.GetTaxDeducationByType("Personal")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprintf("failed to get personal deduction: %v", err)})
	}
	if req.Amount > personal.Maximum_amount {
		return c.JSON(http.StatusBadRequest, Err{Message: "Amount should be less than 100,000"})
	} else if req.Amount < personal.Minimum_amount {
		return c.JSON(http.StatusBadRequest, Err{Message: "Amount should be more than 10,000"})
	}

	return c.JSON(http.StatusOK, Err{Message: "For test only"})
}
