package database

import (
	"teach-tech-ai/database/seeders"

	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeders.ListRoleSeeder(db); err != nil {
		return err
	}

	return nil
}