package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		req          *http.Request
		expectedCode int
		expectedBody string
	}{
		{
			name:         "No credentials provided",
			req:          httptest.NewRequest(http.MethodGet, "/", nil),
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"message":"Unauthorized"}`,
		},
		{
			name:         "Incorrect credentials provided",
			req:          httptest.NewRequest(http.MethodGet, "/", nil),
			expectedCode: http.StatusUnauthorized,
			expectedBody: `{"message":"Unauthorized"}`,
		},
		{
			name:         "Correct credentials provided",
			req:          httptest.NewRequest(http.MethodGet, "/", nil),
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "Incorrect credentials provided" {
				tc.req.SetBasicAuth("incorrect", "credentials")
			} else if tc.name == "Correct credentials provided" {
				tc.req.SetBasicAuth("adminTax", "admin!")
			}

			rec := httptest.NewRecorder()
			c := echo.New().NewContext(tc.req, rec)

			err := middleware.AuthMiddleware(func(c echo.Context) error {
				return nil
			})(c)

			assert.Equal(t, tc.expectedCode, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			assert.NoError(t, err)
		})
	}
}

// func TestAuthMiddleware(t *testing.T) {
// 	// Set environment variables for testing
// 	os.Setenv("ADMIN_USERNAME", "adminTax")
// 	os.Setenv("ADMIN_PASSWORD", "admin!")

// 	// Create a new echo instance
// 	e := echo.New()

// 	// Create a new request with no credentials
// 	reqNoCreds := httptest.NewRequest(http.MethodGet, "/", nil)
// 	recNoCreds := httptest.NewRecorder()
// 	cNoCreds := e.NewContext(reqNoCreds, recNoCreds)

// 	// Test case: No credentials provided
// 	err := middleware.AuthMiddleware(func(c echo.Context) error {
// 		return nil
// 	})(cNoCreds)
// 	assert.Equal(t, http.StatusUnauthorized, recNoCreds.Code)
// 	assert.JSONEq(t, `{"message":"Unauthorized"}`, recNoCreds.Body.String())
// 	assert.NoError(t, err)

// 	// Create a new request with incorrect credentials
// 	reqIncorrectCreds := httptest.NewRequest(http.MethodGet, "/", nil)
// 	reqIncorrectCreds.SetBasicAuth("incorrect", "credentials")
// 	recIncorrectCreds := httptest.NewRecorder()
// 	cIncorrectCreds := e.NewContext(reqIncorrectCreds, recIncorrectCreds)

// 	// Test case: Incorrect credentials provided
// 	err = middleware.AuthMiddleware(func(c echo.Context) error {
// 		return nil
// 	})(cIncorrectCreds)
// 	assert.Equal(t, http.StatusUnauthorized, recIncorrectCreds.Code)
// 	assert.JSONEq(t, `{"message":"Unauthorized"}`, recIncorrectCreds.Body.String())
// 	assert.NoError(t, err)

// 	// Create a new request with correct credentials
// 	reqCorrectCreds := httptest.NewRequest(http.MethodGet, "/", nil)
// 	reqCorrectCreds.SetBasicAuth("adminTax", "admin!")
// 	recCorrectCreds := httptest.NewRecorder()
// 	cCorrectCreds := e.NewContext(reqCorrectCreds, recCorrectCreds)

// 	// Test case: Correct credentials provided
// 	err = middleware.AuthMiddleware(func(c echo.Context) error {
// 		return nil
// 	})(cCorrectCreds)
// 	assert.Equal(t, http.StatusOK, recCorrectCreds.Code)
// 	assert.NoError(t, err)
// }
