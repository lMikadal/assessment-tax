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

func TestDeductionHandler(t *testing.T) {
	t.Run("Test set amount over 100,000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqAmount{
			Amount: 100001.0,
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			dbDeduction: []DbDeduction{
				{
					Type:           "Personal",
					Maximum_amount: 100000.0,
					Minimum_amount: 10000.0,
				},
			},
		}

		handler := New(&mock)
		handler.TaxDeducateHandler(c)

		want := Err{Message: "Amount should be less than 100,000"}
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

	t.Run("Test set amount less than 10,000", func(t *testing.T) {
		e := echo.New()
		MockReq := ReqAmount{
			Amount: 9999.0,
		}
		reqBody, _ := json.Marshal(MockReq)
		req := httptest.NewRequest(http.MethodPost, "/admin/deductions/personal", bytes.NewBuffer(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mock := MockTax{
			dbDeduction: []DbDeduction{
				{
					Type:           "Personal",
					Maximum_amount: 100000.0,
					Minimum_amount: 10000.0,
				},
			},
		}

		handler := New(&mock)
		handler.TaxDeducateHandler(c)

		want := Err{Message: "Amount should be more than 10,000"}
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
