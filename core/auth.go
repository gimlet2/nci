package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type GithubClient struct {
	HttpClient *http.Client
	Token      string
}

type Token struct {
	Token *oauth2.Token
}

type Auth interface {
	GetAuthUrl() string
	IsStateValid(state string) bool
	ExchangeForToken(state string, code string) (*Token, error)
	Client(token *Token) *GithubClient
	ReadToken(r *http.Request) (*Token, error)
	SaveToken(t *Token, w http.ResponseWriter)
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

func (a *authImpl) ExchangeForToken(state string, code string) (*Token, error) {
	if !a.IsStateValid(state) {
		return nil, errors.New("invalid state")
	}

	token, err := a.oauthConf.Exchange(context.TODO(), code)
	if err != nil {
		log.Printf("Token exchane failed: %v", err)
		return nil, errors.New("token exchange failed")
	}
	return &Token{token}, nil
}

func (a *authImpl) Client(token *Token) *GithubClient {
	return &GithubClient{
		HttpClient: a.oauthConf.Client(context.TODO(), token.Token),
		Token:      token.Token.AccessToken,
	}
}

func (a *authImpl) ReadToken(r *http.Request) (*Token, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Printf("Faile to read session cookie: %v", err)
		return nil, err
	}
	decodedCookie, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(cookie.Value)
	if err != nil {
		log.Printf("Faile to read session cookie: %v", err)
		return nil, err
	}
	var t oauth2.Token
	err = json.Unmarshal(decodedCookie, &t)
	if err != nil {
		log.Printf("Faile to read session cookie: %v", err)
		return nil, err
	}
	return &Token{&t}, nil
}

func (a *authImpl) SaveToken(token *Token, w http.ResponseWriter) {
	j, _ := json.Marshal(token.Token)
	cookieValue := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(j)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    cookieValue,
		HttpOnly: true,
		Secure:   true,
		Expires:  token.Token.Expiry,
	})
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
