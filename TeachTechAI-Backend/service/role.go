package service

import (
	"context"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"

	"teach-tech-ai/repository"

	"github.com/mashingan/smapping"
)

type RoleService interface {
	CreateRole(ctx context.Context, roleDTO dto.RoleCreateDto) (entity.Role, error)
	GetAllRole(ctx context.Context) ([]entity.Role, error)
	CheckRole(ctx context.Context, name string) (bool, error)
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