package core

import (
	"errors"
	"log"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type Auth interface {
	GetAuthUrl() string
	IsStateValid(state string) bool
	ExchangeForToken(state string, code string) (*oauth2.Token, error)
	Client(token *oauth2.Token) *http.Client
}

type authImpl struct {
	oauthConf  *oauth2.Config
	oauthState string
}

func (a *authImpl) GetAuthUrl() string {
	return a.oauthConf.AuthCodeURL(a.oauthState, oauth2.AccessTypeOffline)
}

func (a *authImpl) IsStateValid(state string) bool {
	if a.oauthState != state {
		log.Printf("invalid oauth state, expected '%s', got '%s'", a.oauthState, state)
		return false
	}
	return true
}

func (a *authImpl) ExchangeForToken(state string, code string) (*oauth2.Token, error) {
	if !a.IsStateValid(state) {
		return nil, errors.New("Invalide state")
	}

	token, err := a.oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Token exchane failed: %v", err)
		return nil, errors.New("Token exchane failed")
	}
	return token, nil
}

func (a *authImpl) Client(token *oauth2.Token) *http.Client {
	return a.oauthConf.Client(oauth2.NoContext, token)
}

func AuthSetup(config *Config) Auth {
	var auth Auth
	auth = &authImpl{
		oauthConf: &oauth2.Config{
			ClientID:     config.ClientId,
			ClientSecret: config.ClientSecret,
			// select level of access you want https://developer.github.com/v3/oauth/#scopes
			Scopes:   []string{"user:email", "repo"}, // more scopes to ad
			Endpoint: githuboauth.Endpoint,
		},
		oauthState: randomUuidString(),
	}
	return auth
}

func randomUuidString() string {
	s, err := uuid.NewV4()
	if err != nil {
		log.Printf("Failed to generate random Uuid: %v", err)
		return ""
	}
	return s.String()
}
