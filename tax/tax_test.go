//go:build unit

package tax

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

type MockTax struct {
	db  []DB
	err error
}

func (m MockTax) TaxByIncome(income uint) ([]DB, error) {
	return m.db, m.err
}

func TestTaxHandler(t *testing.T) {
	t.Run("Test Income 500000", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			db: []DB{
				{Minimum_salary: 0, Maximum_salary: 150000, Rate: 0},
				{Minimum_salary: 150001, Maximum_salary: 500000, Rate: 10},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTax{Tax: 29000.0}
		gotJson := rec.Body.Bytes()
		// t.Logf("gotJson: %v", string(gotJson))

		var got ResTax
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

}