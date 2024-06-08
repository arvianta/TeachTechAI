package entity

import (
	"teach-tech-ai/helpers"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID             uuid.UUID `gorm:"type:char(36);primary_key;default:not_null" json:"id"`
	GoogleID       string    `gorm:"type:varchar(255);" json:"google_id"`
	Email          string    `gorm:"type:varchar(255);unique;not_null" json:"email"`
	Name           string    `gorm:"type:varchar(255);not_null" json:"name"`
	Phone          string    `gorm:"type:varchar(15);" json:"phone"`
	Password       string    `gorm:"type:varchar(255);" json:"password"`
	ProfilePicture string    `gorm:"type:varchar(255);" json:"profile_picture"`
	AsalInstansi   string    `gorm:"type:varchar(255);" json:"asal_instansi"`
	DateOfBirth    time.Time `gorm:"type:date;" json:"date_of_birth"`
	// IsVerified     bool      `gorm:"type:boolean;default:false" json:"is_verified"`
	SessionToken   string    `gorm:"type:varchar(255);" json:"session_token"`
	STExpires      time.Time `gorm:"type:timestamp;" json:"st_expires"`
	RefreshToken   string    `gorm:"type:varchar(255);" json:"refresh_token"`
	RTExpires      time.Time `gorm:"type:timestamp;" json:"rt_expires"`

	Timestamp

	RoleID string `gorm:"type:char(36);not_null" json:"role_id"`
	Role   Role   `gorm:"foreignkey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}
