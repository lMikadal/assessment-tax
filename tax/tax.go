package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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
	Minimum_salary float64 `postgres:"minimum_salary"`
	Maximum_salary float64 `postgres:"maximum_salary"`
	Rate           float64 `postgres:"rate"`
}

type Tax struct {
	info InfoTax
}

type Err struct {
	Message string `json:"message"`
}

type InfoTax interface {
	TaxByIncome(income uint) ([]DB, error)
}

func New(info InfoTax) Tax {
	return Tax{
		info: info,
	}
}

func (t Tax) TaxHandler(c echo.Context) error {
	// var res DB

	// res, _ = t.info.TaxByIncome(0)

	return c.JSON(http.StatusOK, Err{Message: "success"})
}
