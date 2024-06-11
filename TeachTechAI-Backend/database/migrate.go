package database

import (
	"fmt"
	"teach-tech-ai/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		entity.Role{},
		entity.User{},
		entity.Conversation{},
		entity.AIModel{},
		entity.Message{},
		entity.OTP{},
	); err != nil {
		return err
	}
	fmt.Println("Migration success!")
	
	return nil
}