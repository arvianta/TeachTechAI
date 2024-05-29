package repository

import (
	"context"
	"fmt"
	"teach-tech-ai/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(ctx context.Context, role entity.Role) (entity.Role, error)
	GetAllRole(ctx context.Context) ([]entity.Role, error)
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

func (db *roleConnection) CreateRole(ctx context.Context, role entity.Role) (entity.Role, error) {
	role.ID = uuid.New()
	uc := db.connection.Create(&role)
	if uc.Error != nil {
		return entity.Role{}, uc.Error
	}
	return role, nil
}

func (db *roleConnection) GetAllRole(ctx context.Context) ([]entity.Role, error) {
	var listRole []entity.Role
	tx := db.connection.Find(&listRole)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return listRole, nil
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