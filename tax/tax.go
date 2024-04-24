package tax

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

type ResTaxRefund struct {
	TaxRefund float64 `json:"taxRefund"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

type ResTaxLevel struct {
	Tax       float64 `json:"tax"`
	TaxRefund float64 `json:"taxRefund"`
	TaxLevel  []TaxLevel
}

type DB struct {
	ID             int     `postgres:"id"`
	Minimum_salary float64 `postgres:"minimum_salary"`
	Maximum_salary float64 `postgres:"maximum_salary|NULL"`
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
	GetTax() ([]DB, error)
	GetTaxDeducationByType(deducation_type string) (DbDeduction, error)
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
