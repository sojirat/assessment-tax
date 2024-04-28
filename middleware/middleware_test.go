package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	url := "localhost:8080/admin/deductions/k-receipt"

	tests := []struct {
		name         string
		req          *http.Request
		expectedCode int
	}{
		{
			name:         "No credentials provided",
			req:          httptest.NewRequest(http.MethodPatch, url, nil),
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Incorrect credentials provided",
			req:          httptest.NewRequest(http.MethodPatch, url, nil),
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Correct credentials provided",
			req:          httptest.NewRequest(http.MethodPatch, url, strings.NewReader(`{"amount": 50000}`)),
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(tc.req, rec)

			if tc.name == "Incorrect credentials provided" {
				tc.req.SetBasicAuth("incorrect", "credentials")
			} else if tc.name == "Correct credentials provided" {
				tc.req.SetBasicAuth("adminTax", "admin!")
			}

			os.Setenv("ADMIN_USERNAME", "adminTax")
			os.Setenv("ADMIN_PASSWORD", "admin!")

			err := middleware.AuthMiddleware(func(c echo.Context) error {
				return nil
			})(c)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.NoError(t, err)
		})
	}
}
