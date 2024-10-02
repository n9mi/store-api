package dto

type CustomerAddressResponse struct {
	ID        string `json:"id"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Province  string `json:"province"`
	IsDefault bool   `json:"is_default"`
}
