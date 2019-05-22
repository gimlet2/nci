package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"log"
	"net/http"
)

type GitHub interface {
	CurrentUser() *github.User
}

type githubImpl struct {
	client *github.Client
}

func (g *githubImpl) CurrentUser() *github.User {
	user, _, err := g.client.Users.Get(oauth2.NoContext, "")
	if err != nil {
		log.Printf("client.Users.Get() faled: %v", err)
		return nil
	}
	return user
}

func GitHubSetup(client *http.Client) GitHub {
	var githubInstanse GitHub
	githubInstanse = &githubImpl{
		client: github.NewClient(client),
	}
	return githubInstanse
}
