package tax

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Tax struct {
	info InfoTax
}

type DB struct {
	Minimum_salary int `postgres:"minimum_salary"`
	Maximum_salary int `postgres:"maximum_salary"`
	Rate           int `postgres:"rate"`
}

func New(info InfoTax) Tax {
	return Tax{
		info: info,
	}
}

type InfoTax interface {
	TaxByIncome(income uint) (DB, error)
}

type Err struct {
	Message string `json:"message"`
}

func (t Tax) TaxHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, Err{Message: "test ok"})
}
