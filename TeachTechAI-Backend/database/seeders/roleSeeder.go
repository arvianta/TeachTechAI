package seeders

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"teach-tech-ai/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ListRoleSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./database/seeders/data/roles.json")
	if err != nil {
		return err
	}

	jsonData, _ := io.ReadAll(jsonFile)

	var listRole []entity.Role
	if err := json.Unmarshal(jsonData, &listRole); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.Role{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.Role{}); err != nil {
			return err
		}
	}

	for _, data := range listRole {
		data.ID = uuid.New()
		var role entity.Role
		err := db.Where(&entity.Role{Name: data.Name}).First(&role).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		isData := db.Find(&role, "name = ?", data.Name).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
