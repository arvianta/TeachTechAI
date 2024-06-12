package repository

import (
	"fmt"
	"teach-tech-ai/entity"

	"gorm.io/gorm"
)

type AIModelRepository interface {
	FindAIModelIDByName(name string) (string, error)
}

type aimodelConnection struct {
	connection *gorm.DB
}

func NewAIModelRepository(db *gorm.DB) AIModelRepository {
	return &aimodelConnection{
		connection: db,
	}
}

func (db *aimodelConnection) FindAIModelIDByName(name string) (string, error) {
	var AIModel entity.AIModel
	var AIModelID string
	ux := db.connection.Select("id").Where("name = ?", name).Take(&AIModel).Scan(&AIModelID)

	if ux.Error != nil {
		return "", ux.Error
	}

	// Check if a row was actually found (optional)
	if ux.RowsAffected == 0 {
		return "", fmt.Errorf("AIModel with name '%s' not found", name)
	}

	return AIModelID, nil
}
