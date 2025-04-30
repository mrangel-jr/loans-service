package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/mrangel-jr/loans-service/models"
)

type CustomerLoanResponse struct {
	Customer string        `json:"customer"`
	Loans    []models.Loan `json:"loans"`
}

func SetupRoutes(mux *http.ServeMux) {

	mux.HandleFunc("/customer-loans", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var customerLoan models.CustomerLoan
		err := json.NewDecoder(r.Body).Decode(&customerLoan)
		if err != nil {
			writeJSON(w, "Bad request", http.StatusBadRequest)
			return
		}

		if err := customerLoan.Validate(); err != nil {
			message := map[string]string{"error": err.Error()}
			writeJSON(w, message, http.StatusBadRequest)
			return
		}

		availableLoans := models.GetAvailableLoans(customerLoan)

		var response CustomerLoanResponse

		response.Loans = availableLoans
		response.Customer = customerLoan.Name

		writeJSON(w, response, http.StatusOK)
	})
}

func writeJSON(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
