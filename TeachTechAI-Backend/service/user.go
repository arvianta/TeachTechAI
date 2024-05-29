package service

import (
	"context"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"
	"teach-tech-ai/helpers"
	"teach-tech-ai/repository"
	"time"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type UserService interface {
	RegisterUser(ctx context.Context, userDTO dto.UserCreateDto) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	Verify(ctx context.Context, email string, password string) (bool, error)
	CheckUser(ctx context.Context, email string) (bool, error)
	FindUserRoleByRoleID(roleID uuid.UUID) (string, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) error
	MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error)
	StoreUserToken(userID uuid.UUID, sessionToken string, refreshToken string, atx time.Time, rtx time.Time) error
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
	user.RoleID, _ = us.roleRepository.FindRoleIDByName(ctx, "USER")
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

func (us *userService) Verify(ctx context.Context, email string, password string) (bool, error) {
	res, err := us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return false, err
	}
	CheckPassword, err := helpers.CheckPassword(res.Password, []byte(password))
	if err != nil {
		return false, err
	}
	if res.Email == email && CheckPassword {
		return true, nil
	}
	return false, nil
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

func (us *userService) UpdateUser(ctx context.Context, userDTO dto.UserUpdateDto) error {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(userDTO))
	if err != nil {
		return err
	}
	return us.userRepository.UpdateUser(ctx, user)
}

func (us *userService) MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	return us.userRepository.FindUserByID(ctx, userID)
}

func (us *userService) FindUserRoleByRoleID(roleID uuid.UUID) (string, error) {
	return us.roleRepository.FindRoleNameByID(roleID)
}

func (us *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return us.userRepository.DeleteUser(ctx, userID)
}

func (us *userService) StoreUserToken(userID uuid.UUID, sessionToken string, refreshToken string, atx time.Time, rtx time.Time) error {
	return us.userRepository.StoreUserToken(userID, sessionToken, refreshToken, atx, rtx)
}
