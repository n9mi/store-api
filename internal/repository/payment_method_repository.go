package repository

import "store-api/internal/entity"

type PaymentMethodRepository struct {
	BaseRepository[entity.PaymentMethod]
}

func NewPaymentMethodRepository() *PaymentMethodRepository {
	return new(PaymentMethodRepository)
}
