package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"

	"github.com/mopeps/snippetbox/pkg/models/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session		  *sessions.Session
	snippets      *pgql.SnippetModel
	templateCache map[string]*(template.Template)
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	secret := flag.String("secret", "super-secret-key", "Secret key")
	flag.Parse()
	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() funciton below. We pass openDB() the dsn
	// from the command-line flag.
	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}
	
	cookieStore := sessions.NewCookieStore([]byte(*secret))
	session := sessions.NewSession(cookieStore, "our-session")
	session.Options.Secure = true

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:	   session,
		snippets:      &pgql.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	
	tlsConfig := &tls.Config{
    PreferServerCipherSuites: true,
    CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
