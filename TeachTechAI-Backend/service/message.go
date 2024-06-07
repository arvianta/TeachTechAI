package service

import (
	"context"
	"teach-tech-ai/dto"
	"teach-tech-ai/entity"
	"teach-tech-ai/repository"
	"teach-tech-ai/utils"
	"time"

	"github.com/google/uuid"
)

type MessageService interface {
	CreateMessage(ctx context.Context, msgDTO dto.MessageRequest) (dto.MessageResponse, error)
}

type messageService struct {
	messageRepository 		repository.MessageRepository
	conversationRepository 	repository.ConversationRepository
	aimodelRepository 		repository.AIModelRepository
}

func NewMessageService(mr repository.MessageRepository, cr repository.ConversationRepository, air repository.AIModelRepository) MessageService {
	return &messageService{
		messageRepository: mr,
		conversationRepository: cr,
		aimodelRepository: air,
	}
}

func (ms *messageService) CreateMessage(ctx context.Context, msgDTO dto.MessageRequest) (dto.MessageResponse, error) {
	message := entity.Message{}
	var err error
	//AI Model ID
	aimodelID, err := ms.aimodelRepository.FindAIModelIDByName(msgDTO.AIModelName)
	if err != nil {
		return dto.MessageResponse{}, err
	}
	message.AIModelID, err = uuid.Parse(aimodelID)
	if err != nil {
		return dto.MessageResponse{}, err
	}
	// Conversation ID
	message.ConversationID, err = uuid.Parse(msgDTO.ConversationID)
	if err != nil {
		return dto.MessageResponse{}, err
	}
	// Message Request
	message.Request = msgDTO.Request
	response, err := utils.PromptAI(msgDTO.Request)
	if err != nil {
		return dto.MessageResponse{}, err
	}
	// Response and Datetime
	message.Response = response.GeneratedText
	message.Datetime = time.Now()
	message.NumOfTokens = response.Details.GeneratedTokens
	// Store Message to DB
	message, err = ms.messageRepository.StoreMessage(message)
	if err != nil {
		return dto.MessageResponse{}, err
	}

	res := dto.MessageResponse{
		ConversationID: message.ConversationID.String(),
		Message_ID: message.ID.String(),
		Response: message.Response,
	}

	return res, nil
}