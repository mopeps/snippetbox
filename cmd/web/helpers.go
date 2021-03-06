package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// Create an addDefaultData helper method. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns
// the pointer. Again, we 're not using the *http.Request parameter at the moment,
// but we will do later in the book

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.IsAuthenticated = app.isAuthenticated()

	return td
}
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	// Retrieve the appropiate template set from the cache based on the page name
	// (like 'home.page.tmpl'). If no entry exists in the cache with the
	// provided name, call the serverError helper method that we made earlier.
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template %s does not exists", name))
		return
	}

	// Initialize a new buffer
	buf := new(bytes.Buffer)

	// Execute the template set, passin in any dynamic data
	err := ts.Execute(w, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
	}

	// Write the contents of the buffer to the http.ResponseWriter. Again, this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer
	buf.WriteTo(w)
}

// The serverError helper writes an error message and stack trace to the errorLog
// then sends a generic 500 internal Server Error response to the user .

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%sv\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) isAuthenticated() bool {
	return app.session.Values["authenticatedUserId"] != nil
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
