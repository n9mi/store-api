package repository

type Repositories struct {
	UserRepository            *UserRepository
	RoleRepository            *RoleRepository
	CategoryRepository        *CategoryRepository
	StoreRepository           *StoreRepository
	CustomerAddressRepository *CustomerAddressRepository
	PaymentMethodRepository   *PaymentMethodRepository
	ProductRepository         *ProductRepository
	OrderRepository           *OrderRepository
	OrderItemRepository       *OrderItemRepository
	CartItemRepository        *CartItemRepository
}

func Setup() *Repositories {
	return &Repositories{
		UserRepository:            NewUserRepository(),
		RoleRepository:            NewRoleRepository(),
		CategoryRepository:        NewCategoryRepository(),
		StoreRepository:           NewStoreRepository(),
		CustomerAddressRepository: NewCustomerAddressRepository(),
		PaymentMethodRepository:   NewPaymentMethodRepository(),
		ProductRepository:         NewProductRepository(),
		OrderRepository:           NewOrderRepository(),
		OrderItemRepository:       NewOrderItemRepository(),
		CartItemRepository:        NewCartItemRepository(),
	}
}
