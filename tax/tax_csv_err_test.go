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

func TestCsvErrorHander(t *testing.T) {
	t.Run("Test first line wrong keyword", func(t *testing.T) {
		e := echo.New()

		body := new(bytes.Buffer)
		writer := csv.NewWriter(body)
		writer.Write([]string{"1"})
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

		want := Err{Message: "invalid csv"}
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

	t.Run("Test first line same keyword", func(t *testing.T) {
		e := echo.New()

		body := new(bytes.Buffer)
		writer := csv.NewWriter(body)
		writer.Write([]string{"totalIncome", "totalIncome"})
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

		want := Err{Message: "invalid csv"}
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
