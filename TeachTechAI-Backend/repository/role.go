package repository

import (
	"context"
	"fmt"
	"teach-tech-ai/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	FindRoleByName(ctx context.Context, name string) (entity.Role, error)
	FindRoleIDByName(ctx context.Context, name string) (string, error)
	FindRoleNameByID(id uuid.UUID) (string, error)
}

type roleConnection struct {
	connection *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleConnection{
		connection: db,
	}
}

func (db *roleConnection) FindRoleByName(ctx context.Context, name string) (entity.Role, error) {
	var role entity.Role
	ux := db.connection.Where("name = ?", name).Take(&role)
	if ux.Error != nil {
		return role, ux.Error
	}
	return role, nil
}

func (db *roleConnection) FindRoleIDByName(ctx context.Context, name string) (string, error) {
	var role entity.Role
	var roleID string
	ux := db.connection.Select("id").Where("name = ?", name).Take(&role).Scan(&roleID)

	if ux.Error != nil {
		return "", ux.Error
	}

	// Check if a row was actually found (optional)
	if ux.RowsAffected == 0 {
		return "", fmt.Errorf("role with name '%s' not found", name)
	}

	return roleID, nil
}

func (db *roleConnection) FindRoleNameByID(id uuid.UUID) (string, error) {
	var role entity.Role
	var roleName string
	ux := db.connection.Select("name").Where("id = ?", id).Take(&role).Scan(&roleName)

	if ux.Error != nil {
		return "", ux.Error
	}

	if ux.RowsAffected == 0 {
		return "", fmt.Errorf("role with id '%s' not found", id)
	}

	return roleName, nil
}
