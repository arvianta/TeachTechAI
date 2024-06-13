package controller

import (
	"fmt"
	"teach-tech-ai/dto"
	"teach-tech-ai/utils"

	"net/http"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	SendVerificationOTPByEmail(ctx *gin.Context)
	VerifyEmailWithOTP(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	MeUser(ctx *gin.Context)
	RefreshUser(ctx *gin.Context)
	GetAllUser(ctx *gin.Context)
	UpdateUserInfo(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	UploadUserProfilePicture(ctx *gin.Context)
	GetUserProfilePicture(ctx *gin.Context)
	DeleteUserProfilePicture(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type userController struct {
	jwtService  service.JWTService
	userService service.UserService
}

func NewUserController(us service.UserService, jwts service.JWTService) UserController {
	return &userController{
		userService: us,
		jwtService:  jwts,
	}
}

func (uc *userController) RegisterUser(ctx *gin.Context) {
	var user dto.UserCreateDTO
	err := ctx.ShouldBind(&user)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	checkUser, err := uc.userService.CheckUser(ctx.Request.Context(), user.Email)
	if checkUser {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	_, err = uc.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_REGISTER_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_REGISTER_USER, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) SendVerificationOTPByEmail(ctx *gin.Context) {
	var userVerifyDTO dto.SendUserOTPByEmail
	err := ctx.ShouldBind(&userVerifyDTO)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = uc.userService.SendUserOTPByEmail(ctx.Request.Context(), userVerifyDTO)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_SEND_OTP_EMAIL, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SEND_OTP_EMAIL_SUCCESS, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) VerifyEmailWithOTP(ctx *gin.Context) {
	var userVerifyDTO dto.VerifyUserOTPByEmail
	err := ctx.ShouldBind(&userVerifyDTO)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = uc.userService.VerifyUserOTPByEmail(ctx.Request.Context(), userVerifyDTO)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_VERIFY_EMAIL, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_VERIFY_EMAIL, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) LoginUser(ctx *gin.Context) {
	var userLoginDTO dto.UserLoginDTO
	err := ctx.ShouldBind(&userLoginDTO)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := uc.userService.Verify(ctx.Request.Context(), userLoginDTO.Email, userLoginDTO.Password)
	if !res {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	user, err := uc.userService.FindUserByEmail(ctx.Request.Context(), userLoginDTO.Email)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	roleID, err := uuid.Parse(user.RoleID)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	role, err := uc.userService.FindUserRoleByRoleID(roleID)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	sessionToken, refreshToken, atx, rtx, err := uc.jwtService.GenerateToken(user.ID, role)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}
	userResponse := dto.UserLoginResponseDTO{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
		Role:         role,
	}

	err = uc.userService.StoreUserToken(ctx.Request.Context(), user.ID, sessionToken, refreshToken, atx, rtx)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGIN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	response := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_LOGIN, userResponse)
	ctx.JSON(http.StatusOK, response)
}

func (uc *userController) MeUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	result, err := uc.userService.MeUser(ctx.Request.Context(), userID)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_GET_USER, result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) RefreshUser(ctx *gin.Context) {
	var refreshToken dto.UserRefreshDTO
	err := ctx.ShouldBind(&refreshToken)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newSessionToken, newRefreshToken, atx, rtx, err := uc.jwtService.RefreshToken(refreshToken.RefreshToken)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_REFRESHING_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	role, err := uc.jwtService.GetUserRoleByToken(newSessionToken)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_REFRESHING_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userID, err := uc.jwtService.GetUserIDByToken(newSessionToken)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_REFRESHING_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = uc.userService.StoreUserToken(ctx.Request.Context(), userID, newSessionToken, newRefreshToken, atx, rtx)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_REFRESHING_TOKEN, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_REFRESH_TOKEN, dto.UserLoginResponseDTO{
		SessionToken: newSessionToken,
		RefreshToken: newRefreshToken,
		Role:         role,
	})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) Logout(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	err := uc.jwtService.InvalidateToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGOUT, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_LOGOUT, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

// unused
func (uc *userController) GetAllUser(ctx *gin.Context) {
	result, err := uc.userService.GetAllUser(ctx.Request.Context())
	if err != nil {
		res := utils.BuildErrorResponse("Gagal Mendapatkan List User", err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildSuccessResponse("Berhasil Mendapatkan List User", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateUserInfo(ctx *gin.Context) {
	var user dto.UserUpdateInfoDTO
	err := ctx.ShouldBind(&user)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	user.ID = userID
	err = uc.userService.UpdateUser(ctx.Request.Context(), user)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_UPDATE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_UPDATE_USER, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) ChangePassword(ctx *gin.Context) {
	var user dto.UserChangePassword
	err := ctx.ShouldBind(&user)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	err = uc.userService.ChangePassword(ctx.Request.Context(), userID, user)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_CHANGE_PASSWORD, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_CHANGE_PASSWORD, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) ForgotPassword(ctx *gin.Context) {
	var forgotPassword dto.ForgotPassword
	err := ctx.ShouldBind(&forgotPassword)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = uc.userService.ForgotPassword(ctx.Request.Context(), forgotPassword)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_RESET_PASSWORD, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_RESET_PASSWORD, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) DeleteUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	err := uc.jwtService.InvalidateToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	err = uc.userService.DeleteUser(ctx.Request.Context(), userID)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_DELETE_USER, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_DELETE_USER, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) UploadUserProfilePicture(ctx *gin.Context) {
	var file dto.UploadFileDTO
	err := ctx.ShouldBind(&file)
	if err != nil {
		res := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), utils.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	localFileName := fmt.Sprintf("%s_%s", userID.String(), file.File.Filename)
	localFilePath := "/tmp/" + localFileName

	// Save the file locally
	if err := ctx.SaveUploadedFile(file.File, localFilePath); err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_UPLOAD_PROFILE_PICTURE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if err := uc.userService.UploadUserProfilePicture(ctx.Request.Context(), userID, localFilePath); err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_UPLOAD_PROFILE_PICTURE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_UPLOAD_PROFILE_PICTURE, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) GetUserProfilePicture(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	res, err := uc.userService.GetUserProfilePicture(ctx.Request.Context(), userID)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_GET_PROFILE_PICTURE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	ctx.File(res)
	_ = utils.DeleteTempFile(res)
}

func (uc *userController) DeleteUserProfilePicture(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_PROCESSING_REQUEST, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	err = uc.userService.DeleteUserProfilePicture(ctx.Request.Context(), userID)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_DELETE_PROFILE_PICTURE, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	res := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_DELETE_PROFILE_PICTURE, utils.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}
