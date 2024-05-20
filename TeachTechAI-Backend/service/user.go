package service

import (
	"context"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"

	// "teach-tech-ai/helpers"
	"teach-tech-ai/repository"

	"github.com/mashingan/smapping"
)

type UserService interface {
	RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	// Verify(ctx context.Context, email string, password string) (bool, error)
	CheckUser(ctx context.Context, email string) (bool, error)
	// DeleteUser(ctx context.Context, userID uuid.UUID) (error)
	// UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) (error)
	// MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
	roleRepository repository.RoleRepository
}

func NewUserService(ur repository.UserRepository, rr repository.RoleRepository) UserService {
	return &userService{
		userRepository: ur,
		roleRepository: rr,
	}
}

func (us *userService) RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error) {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(userDTO))
	user.RoleID, _ = us.FindRoleIDByName(ctx, "user") //TODO
	if err != nil {
		return user, err
	}
	return us.userRepository.RegisterUser(ctx, user)
}

func (us *userService) GetAllUser(ctx context.Context) ([]entity.User, error) {
	return us.userRepository.GetAllUser(ctx)
}

func (us *userService) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	return us.userRepository.FindUserByEmail(ctx, email)
}

func (us *userService) FindRoleIDByName(ctx context.Context, name string) (string, error) {
	return us.userRepository.FindRoleIDByName(ctx, name)
}

func (us *userService) CheckUser(ctx context.Context, email string) (bool, error) {
	result, err := us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	if result.Email == "" {
		return false, nil
	}
	return true, nil
}
