package tax

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type ReqTax struct {
	TotalIncome float64 `json:"totalIncome"`
	Wht         float64 `json:"wht"`
	Allowances  []Allowance
}

type ResTax struct {
	Tax float64 `json:"tax"`
}

type DB struct {
	ID             int     `postgres:"id"`
	Minimum_salary float64 `postgres:"minimum_salary"`
	Maximum_salary float64 `postgres:"maximum_salary"`
	Rate           float64 `postgres:"rate"`
	Created_at     string  `postgres:"created_at"`
}

type DbDeduction struct {
	ID             int     `postgres:"id"`
	Type           string  `postgres:"deducation_type"`
	Minimum_amount float64 `postgres:"minimum_amount"`
	Maximum_amount float64 `postgres:"maximum_amount"`
	Amount         float64 `postgres:"amount"`
	Created_at     string  `postgres:"created_at"`
	Updated_at     string  `postgres:"updated_at"`
}

type Tax struct {
	info InfoTax
}

type Err struct {
	Message string `json:"message"`
}

type InfoTax interface {
	TaxByIncome(income uint) ([]DB, error)
	GetTaxDeducation(deducation_type string) (DbDeduction, error)
}

func New(info InfoTax) Tax {
	return Tax{
		info: info,
	}
}

func (t Tax) TaxHandler(c echo.Context) error {
	var req ReqTax
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: "invalid request"})
	}

	if req.TotalIncome < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "TotalIncome must be greater than 0"})
	}

	if req.Wht < 0 {
		return c.JSON(http.StatusBadRequest, Err{Message: "Wht must be greater than 0"})
	}

	tax_rate, err := t.info.TaxByIncome(uint(req.TotalIncome))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprintf("failed to get tax rate: %v", err)})
	}
	var res ResTax
	var rang_now float64
	personal, err := t.info.GetTaxDeducation("Personal")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: fmt.Sprintf("failed to get personal deduction: %v", err)})
	}

	req.TotalIncome -= personal.Amount
	for _, v := range tax_rate {
		rang_now = v.Maximum_salary - v.Minimum_salary
		if v.Rate != 0 {
			rang_now += 1
		}
		if rang_now > req.TotalIncome {
			res.Tax += (req.TotalIncome * v.Rate) / 100
			break
		} else {
			res.Tax += (rang_now * v.Rate) / 100
			req.TotalIncome -= rang_now
		}
	}

	res.Tax -= req.Wht

	return c.JSON(http.StatusOK, res)
}
