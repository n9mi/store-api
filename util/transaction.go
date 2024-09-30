package util

import "gorm.io/gorm"

func RecoverRollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
	}
}
