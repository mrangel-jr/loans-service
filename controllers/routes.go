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
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		availableLoans := models.GetAvailableLoans(customerLoan)

		var response CustomerLoanResponse

		response.Loans = availableLoans
		response.Customer = customerLoan.Name

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})
}
