package tax_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/tax"
	"github.com/stretchr/testify/assert"
)

func TestReadCSVFileHandler(t *testing.T) {
	testCases := []struct {
		name           string
		csvContent     string
		expectedBody   string
		expectedCode   int
		expectedErrMsg string
	}{
		{
			name:         "Valid CSV Content",
			csvContent:   "TotalIncome,WHT,Allowances\n1000,50,100\n2000,100,150",
			expectedBody: `{"Taxes":[{"TotalIncome":1000,"Tax":200},{"TotalIncome":2000,"Tax":400}]}`,
			expectedCode: http.StatusOK,
		},
		{
			name:           "Invalid File Name",
			csvContent:     "TotalIncome,WHT,Allowances\n1000,50,100\n2000,100,150",
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "filename must be taxes.csv",
		},
		{
			name:           "Invalid CSV Content",
			csvContent:     "",
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "value must not be empty",
		},
		{
			name:           "Empty Value",
			csvContent:     "TotalIncome,WHT,Allowances\n1000,50,\n2000,100,150",
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "value must not be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "localhost:8080/tax/calculations/upload-csv", strings.NewReader(tc.csvContent))
			req.Header.Set("Content-Type", "multipart/form-data; boundary=--------------------------1234567890")

			e := echo.New()
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := tax.ReadCSVFileHandler(c)

			assert.Equal(t, tc.expectedCode, rec.Code)

			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			}

			if tc.expectedErrMsg != "" {
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			}
		})
	}
}
