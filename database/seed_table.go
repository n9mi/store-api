package database

import (
	"fmt"
	"math/rand"
	"store-api/internal/entity"
	"store-api/internal/repository"
	"store-api/util"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB, repositories *repository.Repositories) error {
	roles, err := SeedRoles(db, repositories)
	if err != nil {
		return err
	}

	users, err := SeedUsers(db, repositories, roles)
	if err != nil {
		return err
	}

	categories, err := SeedCategories(db, repositories)
	if err != nil {
		return err
	}

	stores, err := SeedStores(db, repositories, users)
	if err != nil {
		return err
	}

	_, err = SeedCustomerAddresses(db, repositories, users)
	if err != nil {
		return err
	}

	_, err = SeedPaymentMethods(db, repositories)
	if err != nil {
		return err
	}

	_, err = SeedProducts(db, repositories, categories, stores)
	if err != nil {
		return err
	}

	return nil
}

func SeedRoles(db *gorm.DB, repositories *repository.Repositories) ([]entity.Role, error) {
	roles := []entity.Role{
		{ID: "role_merchant", Name: "merchant"},
		{ID: "role_customer", Name: "customer"},
	}

	tx := db.Begin()
	defer util.RecoverRollback(tx)
	for i, r := range roles {

		if err := repositories.RoleRepository.Create(tx, &r); err != nil {
			tx.Rollback()

			return nil, err
		}

		roles[i] = r
	}

	return roles, tx.Commit().Error
}

func SeedUsers(db *gorm.DB, repositories *repository.Repositories, roles []entity.Role) ([]entity.User, error) {
	var users []entity.User
	numPerRole := 3

	tx := db.Begin()
	defer util.RecoverRollback(tx)
	for _, r := range roles {
		for i := 1; i <= numPerRole; i++ {
			userPwd, _ := util.HashUserPassword("password")
			user := entity.User{
				ID:       uuid.NewString(),
				Name:     fmt.Sprintf("%s %d", r.Name, i),
				Email:    fmt.Sprintf("%s%d@test.com", r.Name, i),
				Password: userPwd,
			}
			if err := repositories.UserRepository.Create(tx, &user); err != nil {
				tx.Rollback()

				return nil, err
			}
			if err := repositories.UserRepository.AssignRole(tx, &user, &r); err != nil {
				tx.Rollback()

				return nil, err
			}
			users = append(users, user)
		}
	}

	return users, tx.Commit().Error
}

func SeedCategories(db *gorm.DB, repositories *repository.Repositories) ([]entity.Category, error) {
	categories := []entity.Category{
		{ID: uuid.NewString(), Name: "book"},
		{ID: uuid.NewString(), Name: "fashion"},
		{ID: uuid.NewString(), Name: "stationery"},
		{ID: uuid.NewString(), Name: "gardening tools"},
	}

	tx := db.Begin()
	defer util.RecoverRollback(tx)
	for i, c := range categories {
		if err := repositories.CategoryRepository.Create(tx, &c); err != nil {
			tx.Rollback()

			return nil, err
		}

		categories[i] = c
	}

	return categories, tx.Commit().Error
}

func SeedStores(db *gorm.DB, repositories *repository.Repositories, users []entity.User) ([]entity.Store, error) {
	var stores []entity.Store

	tx := db.Begin()
	defer util.RecoverRollback(tx)
	for _, u := range users {
		if isMerchant, err := repositories.UserRepository.HasRole(tx, u.ID, "role_merchant"); isMerchant &&
			err == nil {
			store := entity.Store{
				ID:       uuid.NewString(),
				Name:     fmt.Sprintf("Store %s", util.GenerateRandomString(5)),
				Street:   fmt.Sprintf("Street %s", util.GenerateRandomString(5)),
				City:     fmt.Sprintf("City %s", util.GenerateRandomString(5)),
				Province: fmt.Sprintf("Province %s", util.GenerateRandomString(5)),
				UserID:   u.ID,
			}

			if err := repositories.StoreRepository.Create(tx, &store); err != nil {
				tx.Rollback()

				return nil, err
			}

			stores = append(stores, store)
		} else if err != nil {
			tx.Rollback()

			return nil, err
		}
	}

	return stores, tx.Commit().Error
}

func SeedCustomerAddresses(db *gorm.DB, repositories *repository.Repositories, user []entity.User) ([]entity.CustomerAddress, error) {
	var customerAddresses []entity.CustomerAddress

	tx := db.Begin()
	defer util.RecoverRollback(tx)
	for _, u := range user {
		if isCustomer, err := repositories.UserRepository.HasRole(tx, u.ID, "role_customer"); isCustomer && err == nil {
			customerAddr := entity.CustomerAddress{
				ID:        uuid.NewString(),
				Street:    fmt.Sprintf("Street %s", util.GenerateRandomString(5)),
				City:      fmt.Sprintf("City %s", util.GenerateRandomString(5)),
				Province:  fmt.Sprintf("Province %s", util.GenerateRandomString(5)),
				IsDefault: true,
				UserID:    u.ID,
			}

			if err := repositories.CustomerAddressRepository.Create(tx, &customerAddr); err != nil {
				tx.Rollback()

				return nil, err
			}

			customerAddresses = append(customerAddresses, customerAddr)
		} else if err != nil {
			tx.Rollback()

			return nil, err
		}
	}

	return customerAddresses, tx.Commit().Error
}

func SeedPaymentMethods(db *gorm.DB, repositories *repository.Repositories) ([]entity.PaymentMethod, error) {
	var paymentMethods []entity.PaymentMethod
	numsPaymentMethods := 5

	tx := db.Begin()
	defer util.RecoverRollback(tx)
	for i := 1; i <= numsPaymentMethods; i++ {
		paymentMethod := entity.PaymentMethod{
			ID:   uuid.NewString(),
			Code: fmt.Sprintf("%d%d", i, i),
			Name: fmt.Sprintf("Payment method %d", i),
		}

		if err := repositories.PaymentMethodRepository.Create(tx, &paymentMethod); err != nil {
			tx.Rollback()

			return nil, err
		}

		paymentMethods = append(paymentMethods, paymentMethod)
	}

	return paymentMethods, tx.Commit().Error
}

func SeedProducts(db *gorm.DB, repositories *repository.Repositories, categories []entity.Category, stores []entity.Store) ([]entity.Product, error) {
	var products []entity.Product
	numProductsPerCtg := 3

	tx := db.Begin()
	defer util.RecoverRollback(tx)

	for _, s := range stores {
		for _, c := range categories {
			for i := 1; i <= numProductsPerCtg; i++ {
				product := entity.Product{
					ID:          uuid.NewString(),
					Name:        util.GenerateRandomString(5),
					Description: util.GenerateRandomString(20),
					PriceIdr:    float64(rand.Intn(99)+1) * 1000,
					Stock:       rand.Intn(99) + 1,
					CategoryID:  c.ID,
					StoreID:     s.ID,
				}

				if err := repositories.ProductRepository.Create(tx, &product); err != nil {
					tx.Rollback()

					return nil, err
				}

				products = append(products, product)
			}
		}
	}

	return products, tx.Commit().Error
}
