package repository

import (
	"store-api/internal/entity"

	"gorm.io/gorm"
)

type UserRepository struct {
	BaseRepository[entity.User]
}

func NewUserRepository() *UserRepository {
	return new(UserRepository)
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("lower(email) = lower(?)", email).
		First(user).Error
}

func (r *UserRepository) EmailExists(db *gorm.DB, email string) (bool, error) {
	var count int

	err := db.Model(new(entity.User)).
		Select("1").Where("lower(email) = lower(%s)", email).
		Limit(1).
		Scan(&count).
		Error

	return count == 1, err
}

func (r *UserRepository) AssignRole(db *gorm.DB, user *entity.User, role *entity.Role) error {
	return db.Model(user).Association("Roles").Append(role)
}

func (r *UserRepository) FindByRoles(db *gorm.DB, roleID string) ([]entity.User, error) {
	var users []entity.User

	err := db.Model(&entity.Role{}).
		Where("roles.id = ?", roleID).
		Association("Users").Find(&users)

	return users, err
}

func (r *UserRepository) HasRole(db *gorm.DB, userID string, roleID string) (bool, error) {
	var count int

	err := db.Table("user_roles").
		Select("1").
		Where("user_id = ?", userID).
		Where("role_id = ?", roleID).
		Limit(1).
		Scan(&count).
		Error

	return count == 1, err
}
