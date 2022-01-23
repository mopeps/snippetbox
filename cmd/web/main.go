package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mopeps/snippetbox/pkg/models/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *pgql.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()   

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	
	// To keep the main() function tidy I've put the code for creating a connection
	// pool into the separate openDB() funciton below. We pass openDB() the dsn
	// from the command-line flag.
	db, err := openDB()
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &pgql.SnippetModel{DB: db},
	}

	srv := &http.Server {
		Addr:	*addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
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
