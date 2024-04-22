package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Tax struct {
	info InfoTax
}

type Allowance struct {
	AllowanceType string  `json:"allowance_type"`
	Amount        float64 `json:"amount"`
}

type ReqTax struct {
	TotalIncome float64 `json:"total_income"`
	Wht         float64 `json:"wht"`
	Allowances  []Allowance
}

type ResTax struct {
	Tax float64 `json:"tax"`
}

type DB struct {
	Minimum_salary float64 `json:"minimum_salary"`
	Maximum_salary float64 `json:"maximum_salary"`
	Rate           float64 `json:"rate"`
}

func New(info InfoTax) Tax {
	return Tax{
		info: info,
	}
}

type InfoTax interface {
	TaxByIncome(income uint) ([]DB, error)
}

type Err struct {
	Message string `json:"message"`
}

func (t Tax) TaxHandler(c echo.Context) error {
	// var res DB

	// res, _ = t.info.TaxByIncome(0)

	return c.JSON(http.StatusOK, Err{Message: "success"})
}
