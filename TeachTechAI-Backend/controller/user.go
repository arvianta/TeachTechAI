package controller

import (
	"fmt"
	"teach-tech-ai/common"
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
		response := common.BuildErrorResponse("Gagal Register", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	checkUser, _ := uc.userService.CheckUser(ctx.Request.Context(), user.Email)
	if checkUser {
		res := common.BuildErrorResponse("User Sudah Terdaftar", "false", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	_, err = uc.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Menambahkan User", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) SendVerificationOTPByEmail(ctx *gin.Context) {
	var userVerifyDTO dto.SendUserOTPByEmail
	err := ctx.ShouldBind(&userVerifyDTO)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Mengirim Email Verifikasi", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = uc.userService.SendUserOTPByEmail(ctx.Request.Context(), userVerifyDTO)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Mengirim Email Verifikasi", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mengirim Email Verifikasi", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) VerifyEmailWithOTP(ctx *gin.Context) {
	var userVerifyDTO dto.VerifyUserOTPByEmail
	err := ctx.ShouldBind(&userVerifyDTO)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Verifikasi Email", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = uc.userService.VerifyUserOTPByEmail(ctx.Request.Context(), userVerifyDTO)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Verifikasi Email", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := common.BuildResponse(true, "Berhasil Verifikasi Email", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) LoginUser(ctx *gin.Context) {
	var userLoginDTO dto.UserLoginDTO
	err := ctx.ShouldBind(&userLoginDTO)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := uc.userService.Verify(ctx.Request.Context(), userLoginDTO.Email, userLoginDTO.Password)
	if !res {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	
	user, err := uc.userService.FindUserByEmail(ctx.Request.Context(), userLoginDTO.Email)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	roleID, err := uuid.Parse(user.RoleID)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}
	
	role, err := uc.userService.FindUserRoleByRoleID(roleID)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	sessionToken, refreshToken, atx, rtx, err := uc.jwtService.GenerateToken(user.ID, role)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}
	userResponse := dto.UserLoginResponseDTO{
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
		Role: role,
	}

	err = uc.userService.StoreUserToken(user.ID, sessionToken, refreshToken, atx, rtx)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}
	
	response := common.BuildResponse(true, "Berhasil Login", userResponse)
	ctx.JSON(http.StatusOK, response)
}

func (uc *userController) MeUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	result, err := uc.userService.MeUser(ctx.Request.Context(), userID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan User", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) RefreshUser(ctx *gin.Context) {
	var refreshToken dto.UserRefreshDTO
	err := ctx.ShouldBind(&refreshToken)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Refresh Token", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	newSessionToken, newRefreshToken, atx, rtx, err := uc.jwtService.RefreshToken(refreshToken.RefreshToken)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Refresh Token", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	role, err := uc.jwtService.GetUserRoleByToken(newSessionToken)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Refresh Token", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userID, err := uc.jwtService.GetUserIDByToken(newSessionToken)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Refresh Token", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	err = uc.userService.StoreUserToken(userID, newSessionToken, newRefreshToken, atx, rtx)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Refresh Token", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
	}

	res := common.BuildResponse(true, "Berhasil Refresh Token", dto.UserLoginResponseDTO{
		SessionToken: newSessionToken,
		RefreshToken: newRefreshToken,
		Role: role,
	})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) Logout(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	err := uc.jwtService.InvalidateToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Logout", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := common.BuildResponse(true, "Berhasil Logout", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) GetAllUser(ctx *gin.Context) {
	result, err := uc.userService.GetAllUser(ctx.Request.Context())
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan List User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan List User", result)
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) UpdateUserInfo(ctx *gin.Context) {
	var user dto.UserUpdateInfoDTO
	err := ctx.ShouldBind(&user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	user.ID = userID
	err = uc.userService.UpdateUser(ctx.Request.Context(), user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Mengupdate User", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) ChangePassword(ctx *gin.Context) {
	var user dto.UserChangePassword
	err := ctx.ShouldBind(&user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate Password", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	err = uc.userService.ChangePassword(ctx.Request.Context(), userID, user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate Password", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mengupdate Password", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) ForgotPassword(ctx *gin.Context) {
	var forgotPassword dto.ForgotPassword
	err := ctx.ShouldBind(&forgotPassword)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Reset Password", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	err = uc.userService.ForgotPassword(ctx.Request.Context(), forgotPassword)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Reset Password", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mengirim Password Baru", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) DeleteUser(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	err := uc.jwtService.InvalidateToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Logout", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	userID, err := uc.jwtService.GetUserIDByToken(token)

	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	err = uc.userService.DeleteUser(ctx.Request.Context(), userID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menghapus User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	res := common.BuildResponse(true, "Berhasil Menghapus User", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) UploadUserProfilePicture(ctx *gin.Context) {
	var file dto.UploadFileDTO
	err := ctx.ShouldBind(&file)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mengupdate User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	localFileName := fmt.Sprintf("%s_%s", userID.String(), file.File.Filename)
	localFilePath := "/tmp/" + localFileName

	// Save the file locally
	if err := ctx.SaveUploadedFile(file.File, localFilePath); err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Gagal Menyimpan File", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if err := uc.userService.UploadUserProfilePicture(ctx.Request.Context(), userID, localFilePath); err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Gagal Upload File to Cloud", nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res := common.BuildResponse(true, "Berhasil mengupload file", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}

func (uc *userController) GetUserProfilePicture(ctx *gin.Context) {
	token := ctx.MustGet("token").(string)
	userID, err := uc.jwtService.GetUserIDByToken(token)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	res, err := uc.userService.GetUserProfilePicture(ctx.Request.Context(), userID)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Gagal Download File from Cloud", nil)
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
		response := common.BuildErrorResponse("Gagal Memproses Request", "Token Tidak Valid", nil)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}

	err = uc.userService.DeleteUserProfilePicture(ctx.Request.Context(), userID)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Memproses Request", "Gagal Menghapus File dari Cloud", nil)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	res := common.BuildResponse(true, "Berhasil menghapus foto profil", common.EmptyObj{})
	ctx.JSON(http.StatusOK, res)
}