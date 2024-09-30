package repository

import "store-api/internal/entity"

type RoleRepository struct {
	BaseRepository[entity.Role]
}

func NewRoleRepository() *RoleRepository {
	return new(RoleRepository)
}
