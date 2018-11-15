package main

import (
	// "crypto/sha256"

	"log"
	"net/http"

	"github.com/gorilla/context"
	"github.com/wisawing/GoBackend/handler"
)

// TODO: Handling concurrency?
// TODO: Security?

func main() {
	// fmt.Println("Hello")
	// http.HandleFunc("/foo", fooHandler)
	// http.HandleFunc("/", httpHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	mux := http.NewServeMux()
	mux.HandleFunc("/account", handler.AccountPageHandler)
	mux.HandleFunc("/signup", handler.SignupHandler)
	mux.HandleFunc("/register", handler.RegisterPageHandler) // maybe there is a better name for register vs signup?
	mux.HandleFunc("/login", handler.LoginHandler)
	mux.HandleFunc("/logout", handler.LogoutHandler)
	mux.HandleFunc("/", handler.MainPageHandler)

	// context.ClearHandler to prevent gorilla/session from mem leak.
	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(mux)))
}
