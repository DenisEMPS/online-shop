package domain

type CreateUserAdress struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
}
