package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wisawing/GoBackend/repos"
)

// MainPageHandler read main page
func MainPageHandler(res http.ResponseWriter, req *http.Request) {
	session, _ := sessionStore.Get(req, "session-name") // TODO: Handle error

	var body_b []byte
	var body string
	isAuthed, ok := session.Values["authed"]
	if ok && isAuthed.(bool) {
		http.Redirect(res, req, "/account", 302)
	} else {
		body_b, _ = ioutil.ReadFile("htmls/index.html")
		body = string(body_b[:])
	}
	fmt.Fprintf(res, string(body))
}

func AccountPageHandler(res http.ResponseWriter, req *http.Request) {
	session, _ := sessionStore.Get(req, "session-name") // TODO: Handle error

	if authed, ok := session.Values["authed"]; ok && authed.(bool) {
		var body, _ = ioutil.ReadFile("htmls/account.html")
		user, _ := repos.FindUser(session.Values["user"].(string)) // TODO: Handle wrong user
		body_s := string(body[:])
		body_s = fmt.Sprintf(body_s, user.Username, user.Email)

		fmt.Fprintf(res, body_s)
	} else {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
}

func RegisterPageHandler(res http.ResponseWriter, req *http.Request) {
	var body, _ = ioutil.ReadFile("htmls/signup.html")
	fmt.Fprintf(res, string(body))
}
