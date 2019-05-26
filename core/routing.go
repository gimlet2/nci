package core

import (
	"log"
	"net/http"
	"strings"
)

const htmlIndex = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>
`

const htmlMain = `<html><body>
Hello ${userName}
</body></html>
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
			github := GitHubSetup(auth.Client(token))
			s.Html(w, strings.ReplaceAll(htmlMain, "${userName}", *github.CurrentUser().Login))
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
			github := GitHubSetup(auth.Client(token))
			user := github.CurrentUser()
			if user == nil {
				goHome(w, r)
				return
			}
			log.Printf("Logged in as GitHub user: %s\n", *user.Login)
			auth.SaveToken(token, w)
			http.Redirect(w, r, "/main", http.StatusTemporaryRedirect)
		})
		s.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
			s.Json(w, s.ReadJson(r))
		})
	})
}

func goHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
