package service

import (
	"context"
	"errors"
	"fmt"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"
	"teach-tech-ai/helpers"
	"teach-tech-ai/repository"
	"teach-tech-ai/utils"
	"time"

	"github.com/google/uuid"
	"github.com/mashingan/smapping"
)

type UserService interface {
	RegisterUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error)
	SendUserOTPByEmail(ctx context.Context, userVerifyDTO dto.SendUserOTPByEmail) error
	VerifyUserOTPByEmail(ctx context.Context, userVerifyDTO dto.VerifyUserOTPByEmail) error
	GetAllUser(ctx context.Context) ([]entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	Verify(ctx context.Context, email string, password string) (bool, error)
	CheckUser(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, userDTO dto.UserUpdateInfoDTO) error
	ChangePassword(ctx context.Context, userID uuid.UUID, passwordDTO dto.UserChangePassword) error
	ForgotPassword(ctx context.Context, forgotPasswordDTO dto.ForgotPassword) error
	MeUser(ctx context.Context, userID uuid.UUID) (entity.User, error)
	FindUserRoleByRoleID(roleID uuid.UUID) (string, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	StoreUserToken(ctx context.Context, userID uuid.UUID, sessionToken string, refreshToken string, atx time.Time, rtx time.Time) error
	UploadUserProfilePicture(ctx context.Context, userID uuid.UUID, localFilePath string) error
	GetUserProfilePicture(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteUserProfilePicture(ctx context.Context, userID uuid.UUID) error
}

type userService struct {
	userRepository  repository.UserRepository
	roleRepository  repository.RoleRepository
	otpEmailService OTPEmailService
}

func NewUserService(ur repository.UserRepository, rr repository.RoleRepository, os OTPEmailService) UserService {
	return &userService{
		userRepository:  ur,
		roleRepository:  rr,
		otpEmailService: os,
	}
}

func (us *userService) RegisterUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error) {
	roleID, err := us.roleRepository.FindRoleIDByName(ctx, "USER")
	if err != nil {
		return entity.User{}, err
	}

	user := entity.User{
		Email:      userDTO.Email,
		Name:       userDTO.Name,
		Password:   userDTO.Password,
		RoleID:     roleID,
		IsVerified: false,
	}

	createdUser, err := us.userRepository.RegisterUser(ctx, user)
	if err != nil {
		return user, err
	}

	_, err = us.otpEmailService.SendOTPByEmail(ctx, createdUser.Email)
	if err != nil {
		return createdUser, err
	}

	return createdUser, nil
}

func (us *userService) SendUserOTPByEmail(ctx context.Context, userVerifyDTO dto.SendUserOTPByEmail) error {
	user, err := us.userRepository.FindUserByEmail(ctx, userVerifyDTO.Email)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return errors.New("user already verified")
	}

	_, err = us.otpEmailService.SendOTPByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) VerifyUserOTPByEmail(ctx context.Context, userVerifyDTO dto.VerifyUserOTPByEmail) error {
	user, err := us.userRepository.FindUserByEmail(ctx, userVerifyDTO.Email)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return errors.New("user already verified")
	}

	err = us.otpEmailService.VerifyOTPByEmail(ctx, user.Email, userVerifyDTO.OTP)
	if err != nil {
		return err
	}

	user.IsVerified = true

	err = us.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
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
	if !res.IsVerified {
		return false, errors.New("user belum terverifikasi")
	}
	CheckPassword, err := helpers.CheckPassword(res.Password, []byte(password))
	if err != nil {
		return false, errors.New("email atau password salah")
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
	return true, dto.ErrEmailAlreadyExists
}

func (us *userService) UpdateUser(ctx context.Context, userDTO dto.UserUpdateInfoDTO) error {
	user := entity.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(userDTO))
	if err != nil {
		return err
	}
	return us.userRepository.UpdateUser(ctx, user)
}

func (us *userService) ChangePassword(ctx context.Context, userID uuid.UUID, passwordDTO dto.UserChangePassword) error {
	user, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = helpers.CheckPassword(user.Password, []byte(passwordDTO.OldPassword))
	if err != nil {
		return errors.New("password lama salah")
	}

	hashedNewPassword, err := helpers.HashPassword(passwordDTO.NewPassword)
	if err != nil {
		return err
	}

	_, err = helpers.CheckPassword(user.Password, []byte(passwordDTO.NewPassword))
	if err != nil {
		user.Password = string(hashedNewPassword)
		return us.userRepository.UpdateUser(ctx, user)
	}
	return errors.New("password baru tidak boleh sama dengan password lama")
}

func (us *userService) ForgotPassword(ctx context.Context, forgotPasswordDTO dto.ForgotPassword) error {
	user, err := us.userRepository.FindUserByEmail(ctx, forgotPasswordDTO.Email)
	if err != nil {
		return err
	}

	if !user.IsVerified {
		return errors.New("user not verified")
	}

	newPassword := helpers.GeneratePassword(16, true, true)
	hashedPassword, err := helpers.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = us.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	subject := "TeachTechAI OTP Verification"
	body := fmt.Sprintf("Your new password: %s", newPassword)

	err = utils.SendMail(forgotPasswordDTO.Email, subject, body)
	if err != nil {
		return err
	}

	return nil
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

func (us *userService) StoreUserToken(ctx context.Context, userID uuid.UUID, sessionToken string, refreshToken string, atx time.Time, rtx time.Time) error {
	return us.userRepository.StoreUserToken(ctx, userID, sessionToken, refreshToken, atx, rtx)
}

func (us *userService) UploadUserProfilePicture(ctx context.Context, userID uuid.UUID, localFilePath string) error {
	user, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.ProfilePicture != "" {
		err := utils.DeleteFileFromCloud(ctx, user.ProfilePicture)
		if err != nil {
			return err
		}
	}

	filename, err := utils.UploadFileToCloud(ctx, localFilePath, userID)
	if err != nil {
		return err
	}

	err = us.userRepository.UpdateProfilePicture(ctx, userID, filename)
	if err != nil {
		return err
	}

	_ = utils.DeleteTempFile(localFilePath)

	return nil
}

func (us *userService) GetUserProfilePicture(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	localFilePath := "/tmp/" + user.ProfilePicture

	if user.ProfilePicture == "" {
		return "", errors.New("profile picture not found")
	}

	err = utils.DownloadFileFromCloud(ctx, user.ProfilePicture, localFilePath)
	if err != nil {
		return "", err
	}

	return localFilePath, nil
}

func (us *userService) DeleteUserProfilePicture(ctx context.Context, userID uuid.UUID) error {
	user, err := us.userRepository.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.ProfilePicture == "" {
		return nil
	}

	err = utils.DeleteFileFromCloud(ctx, user.ProfilePicture)
	if err != nil {
		return err
	}

	err = us.userRepository.ClearProfilePicture(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
