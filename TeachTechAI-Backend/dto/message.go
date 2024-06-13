package dto

import (
	"errors"
	"teach-tech-ai/entity"
)

const (
	// Gagal
	MESSAGE_FAILED_CREATING_MESSAGE = "gagal membuat message"
	MESSAGE_FAILED_GET_MESSAGE      = "gagal mengambil message"
	MESSAGE_FAILED_GET_CONVO        = "gagal mengambil conversation"
	MESSAGE_FAILED_DELETE_CONVO     = "gagal menghapus conversation"

	// Berhasil
	MESSAGE_SUCCESS_CREATE_MESSAGE = "message berhasil dibuat"
	MESSAGE_SUCCESS_GET_MESSAGE    = "message berhasil diambil"
	MESSAGE_SUCCESS_GET_CONVO      = "conversation berhasil diambil"
	MESSAGE_SUCCESS_DELETE_CONVO   = "conversation berhasil dihapus"
)

var (
	ErrValidateUserConversation = errors.New("pengguna tidak diizinkan untuk mengakses conversation ini")
)

type (
	MessageRequestDTO struct {
		ConversationID string `json:"conversation_id" form:"conversation_id"`
		Topic          string `json:"topic" form:"topic" binding:"required"`
		Request        string `json:"request" form:"request" binding:"required"`
		AIModelName    string `json:"aimodel" form:"aimodel" binding:"required"`
	}

	MessageResponseDTO struct {
		ConversationID string `json:"conversation_id"`
		Message_ID     string `json:"message_id"`
		Response       string `json:"response"`
	}

	GetMessagesFromConversationDTO struct {
		ConversationID string `json:"conversation_id" form:"conversation_id"`
	}

	GetMessagesFromConversationResponseDTO struct {
		Messages []entity.Message `json:"messages"`
	}
)
