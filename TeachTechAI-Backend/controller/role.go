package controller

import (
	"teach-tech-ai/common"
	"teach-tech-ai/dto"

	// "teach-tech-ai/entity"
	"net/http"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
)

type RoleController interface {
	CreateRole(ctx *gin.Context)
	GetAllRole(ctx *gin.Context)
	// LoginUser(ctx *gin.Context)
	// DeleteUser(ctx *gin.Context)
	// UpdateUser(ctx *gin.Context)
	// MeUser(ctx *gin.Context)
}

type roleController struct {
	// jwtService  service.JWTService
	roleService service.RoleService
}

func NewRoleController(rs service.RoleService) RoleController {
	return &roleController{
		roleService: rs,
		// jwtService:  jwts,
	}
}

func (rc *roleController) CreateRole(ctx *gin.Context) {
	var role dto.RoleCreateDto
	err := ctx.ShouldBind(&role)
	checkRole, _ := rc.roleService.CheckRole(ctx.Request.Context(), role.Name)
	if checkRole {
		res := common.BuildErrorResponse("Role Sudah Ada", "false", common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}
	result, err := rc.roleService.CreateRole(ctx.Request.Context(), role)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Menambahkan Role", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Menambahkan Role", result)
	ctx.JSON(http.StatusOK, res)
}

func (rc *roleController) GetAllRole(ctx *gin.Context) {
	result, err := rc.roleService.GetAllRole(ctx.Request.Context())
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan List Role", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan List Role", result)
	ctx.JSON(http.StatusOK, res)
}
