package core

import (
	"context"
	"log"

	"github.com/google/go-github/github"
	// "net/http"
)

type GitHub interface {
	CurrentUser() *github.User
	ListRepos(user string) []*github.Repository
	SetupRepo(user string, repo string)
}

type githubImpl struct {
	client      *github.Client
	accessToken string
	config      *Config
}

func (g *githubImpl) CurrentUser() *github.User {

	user, _, err := g.client.Users.Get(context.TODO(), "")
	if err != nil {
		log.Printf("client.Users.Get() faled: %v", err)
		return nil
	}
	return user
}

func (g *githubImpl) ListRepos(user string) []*github.Repository {
	repos, _, err := g.client.Repositories.List(context.TODO(), user, nil)
	if err != nil {
		log.Printf("client.Repositories.ListAll() faled: %v", err)
		return nil
	}
	return repos
}

func (g *githubImpl) SetupRepo(user string, repo string) {
	var config map[string]interface{}
	config["url"] = "https://" + g.config.Hostname + "/hook"
	config["content_type"] = "json"
	web := "web"
	active := true
	g.client.Repositories.CreateHook(context.TODO(), user, repo, &github.Hook{
		Name:   &web,
		Active: &active,
		Events: []string{"push", "pull_request"},
		Config: config,
	})
}

func GitHubSetup(config *Config, client *GithubClient) GitHub {
	var githubInstanse GitHub
	githubInstanse = &githubImpl{
		client:      github.NewClient(client.HttpClient),
		accessToken: client.Token,
		config:      config,
	}
	return githubInstanse
}
