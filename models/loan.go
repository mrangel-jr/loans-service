package models

import (
	"errors"
	"strings"
)

type Loan struct {
	Type           string `json:"type"`
	InterestedRate int    `json:"interested_rate"`
}

type CustomerLoan struct {
	Age      int     `json:"age"`
	Cpf      string  `json:"cpf"`
	Name     string  `json:"name"`
	Income   float64 `json:"income"`
	Location string  `json:"location"`
}

func getPersonalLoan(customerLoan CustomerLoan) []Loan {
	isEligible := customerLoan.Income <= 3000 || customerLoan.Age < 30 && customerLoan.Income > 3000 && customerLoan.Income < 5000 && customerLoan.Location == "SP"
	if isEligible {
		return []Loan{
			{
				Type:           "PERSONAL",
				InterestedRate: 4,
			},
		}
	}
	return []Loan{}
}

func getGuaranteedLoan(customerLoan CustomerLoan) []Loan {
	isEligible := customerLoan.Income <= 3000 || customerLoan.Age < 30 && customerLoan.Income > 3000 && customerLoan.Income < 5000 && customerLoan.Location == "SP"
	if isEligible {
		return []Loan{
			{
				Type:           "GUARANTEED",
				InterestedRate: 3,
			},
		}
	}
	return []Loan{}
}

func getConsigmentLoan(customerLoan CustomerLoan) []Loan {
	isEligible := customerLoan.Income >= 5000
	if isEligible {
		return []Loan{
			{
				Type:           "CONSIGMENT",
				InterestedRate: 2,
			},
		}
	}
	return []Loan{}
}

func GetAvailableLoans(customerLoan CustomerLoan) []Loan {
	// This function should contain the logic to determine available loans
	// based on the customerLoan details.
	loans := []Loan{}
	loans = append(loans, getPersonalLoan(customerLoan)...)
	loans = append(loans, getGuaranteedLoan(customerLoan)...)
	loans = append(loans, getConsigmentLoan(customerLoan)...)
	return loans
}

func (c CustomerLoan) Validate() error {
	if c.Age <= 0 {
		return errors.New("age is required and must be greater than 0")
	}
	if strings.TrimSpace(c.Cpf) == "" {
		return errors.New("cpf is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("name is required")
	}
	if c.Income <= 0 {
		return errors.New("income is required and must be greater than 0")
	}
	if strings.TrimSpace(c.Location) == "" {
		return errors.New("location is required")
	}
	return nil
}
