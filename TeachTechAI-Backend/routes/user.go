package routes

import (
	"teach-tech-ai/controller"
	"teach-tech-ai/middleware"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, UserController controller.UserController, jwtService service.JWTService) {
	userRoutes := router.Group("/api/user")
	{
		userRoutes.POST("/register", UserController.RegisterUser)
		userRoutes.POST("/send-otp", UserController.SendVerificationOTPByEmail)
		userRoutes.POST("/verify-otp", UserController.VerifyEmailWithOTP)
		userRoutes.POST("/login", UserController.LoginUser)
		userRoutes.GET("/me", middleware.Authenticate(jwtService), UserController.MeUser)
		userRoutes.POST("/refresh", UserController.RefreshUser)
		userRoutes.PATCH("/update", middleware.Authenticate(jwtService), UserController.UpdateUserInfo)
		userRoutes.PATCH("/change-password", middleware.Authenticate(jwtService), UserController.ChangePassword)
		userRoutes.POST("/forgot-password", UserController.ForgotPassword)
		userRoutes.GET("/profile-picture", middleware.Authenticate(jwtService), UserController.GetUserProfilePicture)
		userRoutes.POST("/upload-profile-picture", middleware.Authenticate(jwtService), UserController.UploadUserProfilePicture)
		userRoutes.POST("/delete-profile-picture", middleware.Authenticate(jwtService), UserController.DeleteUserProfilePicture)
		userRoutes.POST("/logout", middleware.Authenticate(jwtService), UserController.Logout)
		userRoutes.DELETE("/delete", middleware.Authenticate(jwtService), UserController.DeleteUser)
	}
}
