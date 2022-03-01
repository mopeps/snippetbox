package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := mux.NewRouter()
	mux.HandleFunc("/", app.home).Methods("GET")
	mux.HandleFunc("/snippet/create", app.createSnippet).Methods("POST")
	mux.HandleFunc("/snippet/create", app.createSnippetForm).Methods("GET")
	mux.HandleFunc("/snippet/{id}", app.showSnippet).Methods("GET")

	mux.HandleFunc("/user/signup", app.signupUserForm).Methods("GET")
	mux.HandleFunc("/user/signup", app.signupUser).Methods("POST")
	mux.HandleFunc("/user/signup", app.loginUserForm).Methods("GET")
	mux.HandleFunc("/user/signup", app.loginUser).Methods("POST")
	mux.HandleFunc("/user/signup", app.logoutUser).Methods("POST")

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))
	// Pass the servemux as the 'next' parameter to the secureHeaders middleware
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else
	return standardMiddleware.Then(mux)
}
