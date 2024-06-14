package service

import (
	"context"
	"fmt"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"
	"teach-tech-ai/helpers"
	"teach-tech-ai/repository"
	"teach-tech-ai/utils"
	"time"

	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/mashingan/smapping"
)

type UserService interface {
	RegisterUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error)
	SendUserOTPByEmail(ctx context.Context, email string) error
	VerifyUserOTPByEmail(ctx context.Context, userVerifyDTO dto.VerifyUserOTPByEmail) error
	LoginRegisterWithOAuth(ctx context.Context, user goth.User) (dto.UserLoginResponseDTO, error)
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
	jwtService      JWTService
}

func NewUserService(ur repository.UserRepository, rr repository.RoleRepository, os OTPEmailService, jwts JWTService) UserService {
	return &userService{
		userRepository:  ur,
		roleRepository:  rr,
		otpEmailService: os,
		jwtService:      jwts,
	}
}

func (us *userService) RegisterUser(ctx context.Context, userDTO dto.UserCreateDTO) (entity.User, error) {
	checkUser, err := us.CheckUser(ctx, userDTO.Email)
	if checkUser {
		if err == dto.ErrAccountAlreadyVerified {
			return entity.User{}, err
		}
		if err == dto.ErrAccountNotVerified {
			_, err := us.otpEmailService.SendOTPByEmail(ctx, userDTO.Email, userDTO.Name)
			if err != nil {
				return entity.User{}, err
			}
			return entity.User{}, dto.ErrAccountNotVerifiedWhenRegister
		}
	}

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

	_, err = us.otpEmailService.SendOTPByEmail(ctx, createdUser.Email, createdUser.Name)
	if err != nil {
		return createdUser, err
	}

	return createdUser, nil
}

func (us *userService) SendUserOTPByEmail(ctx context.Context, email string) error {
	user, err := us.userRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user.IsVerified {
		return dto.ErrAccountAlreadyVerified
	}

	_, err = us.otpEmailService.SendOTPByEmail(ctx, user.Email, user.Name)
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
		return dto.ErrAccountAlreadyVerified
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

func (us *userService) LoginRegisterWithOAuth(ctx context.Context, userGoth goth.User) (dto.UserLoginResponseDTO, error) {
	userData, err := us.userRepository.FindUserByEmail(ctx, userGoth.Email)
	if err != nil && err != dto.ErrUserNotFoundGorm {
		return dto.UserLoginResponseDTO{}, err
	} else if err == dto.ErrUserNotFoundGorm {
		// Register the user if not found
		roleID, err := us.roleRepository.FindRoleIDByName(ctx, "USER")
		if err != nil {
			return dto.UserLoginResponseDTO{}, err
		}

		userData = entity.User{
			Email:      userGoth.Email,
			GoogleID:   userGoth.UserID,
			RoleID:     roleID,
			IsVerified: true,
		}

		_, err = us.userRepository.RegisterUser(ctx, userData)
		if err != nil {
			return dto.UserLoginResponseDTO{}, err
		}
	} else {
		if userData.GoogleID == "" || !userData.IsVerified {
			userData.GoogleID = userGoth.UserID
			userData.IsVerified = true
			err = us.userRepository.UpdateUser(ctx, userData)
			if err != nil {
				return dto.UserLoginResponseDTO{}, err
			}
		}
	}

	// Handle Login
	roleID, err := uuid.Parse(userData.RoleID)
	if err != nil {
		return dto.UserLoginResponseDTO{}, err
	}

	role, err := us.roleRepository.FindRoleNameByID(roleID)
	if err != nil {
		return dto.UserLoginResponseDTO{}, err
	}

	sessionToken, refreshToken, atx, rtx, err := us.jwtService.GenerateToken(userData.ID, role)
	if err != nil {
		return dto.UserLoginResponseDTO{}, err
	}

	userResponse := dto.UserLoginResponseDTO{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
		Role:         role,
	}

	err = us.StoreUserToken(ctx, userData.ID, sessionToken, refreshToken, atx, rtx)
	if err != nil {
		return dto.UserLoginResponseDTO{}, err
	}

	return userResponse, nil
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
		return false, dto.ErrAccountNotVerified
	}
	CheckPassword, err := helpers.CheckPassword(res.Password, []byte(password))
	if err != nil {
		return false, dto.ErrEmailOrPassword
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

	if result.IsVerified {
		return true, dto.ErrAccountAlreadyVerified
	} else if !result.IsVerified {
		return true, dto.ErrAccountNotVerified
	}
	return false, nil
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
		return dto.ErrInvalidOldPassword
	}

	hashedNewPassword, err := helpers.HashPassword(passwordDTO.NewPassword)
	if err != nil {
		return err
	}

	_, err = helpers.CheckPassword(user.Password, []byte(passwordDTO.NewPassword))
	if err == nil {
		return dto.ErrPasswordSame
	}

	user.Password = string(hashedNewPassword)
	return us.userRepository.UpdateUser(ctx, user)
}

func (us *userService) ForgotPassword(ctx context.Context, forgotPasswordDTO dto.ForgotPassword) error {
	user, err := us.userRepository.FindUserByEmail(ctx, forgotPasswordDTO.Email)
	if err != nil {
		return err
	}

	if !user.IsVerified {
		return dto.ErrAccountNotVerified
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
		return "", dto.ErrProfilePictureNotFound
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
