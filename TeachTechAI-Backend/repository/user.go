package repository

import (
	"context"
	"teach-tech-ai/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(ctx context.Context, user entity.User) (entity.User, error)
	GetAllUser(ctx context.Context) ([]entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (entity.User, error)
	FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, user entity.User) error
	StoreUserToken(userID uuid.UUID, sessionToken string, refreshToken string, atx time.Time, rtx time.Time) error
	FindUserRoleIDByID(userID uuid.UUID) (uuid.UUID, error)
	InvalidateUserToken(userID uuid.UUID) error
	GetUserSessionToken(userID uuid.UUID) (string, error)
	UpdateProfilePicture(ctx context.Context, userID uuid.UUID, url string) error
	ClearProfilePicture(ctx context.Context, userID uuid.UUID) error
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) RegisterUser(ctx context.Context, user entity.User) (entity.User, error) {
	uc := db.connection.WithContext(ctx).Create(&user)
	if uc.Error != nil {
		return entity.User{}, uc.Error
	}
	return user, nil
}

func (db *userConnection) GetAllUser(ctx context.Context) ([]entity.User, error) {
	var listUser []entity.User
	tx := db.connection.WithContext(ctx).Find(&listUser)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return listUser, nil
}

func (db *userConnection) FindUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var user entity.User
	ux := db.connection.WithContext(ctx).Where("email = ?", email).Take(&user)
	if ux.Error != nil {
		return user, ux.Error
	}
	return user, nil
}

func (db *userConnection) FindUserByID(ctx context.Context, userID uuid.UUID) (entity.User, error) {
	var user entity.User
	ux := db.connection.WithContext(ctx).Where("id = ?", userID).Take(&user)
	if ux.Error != nil {
		return user, ux.Error
	}
	return user, nil
}

func (db *userConnection) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	uc := db.connection.WithContext(ctx).Delete(&entity.User{}, &userID)
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}

func (db *userConnection) UpdateUser(ctx context.Context, user entity.User) error {
	uc := db.connection.WithContext(ctx).Updates(&user)
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}

func (db *userConnection) StoreUserToken(userID uuid.UUID, sessionToken string, refreshToken string, atx time.Time, rtx time.Time) error {
	user := entity.User{ID: userID}
	uc := db.connection.Model(&user).Updates(map[string]interface{}{
		"session_token": sessionToken,
		"refresh_token": refreshToken,
		"st_expires": atx,
		"rt_expires": rtx,
	})
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}

func (db *userConnection) FindUserRoleIDByID(userID uuid.UUID) (uuid.UUID, error) {
	var user entity.User
	ux := db.connection.Where("id = ?", userID).Take(&user)
	if ux.Error != nil {
		return uuid.Nil, ux.Error
	}
	roleID, err := uuid.Parse(user.RoleID)
	if err != nil {
		return uuid.Nil, err
	}
	return roleID, nil
}

func (db *userConnection) InvalidateUserToken(userID uuid.UUID) (error) {
	user := entity.User{ID: userID}
	uc := db.connection.Model(&user).Updates(map[string]interface{}{
		"session_token": "",
		"refresh_token": "",
		"st_expires": time.Time{},
		"rt_expires": time.Time{},
	})
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}

func (db *userConnection) GetUserSessionToken(userID uuid.UUID) (string, error) {
	var user entity.User
	ux := db.connection.Where("id = ?", userID).Take(&user)
	if ux.Error != nil {
		return "", ux.Error
	}
	return user.SessionToken, nil
}

func (db *userConnection) UpdateProfilePicture(ctx context.Context, userID uuid.UUID, url string) error {
	uc := db.connection.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("profile_picture", url)
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}

func (db *userConnection) ClearProfilePicture(ctx context.Context, userID uuid.UUID) error {
	uc := db.connection.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Update("profile_picture", "")
	if uc.Error != nil {
		return uc.Error
	}
	return nil
}

