package controller

import (
	"fmt"
	"net/http"
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
}

func NewOAuthController(ots service.OAuthService) OAuthController {
	return &oauthController{
		oauthService: ots,
	}
}

// type contextKey string

func (oc *oauthController) GetAuthCallbackFunction(ctx *gin.Context) {
	// provider := ctx.Param("provider")
	// ctx.Request = ctx.Request.WithContext(context.WithValue(context.Background(), contextKey("provider"), provider))
	q := ctx.Request.URL.Query()
	q.Add("provider", "google")
	ctx.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
	if err != nil {
		response := utils.BuildErrorResponse("Gagal Login", "OAuth Error", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	fmt.Println(user)

	response := utils.BuildSuccessResponse("Berhasil Login", user)
	ctx.JSON(http.StatusOK, response)
}

func (oc *oauthController) Logout(ctx *gin.Context) {
	// provider := ctx.Param("provider")
	// ctx.Request = ctx.Request.WithContext(context.WithValue(context.Background(), contextKey("provider"), provider))
	q := ctx.Request.URL.Query()
	q.Add("provider", "google")
	ctx.Request.URL.RawQuery = q.Encode()

	gothic.Logout(ctx.Writer, ctx.Request)

	response := utils.BuildSuccessResponse("Berhasil Logout", nil)
	ctx.JSON(http.StatusOK, response)
}

func (oc *oauthController) Authenticate(ctx *gin.Context) {
	// provider := ctx.Param("provider")
	// ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), contextKey("provider"), provider))
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
