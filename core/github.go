package core

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/google/go-github/github"
)

type GitHub interface {
	CurrentUser() (string, error)
	ListRepos(user string) []*github.Repository
	SetupRepo(user string, repo string)
}

type githubImpl struct {
	client      *github.Client
	accessToken string
	config      *Config
}

func (g *githubImpl) CurrentUser() (string, error) {
	user, _, err := g.client.Users.Get(context.TODO(), "")
	if err != nil {
		log.Printf("client.Users.Get() faled: %v", err)
		return "", err
	}
	return *user.Login, nil
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
	config := make(map[string]interface{})
	config["url"] = "https://" + g.config.Hostname + "/hook"
	config["content_type"] = "json"
	web := "web"
	active := true
	// g.client.
	g.client.Repositories.CreateHook(context.TODO(), user, repo, &github.Hook{
		Name:   &web,
		Active: &active,
		Events: []string{"push", "pull_request"},
		Config: config,
	})
}

func Hook(r *http.Request) (string, interface{}, error) {
	hookType := github.WebHookType(r)
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read webhook: %v", err)
		return "", nil, err
	}
	p, err := github.ParseWebHook(hookType, b)
	if err != nil {
		log.Printf("Failed to read webhook: %v", err)
		return "", nil, err
	}
	return hookType, p, nil
}

func RepoName(hookType string, hook interface{}) (string, error) {
	switch hookType {
	case "PullRequestEvent":
		pullRequest := hook.(github.PullRequestEvent)
		return *pullRequest.Repo.Name, nil
	case "PushEvent":
		pushRequest := hook.(github.PullRequestEvent)
		return *pushRequest.Repo.Name, nil
	}
	return "", errors.New("unsupported event type")
}

func HookType(hook *interface{}) string {
	return reflect.TypeOf(hook).Name()
}

func GitHubSetup(config *Config, client *GithubClient) GitHub {
	var githubInstance GitHub
	var c *github.Client
	var t string
	if client == nil {
		c = nil
		t = ""
	} else {
		c = github.NewClient(client.HttpClient)
		t = client.Token
	}
	githubInstance = &githubImpl{
		client:      c,
		accessToken: t,
		config:      config,
	}
	return githubInstance
}

func ExtractRepositoriesNames(repos []*github.Repository) []string {
	m := func(r *github.Repository) string {
		s := strings.ReplaceAll(htmlRepo, "${name}", *r.Name)
		return s
	}
	return mapRepos(repos, m)
}

func mapRepos(repos []*github.Repository, m func(r *github.Repository) string) []string {
	var result []string
	for _, r := range repos {
		result = append(result, m(r))
	}
	return result
}

// Checkout example
//c := "git clone -u user:" + token.Token.AccessToken + " " + *pull.Repo.CloneURL + " ./test"
//log.Printf("Checkout: %s", c)
//exec.Command("git", "clone", "-u", "gimlet2:"+token.Token.AccessToken, *pull.Repo.CloneURL, "test").Start()
