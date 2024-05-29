package middleware

import (
	"net/http"
	"strings"
	"teach-tech-ai/common"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService service.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Ditemukan", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !strings.Contains(authHeader, "Bearer ") {
			response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !token.Valid {
			response := common.BuildErrorResponse("Gagal Memproses Request", "Akses Ditolak", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		isValid, err := jwtService.ValidateTokenWithDB(authHeader)
		if err != nil {
			response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !isValid {
			response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
		userID, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			response := common.BuildErrorResponse("Gagal Memproses Request", err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		

		ctx.Set("token", authHeader)
		ctx.Set("userID", userID)
		ctx.Next()
	}
}