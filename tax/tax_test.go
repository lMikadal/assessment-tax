//go:build unit

package tax

import (
	"bytes"
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

func (m MockTax) GetTaxDeducation(deducation_type string) (DbDeduction, error) {
	return DbDeduction{
		ID:             1,
		Type:           "Personal",
		Minimum_amount: 10000,
		Maximum_amount: 100000,
		Amount:         60000,
		Created_at:     "2021-09-01",
		Updated_at:     "2021-09-01",
	}, nil
}

func TestTaxHandler(t *testing.T) {
	t.Run("Test Income 500000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 500000.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
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

		var got ResTax
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income less than 0", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: -1.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			db: []DB{},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := Err{Message: "totalIncome must be greater than 0"}
		gotJson := rec.Body.Bytes()

		var got Err
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("got: %v, want: %v", rec.Code, http.StatusBadRequest)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Wht less than 0", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 1.0,
			Wht:         -1.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			db: []DB{},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := Err{Message: "Wht must be greater than 0"}
		gotJson := rec.Body.Bytes()

		var got Err
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("got: %v, want: %v", rec.Code, http.StatusBadRequest)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Wht great than totalIncome", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 1.0,
			Wht:         2.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			db: []DB{},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := Err{Message: "Wht must be less than totalIncome"}
		gotJson := rec.Body.Bytes()

		var got Err
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("got: %v, want: %v", rec.Code, http.StatusBadRequest)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test allowanceType no have", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 1.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "test",
					Amount:        0.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			db: []DB{},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := Err{Message: "Not found allowanceType"}
		gotJson := rec.Body.Bytes()

		var got Err
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("got: %v, want: %v", rec.Code, http.StatusBadRequest)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Allowances more than 2", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 1.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        0.0,
				},
				{
					AllowanceType: "k-receipt",
					Amount:        0.0,
				},
				{
					AllowanceType: "test",
					Amount:        0.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			db: []DB{},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := Err{Message: "Allowances must be less than or equal to 2"}
		gotJson := rec.Body.Bytes()

		var got Err
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if rec.Code != http.StatusBadRequest {
			t.Errorf("got: %v, want: %v", rec.Code, http.StatusBadRequest)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})
}
