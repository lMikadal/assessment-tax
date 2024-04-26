//go:build unit

package tax

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestCsvHander(t *testing.T) {
	t.Run("Test normal story6", func(t *testing.T) {
		e := echo.New()

		body := new(bytes.Buffer)
		writer := csv.NewWriter(body)
		writer.Write([]string{"totalIncome", "wht", "donation"})
		writer.Write([]string{"500000", "0", "0"})
		writer.Write([]string{"600000", "40000", "20000"})
		writer.Write([]string{"750000", "50000", "15000"})
		writer.Flush()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
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
		handler.UploadCSVHandler(c)

		want := ResAllCsv{
			Taxes: []ResCsvTax{
				{
					TotalIncome: 500000.0,
					Tax:         29000.0,
				},
				{
					TotalIncome: 600000.0,
					TaxRefund:   2000.0,
				},
				{
					TotalIncome: 750000.0,
					Tax:         11250.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResAllCsv
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("got: %v, want: %v", rec.Code, http.StatusOK)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test another stroy", func(t *testing.T) {
		e := echo.New()

		body := new(bytes.Buffer)
		writer := csv.NewWriter(body)
		writer.Write([]string{"totalIncome", "wht", "donation"})
		writer.Write([]string{"1000000", "0", "0"})
		writer.Write([]string{"500000", "25000", "0"})
		writer.Write([]string{"500000", "0", "200000"})
		writer.Flush()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
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
		handler.UploadCSVHandler(c)

		want := ResAllCsv{
			Taxes: []ResCsvTax{
				{
					TotalIncome: 1000000.0,
					Tax:         101000.0,
				},
				{
					TotalIncome: 500000.0,
					Tax:         4000.0,
				},
				{
					TotalIncome: 500000.0,
					Tax:         19000.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResAllCsv
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("got: %v, want: %v", rec.Code, http.StatusOK)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})

	t.Run("Test have field totalIncome", func(t *testing.T) {
		e := echo.New()

		body := new(bytes.Buffer)
		writer := csv.NewWriter(body)
		writer.Write([]string{"totalIncome"})
		writer.Write([]string{"1000000"})
		writer.Write([]string{"2000000"})
		writer.Write([]string{"3000000"})
		writer.Flush()

		req := httptest.NewRequest(http.MethodPost, "/tax/calculations/upload-csv", body)
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
		handler.UploadCSVHandler(c)

		want := ResAllCsv{
			Taxes: []ResCsvTax{
				{
					TotalIncome: 1000000.0,
					Tax:         101000.0,
				},
				{
					TotalIncome: 2000000.0,
					Tax:         298000.0,
				},
				{
					TotalIncome: 3000000.0,
					Tax:         639000.0,
				},
			},
		}
		gotJson := rec.Body.Bytes()

		var got ResAllCsv
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("failed to unmarshal json: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("got: %v, want: %v", rec.Code, http.StatusOK)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	})
}
