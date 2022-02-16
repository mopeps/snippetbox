package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/mopeps/snippetbox/pkg/forms"
	"github.com/mopeps/snippetbox/pkg/models"
)

// Define a templateData type to act as athe holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as as the build progresses.

type templateData struct {
	CurrentYear int
	Flash interface{}
	Form *forms.Form
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

// Create a humanDate function which returns a nicely formatted string 
// representation of a time.Time object.

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob funtion to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one
	for _, page := range pages {
		// Extract the file name(like 'home.page.tmpl') from the full filepath
		// and assign it to the name variable
		name := filepath.Base(page)

		// Parse the page template file in to a template set.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseHlov method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout  at the moment)

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParserGlob method to add any 'partial' templates to the 
		// template set (in our case, it's just the 'footer' partial at the
		// moment)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'home.page.tmpl') as the key.
		cache[name] = ts

		
	}
	
	// return the map
	return cache, nil
}
