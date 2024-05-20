package service

import (
	"context"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"

	// "teach-tech-ai/helpers"
	"teach-tech-ai/repository"

	"github.com/mashingan/smapping"
)

type RoleService interface {
	CreateRole(ctx context.Context, roleDTO dto.RoleCreateDto) (entity.Role, error)
	GetAllRole(ctx context.Context) ([]entity.Role, error)
	FindRoleByName(ctx context.Context, name string) (entity.Role, error)
	FindRoleIDByName(ctx context.Context, name string) (string, error)
	// Verify(ctx context.Context, email string, password string) (bool, error)
	CheckRole(ctx context.Context, name string) (bool, error)
	// DeleteUser(ctx context.Context, userID uuid.UUID) (error)
	// UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) (error)
	// MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error)
}

type roleService struct {
	roleRepository repository.RoleRepository
}

func NewRoleService(rr repository.RoleRepository) RoleService {
	return &roleService{
		roleRepository: rr,
	}
}

func (rs *roleService) CreateRole(ctx context.Context, roleDTO dto.RoleCreateDto) (entity.Role, error) {
	role := entity.Role{}
	err := smapping.FillStruct(&role, smapping.MapFields(roleDTO))
	if err != nil {
		return role, err
	}
	return rs.roleRepository.CreateRole(ctx, role)
}

func (rs *roleService) GetAllRole(ctx context.Context) ([]entity.Role, error) {
	return rs.roleRepository.GetAllRole(ctx)
}

func (rs *roleService) FindRoleByName(ctx context.Context, name string) (entity.Role, error) {
	return rs.roleRepository.FindRoleByName(ctx, name)
}

func (rs *roleService) FindRoleIDByName(ctx context.Context, name string) (string, error) {
	return rs.roleRepository.FindRoleIDByName(ctx, name)
}

func (rs *roleService) CheckRole(ctx context.Context, name string) (bool, error) {
	result, err := rs.roleRepository.FindRoleByName(ctx, name)
	if err != nil {
		return false, err
	}

	if result.Name == "" {
		return false, nil
	}
	return true, nil
}
