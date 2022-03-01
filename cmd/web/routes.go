package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	router := mux.NewRouter()
	dynamicRouter := router.PathPrefix("/").Subrouter()
	router.HandleFunc("/", app.home).Methods("GET")
	router.HandleFunc("/snippet/{id}", app.showSnippet).Methods("GET")
	router.HandleFunc("/user/login", app.loginUserForm).Methods("GET")
	router.HandleFunc("/user/login", app.loginUser).Methods("POST")

	dynamicRouter.Use(app.requireAuthentification)
	dynamicRouter.HandleFunc("/snippet/create", app.createSnippet).Methods("POST")
	dynamicRouter.HandleFunc("/snippet/create", app.createSnippetForm).Methods("GET")
	dynamicRouter.HandleFunc("/user/signup", app.signupUserForm).Methods("GET")
	dynamicRouter.HandleFunc("/user/signup", app.signupUser).Methods("POST")
	dynamicRouter.HandleFunc("/user/logout", app.logoutUser).Methods("POST")

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(router)
}
