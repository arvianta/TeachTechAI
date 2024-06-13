package middleware

import (
	"net/http"
	"strings"
	"teach-tech-ai/dto"
	"teach-tech-ai/service"
	"teach-tech-ai/utils"

	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_FOUND, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !strings.Contains(authHeader, "Bearer ") {
			response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !token.Valid {
			response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, dto.MESSAGE_FAILED_DENIED_ACCESS, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		isValid, err := jwtService.ValidateTokenWithDB(authHeader)
		if err != nil {
			response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !isValid {
			response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
		userID, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		ctx.Set("token", authHeader)
		ctx.Set("userID", userID)
		ctx.Next()
	}
}
