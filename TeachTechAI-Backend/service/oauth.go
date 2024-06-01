package service

import (
	"os"

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
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		clientID = "clientID"
	}
	return clientID
}

func getClientSecret() string {
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = "clientSecret"
	}
	return clientSecret
}

func getCallBackURL() string {
	callBackURL := os.Getenv("CALLBACK_URL")
	if callBackURL == "" {
		callBackURL = "callBackURL"
	}
	return callBackURL
}

func getGoogleSecretKey() string {
	googleSecretKey := os.Getenv("GOOGLE_SECRET")
	if googleSecretKey == "" {
		googleSecretKey = "googleSecretKey"
	}
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
