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
	// CreateMessageStream(ctx context.Context, msgDTO dto.MessageRequestDTO) (chan string, error)
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
	aimodelID, err := ms.aimodelRepository.FindAIModelIDByName(ctx, msgDTO.AIModelName)
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
	message, err = ms.messageRepository.StoreMessage(ctx, message)
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

// func (ms *messageService) CreateMessageStream(ctx context.Context, msgDTO dto.MessageRequestDTO) (chan string, error) {
// 	responseChan := make(chan string)

// 	var (
// 		completeMessage string
// 		numOfTokens     int
// 		finishReason    string
// 	)

// 	// Call utility function to start streaming AI responses
// 	go func() {
// 		defer close(responseChan)

// 		msgContent, tokens, reason, err := utils.PromptAIStream(ctx, msgDTO.Request, msgDTO.AIModelName, responseChan)
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		completeMessage = msgContent
// 		numOfTokens = tokens
// 		finishReason = reason

// 		log.Println(completeMessage)
// 		log.Println(numOfTokens)
// 		log.Println(finishReason)
// 	}()

// 	// Store Message to DB (if needed)

// 	// This part can be removed if you're streaming directly to the client
// 	// Create message entity
// 	// message := entity.Message{
// 	// 	Request:        msgDTO.Request,
// 	// 	Response:       completeMessage,
// 	// 	Datetime:       time.Now(),
// 	// 	NumOfTokens:    numOfTokens,
// 	// 	FinishReason:   finishReason,
// 	// }

// 	// Store Message to DB
// 	// _, err := ms.messageRepository.StoreMessage(message)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	return responseChan, nil
// }

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
