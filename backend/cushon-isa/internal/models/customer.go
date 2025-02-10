package models

type RetailCustomer struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

type CreateRetailCustomerRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
}

func NewRetailCustomer(firstName, lastName, email string) RetailCustomer {
	return RetailCustomer{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
}
