package entity

import (
	"teach-tech-ai/helpers"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
    ID             uuid.UUID `gorm:"type:char(36);primary_key;not_null" json:"id"`
    GoogleID       string    `gorm:"type:varchar(255);not_null" json:"google_id"`
    Email          string    `gorm:"type:varchar(255);unique;not_null" json:"email"`
    Name           string    `gorm:"type:varchar(255);not_null" json:"name"`
    Phone          string    `gorm:"type:varchar(15);not_null" json:"phone"`
    Password       string    `gorm:"type:varchar(255);not_null" json:"password"`
    ProfilePicture string    `gorm:"type:varchar(255);not_null" json:"profile_picture"`
    AsalInstansi   int       `gorm:"not_null" json:"asal_instansi"`
    DateOfBirth    time.Time `gorm:"type:date;not_null" json:"date_of_birth"`
    CreatedAt      time.Time `gorm:"type:timestamp;not_null" json:"created_at"`
    UpdatedAt      time.Time `gorm:"type:timestamp;not_null" json:"updated_at"`
    RoleID         uuid.UUID `gorm:"type:char(36);not_null" json:"role_id"`
    Role           Role      `gorm:"foreignkey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	var err error
	u.Password, err = helpers.HashPassword(u.Password)
	if err != nil {
		return err
	}
	return nil
}