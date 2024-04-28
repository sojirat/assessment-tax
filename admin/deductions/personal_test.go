package deductions_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/sojirat/assessment-tax/admin/deductions"
)

func TestUpdatePersonalDeductionHandler(t *testing.T) {
	tests := []struct {
		name           string
		inputAmount    float64
		expectedStatus int
	}{
		{
			name:           "Invalid amount (less than 10,000)",
			inputAmount:    5000.0,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Valid amount",
			inputAmount:    70000.0,
			expectedStatus: http.StatusOK,
		},
	}

	e := echo.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodPatch, "localhost:8080/admin/deductions/personal", strings.NewReader(fmt.Sprintf(`{"amount": %f}`, tt.inputAmount)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := deductions.UpdatePersonalDeductionHandler(c)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status %d; got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedStatus == http.StatusOK && err != nil {
				t.Errorf("Expected no error; got %v", err)
			}
		})
	}
}
