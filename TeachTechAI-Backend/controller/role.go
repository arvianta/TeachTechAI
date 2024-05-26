package controller

import (
	"teach-tech-ai/common"
	"teach-tech-ai/dto"

	// "teach-tech-ai/entity"
	"net/http"
	"teach-tech-ai/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleController interface {
	CreateRole(ctx *gin.Context)
	GetAllRole(ctx *gin.Context)
	GetRoleNameByID(ctx *gin.Context)
}

type roleController struct {
	userService service.UserService
	roleService service.RoleService
}

func NewRoleController(rs service.RoleService, us service.UserService) RoleController {
	return &roleController{
		roleService: rs,
		userService: us,
	}
}

func (rc *roleController) CreateRole(ctx *gin.Context) {
	var role dto.RoleCreateDto
	err := ctx.ShouldBind(&role)
	if err != nil {
		response := common.BuildErrorResponse("Gagal Login", err.Error(), common.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	
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

func (rc *roleController) GetRoleNameByID(ctx *gin.Context) {
	roleID := ctx.Param("role_id")
	roleUUID, err := uuid.Parse(roleID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan Role", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := rc.userService.FindUserRoleByRoleID(roleUUID)
	if err != nil {
		res := common.BuildErrorResponse("Gagal Mendapatkan Role", err.Error(), common.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	res := common.BuildResponse(true, "Berhasil Mendapatkan Role", result)
	ctx.JSON(http.StatusOK, res)
}