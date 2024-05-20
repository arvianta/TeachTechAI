package controller

import (
	"teach-tech-ai/common"
	"teach-tech-ai/dto"

	// "teach-tech-ai/entity"
	"net/http"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	RegisterUser(ctx *gin.Context)
	GetAllUser(ctx *gin.Context)
	// LoginUser(ctx *gin.Context)
	// DeleteUser(ctx *gin.Context)
	// UpdateUser(ctx *gin.Context)
	// MeUser(ctx *gin.Context)
}

type userController struct {
	// jwtService  service.JWTService
	userService service.UserService
}

// func NewUserController(us service.UserService, jwts service.JWTService) UserController {
// 	return &userController{
// 		userService: us,
// 		jwtService:  jwts,
// 	}
// }

func NewUserController(us service.UserService) UserController {
	return &userController{
		userService: us,
		// jwtService:  jwts,
	}
}

func (uc *userController) RegisterUser(ctx *gin.Context) {
	var user dto.UserCreateDto
	err := ctx.ShouldBind(&user)
	checkUser, _ := uc.userService.CheckUser(ctx.Request.Context(), user.Email)
	if checkUser {
		res := common.BuildErrorResponse("User Sudah Terdaftar", "false", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result, err := uc.userService.RegisterUser(ctx.Request.Context(), user)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan User", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Menambahkan User", result)
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
