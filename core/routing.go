package core

import (
	"log"
	"net/http"
	"strings"

	githubLib "github.com/google/go-github/github"
)

const htmlIndex = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>
`

const htmlMain = `<html><body>
Hello ${userName}
</body></html>
`

const htmlReposList = `<html><body>
<h1>Repos</h1>
<ul>
${repos}
</ul>
</body></html>
`
const htmlRepo = `
<li>${name}</li>
`

func SetupRouting(config *Config) {

	auth := AuthSetup(config)
	// var github GitHub
	ServerSetup(config, func(s Server) {
		s.Get("/", func(w http.ResponseWriter, r *http.Request) {
			s.Html(w, htmlIndex)
		})
		s.Get("/main", func(w http.ResponseWriter, r *http.Request) {
			token := auth.ReadToken(r)
			if token == nil {
				goHome(w, r)
				return
			}
			github := GitHubSetup(config, auth.Client(token))
			s.Html(w, strings.ReplaceAll(htmlMain, "${userName}", *github.CurrentUser().Login))
		})
		s.Get("/repos", func(w http.ResponseWriter, r *http.Request) {
			token := auth.ReadToken(r)
			if token == nil {
				goHome(w, r)
				return
			}
			github := GitHubSetup(config, auth.Client(token))
			user := github.CurrentUser()
			if user == nil {
				goHome(w, r)
				return
			}
			repos := github.ListRepos(*user.Login)
			m := func(r *githubLib.Repository) string {
				s := strings.ReplaceAll(htmlRepo, "${name}", *r.Name)
				return s
			}
			s.Html(w, strings.ReplaceAll(htmlReposList, "${repos}", strings.Join(MapRepos(repos, m), "")))
		})
		s.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			url := auth.GetAuthUrl()
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		})
		s.Get("/github_oauth_cb", func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.ExchangeForToken(r.FormValue("state"), r.FormValue("code"))
			if err != nil {
				goHome(w, r)
				return
			}
			github := GitHubSetup(config, auth.Client(token))
			user := github.CurrentUser()
			if user == nil {
				goHome(w, r)
				return
			}
			log.Printf("Logged in as GitHub user: %s\n", *user.Login)
			auth.SaveToken(token, w)
			http.Redirect(w, r, "/main", http.StatusTemporaryRedirect)
		})
		s.Post("/hook", func(w http.ResponseWriter, r *http.Request) {
			s.Json(w, s.ReadJson(r))
		})
	})
}

func goHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func MapRepos(repos []*githubLib.Repository, m func(r *githubLib.Repository) string) []string {
	var result []string
	for _, r := range repos {
		result = append(result, m(r))
	}
	return result
}
