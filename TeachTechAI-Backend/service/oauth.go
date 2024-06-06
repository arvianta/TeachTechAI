package service

import (
	"teach-tech-ai/utils"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type OAuthService interface {
	InitOAuth() error
}

type oauthService struct {
	googleClientID     string
	googleClientSecret string
	googleCallBackURL  string
	secretKey 	 string
	maxAge 		 int	
}

func NewOAuthService() OAuthService {
	return &oauthService{
		googleClientID: getClientID(),
		googleClientSecret: getClientSecret(),
		googleCallBackURL: getCallBackURL(),
		secretKey: getGoogleSecretKey(),
		maxAge: 86400 * 30,
	}
}

func getClientID() string {
	clientID := utils.MustGetenv("GOOGLE_CLIENT_ID")
	return clientID
}

func getClientSecret() string {
	clientSecret := utils.MustGetenv("GOOGLE_CLIENT_SECRET")
	return clientSecret
}

func getCallBackURL() string {
	callBackURL := utils.MustGetenv("CALLBACK_URL")
	return callBackURL
}

func getGoogleSecretKey() string {
	googleSecretKey := utils.MustGetenv("GOOGLE_SECRET")
	return googleSecretKey
}

func (o *oauthService) InitOAuth() error {
	store := sessions.NewCookieStore([]byte(o.secretKey))
	store.MaxAge(o.maxAge)
	
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false

	gothic.Store = store
	
	goth.UseProviders(
		google.New(
			o.googleClientID,
			o.googleClientSecret,
			o.googleCallBackURL,
		),
	)

	return nil
}
