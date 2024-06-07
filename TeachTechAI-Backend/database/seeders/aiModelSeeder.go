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

func ListAIModelSeeder(db *gorm.DB) error {
	jsonFile, err := os.Open("./database/seeders/data/aimodels.json")
	if err != nil {
		return err
	}

	jsonData, _ := io.ReadAll(jsonFile)

	var listAIModel []entity.AIModel
	if err := json.Unmarshal(jsonData, &listAIModel); err != nil {
		return err
	}

	hasTable := db.Migrator().HasTable(&entity.AIModel{})
	if !hasTable {
		if err := db.Migrator().CreateTable(&entity.AIModel{}); err != nil {
			return err
		}
	}

	for _, data := range listAIModel {
		data.ID = uuid.New()
		var aiModel entity.AIModel
		err := db.Where(&entity.AIModel{Name: data.Name}).First(&aiModel).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		isData := db.Find(&aiModel, "name = ?", data.Name).RowsAffected
		if isData == 0 {
			if err := db.Create(&data).Error; err != nil {
				return err
			}
		}
	}

	return nil
}