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

func TestGetAvailableLoans(t *testing.T) {
	tests := []struct {
		name           string
		customerLoan   models.CustomerLoan
		expectedLoans  []models.Loan
		expectedStatus int
		errorMessage   string
		wantErr        bool
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
			errorMessage:   "",
			wantErr:        false,
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
			errorMessage:   "",
			wantErr:        false,
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
			errorMessage:   "",
			wantErr:        false,
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
			errorMessage:   "",
			wantErr:        false,
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
			wantErr:        false,
		},
		{
			name: "Invalid customer: missing location",
			customerLoan: models.CustomerLoan{
				Age:    29,
				Cpf:    "12345678900",
				Name:   "John Doe",
				Income: 4400,
			},
			expectedLoans:  []models.Loan{},
			expectedStatus: http.StatusBadRequest,
			errorMessage:   "location is required",
			wantErr:        true,
		},
		{
			name: "Invalid customer: missing age",
			customerLoan: models.CustomerLoan{
				Cpf:      "12345678900",
				Name:     "John Doe",
				Income:   4400,
				Location: "SP",
			},
			expectedLoans:  []models.Loan{},
			expectedStatus: http.StatusBadRequest,
			errorMessage:   "age is required and must be greater than 0",
			wantErr:        true,
		},
		{
			name: "Invalid customer: missing name",
			customerLoan: models.CustomerLoan{
				Age:      29,
				Cpf:      "12345678900",
				Income:   4400,
				Location: "SP",
			},
			expectedLoans:  []models.Loan{},
			expectedStatus: http.StatusBadRequest,
			errorMessage:   "name is required",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.customerLoan)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/customer-loans", bytes.NewBuffer(body))
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			mux := http.NewServeMux()
			controllers.SetupRoutes(mux)
			mux.ServeHTTP(recorder, req)

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			var response struct {
				Customer string        `json:"customer"`
				Loans    []models.Loan `json:"loans"`
			}
			err = json.Unmarshal(recorder.Body.Bytes(), &response)
			if tt.wantErr {
				var errJson map[string]string
				_ = json.Unmarshal(recorder.Body.Bytes(), &errJson)
				assert.Empty(t, response.Customer)
				assert.Empty(t, response.Loans)
				assert.Equal(t, tt.errorMessage, errJson["error"])
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.customerLoan.Name, response.Customer)
			assert.Equal(t, tt.expectedLoans, response.Loans)
		})
	}
}
