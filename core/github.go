package core

import (
	"github.com/google/go-github/github"
	// "golang.org/x/oauth2"
	"context"

	"log"
	"net/http"
)

type GitHub interface {
	CurrentUser() *github.User
	ListRepos() []*github.Repository
}

type githubImpl struct {
	client *github.Client
}

func (g *githubImpl) CurrentUser() *github.User {
	user, _, err := g.client.Users.Get(context.TODO(), "")
	if err != nil {
		log.Printf("client.Users.Get() faled: %v", err)
		return nil
	}
	return user
}

func (g *githubImpl) ListRepos() []*github.Repository {
	repos, _, err := g.client.Repositories.ListAll(context.TODO(), nil)
	if err != nil {
		log.Printf("client.Repositories.ListAll() faled: %v", err)
		return nil
	}
	return repos
}

func GitHubSetup(client *http.Client) GitHub {
	var githubInstanse GitHub
	githubInstanse = &githubImpl{
		client: github.NewClient(client),
	}
	return githubInstanse
}
