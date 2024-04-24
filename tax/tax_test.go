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
	dbDeduction []DbDeduction
	err         error
}

func (m MockTax) GetTax() ([]DB, error) {
	return []DB{
		{Minimum_salary: 0, Maximum_salary: 150000, Rate: 0},
		{Minimum_salary: 150001, Maximum_salary: 500000, Rate: 10},
		{Minimum_salary: 500001, Maximum_salary: 1000000, Rate: 15},
		{Minimum_salary: 1000001, Maximum_salary: 2000000, Rate: 20},
		{Minimum_salary: 2000001, Maximum_salary: 0, Rate: 35},
	}, m.err
}

func (m MockTax) GetTaxDeducationByType(deducation_type string) (DbDeduction, error) {
	for _, v := range m.dbDeduction {
		if v.Type == deducation_type {
			return v, nil
		}
	}
	return DbDeduction{}, m.err
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
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTaxLevel{
			Tax:       29000.0,
			TaxRefund: 0.0,
			TaxLevel: []TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   29000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   0.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResTaxLevel
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income 1000000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 1000000.0,
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
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTaxLevel{
			Tax:       101000.0,
			TaxRefund: 0.0,
			TaxLevel: []TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   35000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   66000.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResTaxLevel
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income 2000000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 2000000.0,
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
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTaxLevel{
			Tax:       298000.0,
			TaxRefund: 0.0,
			TaxLevel: []TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   35000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   75000.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   188000.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResTaxLevel
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income 3000000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 3000000.0,
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
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTaxLevel{
			Tax:       639000.0,
			TaxRefund: 0.0,
			TaxLevel: []TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   35000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   75000.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   200000.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   329000.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResTaxLevel
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

		mock := MockTax{}

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

		mock := MockTax{}

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

		mock := MockTax{}

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

	t.Run("Test AllowanceType no have", func(t *testing.T) {
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

		mock := MockTax{}

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

		mock := MockTax{}

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

	t.Run("Test duplication allowance type", func(t *testing.T) {
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

		mock := MockTax{}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := Err{Message: "Duplicate allowanceType"}
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

	t.Run("Test Allowance amount less than 0", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 1.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        -1.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := Err{Message: "Amount must be greater than 0"}
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

	t.Run("Test Income wht great tax", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 500000.0,
			Wht:         30000.0,
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
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTaxLevel{
			Tax:       0.0,
			TaxRefund: 1000.0,
			TaxLevel: []TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   29000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   0.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResTaxLevel
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income 50,000 and wht 25,000 output tax 4,000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 500000.0,
			Wht:         25000.0,
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
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTax{Tax: 4000.0}
		gotJson := rec.Body.Bytes()

		var got ResTax
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income 50,000 and donation 200,000 output tax 19,000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 500000.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        200000.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
				{
					ID:             1,
					Type:           "Donation",
					Minimum_amount: 0,
					Maximum_amount: 100000,
					Amount:         100000,
					Created_at:     "2021-09-01",
					Updated_at:     "2021-09-01",
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTaxLevel{
			Tax:       19000.0,
			TaxRefund: 0.0,
			TaxLevel: []TaxLevel{
				{
					Level: "0-150,000",
					Tax:   0.0,
				},
				{
					Level: "150,001-500,000",
					Tax:   19000.0,
				},
				{
					Level: "500,001-1,000,000",
					Tax:   0.0,
				},
				{
					Level: "1,000,001-2,000,000",
					Tax:   0.0,
				},
				{
					Level: "2,000,001 ขึ้นไป",
					Tax:   0.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResTaxLevel
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income 50,000 and donation 200,000 and Wht 20,000 output taxRefund 1,000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 500000.0,
			Wht:         20000.0,
			Allowances: []Allowance{
				{
					AllowanceType: "donation",
					Amount:        200000.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
				{
					Type:   "Donation",
					Amount: 100000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTaxRefund{TaxRefund: 1000.0}
		gotJson := rec.Body.Bytes()

		var got ResTaxRefund
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income 50,000 and K-Receipt 100,000 output tax 24,000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 500000.0,
			Wht:         0.0,
			Allowances: []Allowance{
				{
					AllowanceType: "k-receipt",
					Amount:        100000.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
				{
					Type:   "K-Receipt",
					Amount: 50000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTax{Tax: 24000.0}
		gotJson := rec.Body.Bytes()

		var got ResTax
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test Income 50,000 and K-Receipt 100,000 and have Wht 30,000 output taxRefund 6,000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqTax{
			TotalIncome: 500000.0,
			Wht:         30000.0,
			Allowances: []Allowance{
				{
					AllowanceType: "k-receipt",
					Amount:        100000.0,
				},
			},
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/tax/calculations", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			dbDeduction: []DbDeduction{
				{
					Type:   "Personal",
					Amount: 60000,
				},
				{
					Type:   "K-Receipt",
					Amount: 50000,
				},
			},
		}

		handler := New(&mock)
		handler.TaxHandler(c)

		want := ResTaxRefund{TaxRefund: 6000.0}
		gotJson := rec.Body.Bytes()

		var got ResTaxRefund
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

}
