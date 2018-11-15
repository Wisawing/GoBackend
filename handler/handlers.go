package handler

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/wisawing/GoBackend/repos"
)

var sessionStore *sessions.CookieStore

const sessionName string = "session-name"

func init() {
	var key [32]byte
	rand.Read(key[:])
	sessionStore = sessions.NewCookieStore(key[:])
}

func CreateSession(user repos.User, res http.ResponseWriter, req *http.Request) {
	session, _ := sessionStore.Get(req, sessionName) // TODO: handle err
	session.Values["authed"] = true
	session.Values["user"] = user.Username
	session.Options = &sessions.Options{
		MaxAge:   100,
		HttpOnly: true,
	}
	session.Save(req, res)
}

func SignupHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	username, pwd, confirmPwd, email := "", "", "", ""

	username = req.FormValue("username")
	pwd = req.FormValue("password")
	confirmPwd = req.FormValue("confirm")
	email = req.FormValue("email")

	if len(username) == 0 ||
		len(pwd) == 0 ||
		len(confirmPwd) == 0 ||
		len(email) == 0 {
		fmt.Fprintf(res, "Error empty input")
		return
	}

	if pwd == confirmPwd {
		user := repos.Register(username, pwd, email)
		CreateSession(user, res, req)
		fmt.Fprintf(res, "Registration Completed")

		return
	} else {
		fmt.Fprintf(res, "Password and confirmation mismatch")
		return
	}
}

func LoginHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	username := req.FormValue("username")
	pwd := req.FormValue("password")

	user, ok := repos.ValidateUser(username, pwd)
	if ok {
		CreateSession(user, res, req)

		http.Redirect(res, req, "/account", http.StatusSeeOther)
	} else {
		// wrong password TODO: shows error
		http.Redirect(res, req, "/register", http.StatusFound)
	}
}

func LogoutHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	session, _ := sessionStore.Get(req, "session-name") // TODO: Handle error
	session.Values["authed"] = false

	session.Save(req, res)
}
