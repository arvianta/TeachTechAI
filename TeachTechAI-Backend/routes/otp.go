package routes

import (
	"teach-tech-ai/controller"

	"github.com/gin-gonic/gin"
)

func OTPRoutes(router *gin.Engine, OTPController controller.OTPTwilioController) {
	otpRoutes := router.Group("/api/otp")
	{
		otpRoutes.POST("/send", OTPController.SendSMS)
		otpRoutes.POST("/verify", OTPController.VerifySMS)
	}
}
