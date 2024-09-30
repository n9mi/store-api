package database

import (
	"store-api/internal/entity"

	"gorm.io/gorm"
)

func CreateTable(db *gorm.DB) error {
	if err := db.AutoMigrate(&entity.Role{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Category{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Store{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.CustomerAddress{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.PaymentMethod{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Product{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.Order{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&entity.OrderItem{}); err != nil {
		return err
	}

	return nil
}

func DropTable(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&entity.OrderItem{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Order{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Product{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.PaymentMethod{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.CustomerAddress{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Store{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Category{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable("user_roles"); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.User{}); err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&entity.Role{}); err != nil {
		return err
	}

	return nil
}
