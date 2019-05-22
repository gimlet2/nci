package core


import (
	"fmt"
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

func SetupRouting() {

	auth := AuthSetup()
	var github GitHub
	ServerSetup(func(s Server) {
		s.Get("/", func(w http.ResponseWriter, r *http.Request) {
			s.Html(w, htmlIndex)
		})
		s.Get("/main", func(w http.ResponseWriter, r *http.Request) {
			if github == nil {

			}
			s.Html(w, strings.ReplaceAll(htmlMain, "${userName}", *github.CurrentUser().Login))
		})
		s.Get("/login", func(w http.ResponseWriter, r *http.Request) {
			url := auth.GetAuthUrl()
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		})
		s.Get("/github_oauth_cb", func(w http.ResponseWriter, r *http.Request) {
			token, err := auth.ExchangeForToken(r.FormValue("state"), r.FormValue("code"))
			if err != nil {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}

			github = GitHubSetup(auth.Client(token))
			user := github.CurrentUser()
			if user == nil {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
			http.Redirect(w, r, "/main", http.StatusTemporaryRedirect)
		})
		s.Post("/webhook", func(w http.ResponseWriter, r *http.Request) {
			s.Json(w, s.ReadJson(r))
		})
	})
}
