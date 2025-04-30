package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mrangel-jr/loans-service/controllers"
	"github.com/mrangel-jr/loans-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupServer(t *testing.T, customerLoan models.CustomerLoan) *httptest.ResponseRecorder {
	body, err := json.Marshal(customerLoan)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "/customer-loans", bytes.NewBuffer(body))
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	mux := http.NewServeMux()
	controllers.SetupRoutes(mux)
	mux.ServeHTTP(recorder, req)
	return recorder
}

func TestGetAvailableLoans(t *testing.T) {
	tests := []struct {
		name           string
		customerLoan   models.CustomerLoan
		expectedLoans  []models.Loan
		expectedStatus int
	}{
		{
			name: "Valid customer with personal and guaranteed loans",
			customerLoan: models.CustomerLoan{
				Age:      25,
				Cpf:      "12345678900",
				Name:     "John Doe",
				Income:   4000,
				Location: "SP",
			},
			expectedLoans: []models.Loan{
				{
					Type:           "PERSONAL",
					InterestedRate: 4,
				},
				{
					Type:           "GUARANTEED",
					InterestedRate: 3,
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid customer with consignment loan",
			customerLoan: models.CustomerLoan{
				Age:      35,
				Cpf:      "12345678900",
				Name:     "John Doe",
				Income:   6000,
				Location: "RJ",
			},
			expectedLoans: []models.Loan{
				{
					Type:           "CONSIGMENT",
					InterestedRate: 2,
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid customer with personal and guaranteed loans by income",
			customerLoan: models.CustomerLoan{
				Age:      35,
				Cpf:      "12345678900",
				Name:     "John Doe",
				Income:   2700,
				Location: "RJ",
			},
			expectedLoans: []models.Loan{
				{
					Type:           "PERSONAL",
					InterestedRate: 4,
				},
				{
					Type:           "GUARANTEED",
					InterestedRate: 3,
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid customer with personal and guaranteed loans by income, location and age",
			customerLoan: models.CustomerLoan{
				Age:      29,
				Cpf:      "12345678900",
				Name:     "John Doe",
				Income:   4400,
				Location: "SP",
			},
			expectedLoans: []models.Loan{
				{
					Type:           "PERSONAL",
					InterestedRate: 4,
				},
				{
					Type:           "GUARANTEED",
					InterestedRate: 3,
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Valid customer without loans",
			customerLoan: models.CustomerLoan{
				Age:      29,
				Cpf:      "12345678900",
				Name:     "John Doe",
				Income:   4400,
				Location: "RJ",
			},
			expectedLoans:  []models.Loan{},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := setupServer(t, tt.customerLoan)

			var response struct {
				Customer string        `json:"customer"`
				Loans    []models.Loan `json:"loans"`
			}

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			err := json.Unmarshal(recorder.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tt.customerLoan.Name, response.Customer)
			assert.Equal(t, tt.expectedLoans, response.Loans)
		})
	}
}
func TestInvalidInGetAvailableLoans(t *testing.T) {
	tests := []struct {
		name           string
		customerLoan   models.CustomerLoan
		expectedStatus int
		errorMessage   string
	}{

		{
			name: "Invalid customer: missing location",
			customerLoan: models.CustomerLoan{
				Age:    29,
				Cpf:    "12345678900",
				Name:   "John Doe",
				Income: 4400,
			},
			expectedStatus: http.StatusBadRequest,
			errorMessage:   "location is required",
		},
		{
			name: "Invalid customer: missing age",
			customerLoan: models.CustomerLoan{
				Cpf:      "12345678900",
				Name:     "John Doe",
				Income:   4400,
				Location: "SP",
			},
			expectedStatus: http.StatusBadRequest,
			errorMessage:   "age is required and must be greater than 0",
		},
		{
			name: "Invalid customer: missing name",
			customerLoan: models.CustomerLoan{
				Age:      29,
				Cpf:      "12345678900",
				Income:   4400,
				Location: "SP",
			},
			expectedStatus: http.StatusBadRequest,
			errorMessage:   "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var errJson map[string]string

			recorder := setupServer(t, tt.customerLoan)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			_ = json.Unmarshal(recorder.Body.Bytes(), &errJson)
			assert.Equal(t, tt.errorMessage, errJson["error"])

		})
	}
}
