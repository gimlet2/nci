package core

import (
	"log"
	"net/http"
	"strings"
)

const htmlLogin = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>
`

const htmlMain = `<html><body>
Hello ${userName}
<div>
<a href="/repos">Repos</a>
</div>
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
<li>
<form method="POST" action="/repo/${name}/hook">
${name}
<input type="submit" value="Enable"/>
</form>
</li>
`

func SetupRouting(config *Config) {

	auth := AuthSetup(config)
	storage := StorageSetup()
	// var github GitHub
	ServerSetup(config, func(s Server) {
		s.Get("/", func(w http.ResponseWriter, r *http.Request) {
			s.Html(w, htmlLogin)
		})
		s.Get("/main", func(w http.ResponseWriter, r *http.Request) {
			user, _, _, err := currentUser(auth, r, config)
			if err != nil {
				goHome(w, r)
				return
			}
			s.Html(w, strings.ReplaceAll(htmlMain, "${userName}", user))
		})
		s.Get("/repos", func(w http.ResponseWriter, r *http.Request) {
			user, _, github, err := currentUser(auth, r, config)
			if err != nil {
				goHome(w, r)
				return
			}
			repos := github.ListRepos(user)
			if r.Header.Get("Content-type") == "application/json" {
				s.Json(w, ExtractRepositoriesNames(repos))
			} else {
				s.Html(w, strings.ReplaceAll(htmlReposList, "${repos}", strings.Join(ExtractRepositoriesNames(repos), "")))
			}
		})
		s.Post("/repo/{name}/hook", func(w http.ResponseWriter, r *http.Request) {
			user, token, github, err := currentUser(auth, r, config)
			if err != nil {
				goHome(w, r)
				return
			}
			repo := s.PathParam(r, "name")
			github.SetupRepo(user, repo)
			storage.Save(repo, token)
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
			user, err := github.CurrentUser()
			if err != nil {
				goHome(w, r)
				return
			}
			log.Printf("Logged in as GitHub user: %s\n", user)
			auth.SaveToken(token, w)
			http.Redirect(w, r, "/main", http.StatusTemporaryRedirect)
		})
		s.Post("/hook", func(w http.ResponseWriter, r *http.Request) {
			hookType, hook, err := Hook(r)
			if err != nil {
				return
			}
			repoName, err := RepoName(hookType, hook)
			if err != nil {
				return
			}
			token, err := storage.Read(repoName)
			GitHubSetup(config, auth.Client(token))

			s.Json(w, hook)
		})
	})
}

func currentUser(auth Auth, r *http.Request, config *Config) (string, *Token, GitHub, error) {
	token, err := auth.ReadToken(r)
	if err != nil {
		return "", nil, nil, err
	}
	github := GitHubSetup(config, auth.Client(token))
	user, err := github.CurrentUser()
	if err != nil {
		return "", nil, nil, err
	}
	return user, token, github, nil
}

func goHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
