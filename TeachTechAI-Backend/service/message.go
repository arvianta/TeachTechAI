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
	CreateMessage(ctx context.Context, msgDTO dto.MessageRequestDTO) (dto.MessageResponseDTO, error)
	GetMessagesFromConversation(ctx context.Context, msgDTO dto.GetMessagesFromConversationDTO) (dto.GetMessagesFromConversationResponseDTO, error)
}

type messageService struct {
	messageRepository      repository.MessageRepository
	conversationRepository repository.ConversationRepository
	aimodelRepository      repository.AIModelRepository
}

func NewMessageService(mr repository.MessageRepository, cr repository.ConversationRepository, air repository.AIModelRepository) MessageService {
	return &messageService{
		messageRepository:      mr,
		conversationRepository: cr,
		aimodelRepository:      air,
	}
}

func (ms *messageService) CreateMessage(ctx context.Context, msgDTO dto.MessageRequestDTO) (dto.MessageResponseDTO, error) {
	message := entity.Message{}
	var err error
	//AI Model ID
	aimodelID, err := ms.aimodelRepository.FindAIModelIDByName(msgDTO.AIModelName)
	if err != nil {
		return dto.MessageResponseDTO{}, err
	}
	message.AIModelID, err = uuid.Parse(aimodelID)
	if err != nil {
		return dto.MessageResponseDTO{}, err
	}
	// Conversation ID
	message.ConversationID, err = uuid.Parse(msgDTO.ConversationID)
	if err != nil {
		return dto.MessageResponseDTO{}, err
	}
	// Message Request
	message.Request = msgDTO.Request
	response, err := utils.PromptAI(msgDTO.Request, msgDTO.AIModelName)
	if err != nil {
		return dto.MessageResponseDTO{}, err
	}
	// Response and Datetime
	message.Response = response.Choices[0].Message.Content
	message.Datetime = time.Now()
	message.NumOfTokens = response.Usage.CompletionTokens
	message.FinishReason = response.Choices[0].FinishReason
	// Store Message to DB
	message, err = ms.messageRepository.StoreMessage(message)
	if err != nil {
		return dto.MessageResponseDTO{}, err
	}

	res := dto.MessageResponseDTO{
		ConversationID: message.ConversationID.String(),
		Message_ID:     message.ID.String(),
		Response:       message.Response,
	}

	return res, nil
}

func (ms *messageService) GetMessagesFromConversation(ctx context.Context, msgDTO dto.GetMessagesFromConversationDTO) (dto.GetMessagesFromConversationResponseDTO, error) {
	convID, err := uuid.Parse(msgDTO.ConversationID)
	if err != nil {
		return dto.GetMessagesFromConversationResponseDTO{}, err
	}

	var responseMessages dto.GetMessagesFromConversationResponseDTO

	responseMessages.Messages, err = ms.messageRepository.GetMessagesFromConversation(ctx, convID)
	if err != nil {
		return dto.GetMessagesFromConversationResponseDTO{}, err
	}

	return responseMessages, nil
}
