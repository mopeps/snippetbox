package main

import (
	"net/http"
	
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Create a middleware chain containing our 'standar' middleware
	// which will be used for every request our application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := mux.NewRouter()
	mux.HandleFunc("/", app.home).Methods("GET")
	mux.HandleFunc("/snippet/create", app.createSnippet).Methods("POST")
	mux.HandleFunc("/snippet/create", app.createSnippetForm).Methods("GET")
	mux.HandleFunc("/snippet/{id}", app.showSnippet).Methods("GET")

	fileServer := http.FileServer(http.Dir("./ui/static/"))
    mux.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServer))
	// Pass the servemux as the 'next' parameter to the secureHeaders middleware
	// Because secureHeaders is just a function, and the function returns a 
	// http.Handler we don't need to do anything else
	return standardMiddleware.Then(mux)
}

