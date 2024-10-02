package dto

type CustomerAddressRequest struct {
	Street    string `json:"street" validate:"required,min=1,max=100"`
	City      string `json:"city" validate:"required,min=1,max=100"`
	Province  string `json:"province" validate:"required,min=1,max=100"`
	IsDefault bool   `json:"is_default"`
}

type CustomerAddressResponse struct {
	ID        string `json:"id"`
	Street    string `json:"street"`
	City      string `json:"city"`
	Province  string `json:"province"`
	IsDefault bool   `json:"is_default"`
}
