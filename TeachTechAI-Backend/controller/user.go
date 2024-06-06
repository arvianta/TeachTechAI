package controller

import (
	"teach-tech-ai/common"
	"teach-tech-ai/dto"

	"net/http"
	"teach-tech-ai/entity"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	MeUser(ctx *gin.Context)
	RefreshUser(ctx *gin.Context)
	GetAllUser(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
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
	var user dto.UserCreateDto
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
	//

	res := common.BuildResponse(true, "Berhasil Menambahkan User", common.EmptyObj{})
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

	res, _ := uc.userService.Verify(ctx.Request.Context(), userLoginDTO.Email, userLoginDTO.Password)
	if !res {
		response := common.BuildErrorResponse("Gagal Login", "Email atau Password Salah", common.EmptyObj{})
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
	userResponse := entity.Authorization{
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

	res := common.BuildResponse(true, "Berhasil Refresh Token", entity.Authorization{
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

func (uc *userController) UpdateUser(ctx *gin.Context) {
	var user dto.UserUpdateDto
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