package controller

import (
	"net/http"
	"teach-tech-ai/dto"
	"teach-tech-ai/service"
	"teach-tech-ai/utils"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
)

type OAuthController interface {
	GetAuthCallbackFunction(ctx *gin.Context)
	Logout(ctx *gin.Context)
	Authenticate(ctx *gin.Context)
}

type oauthController struct {
	oauthService service.OAuthService
	userService  service.UserService
}

func NewOAuthController(ots service.OAuthService, us service.UserService) OAuthController {
	return &oauthController{
		oauthService: ots,
		userService:  us,
	}
}

// type contextKey string

func (oc *oauthController) GetAuthCallbackFunction(ctx *gin.Context) {
	q := ctx.Request.URL.Query()
	q.Add("provider", "google")
	ctx.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGIN, err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	res, err := oc.userService.LoginRegisterWithOAuth(ctx, user)
	if err != nil {
		response := utils.BuildErrorResponse(dto.MESSAGE_FAILED_LOGIN, err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := utils.BuildSuccessResponse(dto.MESSAGE_SUCCESS_LOGIN, res)
	ctx.JSON(http.StatusOK, response)
}

func (oc *oauthController) Logout(ctx *gin.Context) {
	q := ctx.Request.URL.Query()
	q.Add("provider", "google")
	ctx.Request.URL.RawQuery = q.Encode()

	gothic.Logout(ctx.Writer, ctx.Request)

	response := utils.BuildSuccessResponse("Berhasil Logout", nil)
	ctx.JSON(http.StatusOK, response)
}

func (oc *oauthController) Authenticate(ctx *gin.Context) {
	q := ctx.Request.URL.Query()
	q.Add("provider", "google")
	ctx.Request.URL.RawQuery = q.Encode()

	if gothUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request); err == nil {
		response := utils.BuildSuccessResponse("Authentication successful", gothUser)
		ctx.JSON(http.StatusOK, response)
	} else {
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	}
}
