package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brunoquindeler/go-web-app-elton-minetto/core/beer"
	"github.com/brunoquindeler/go-web-app-elton-minetto/web/handlers"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "data/beer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if _, err := db.Exec(beer.CreateDBQuery); err != nil {
		log.Fatalf("erro ao criar banco de dados: %s", err.Error())
	}

	service := beer.NewService(db)

	router := mux.NewRouter()

	middleware := negroni.New(
		negroni.NewLogger(),
	)

	handlers.MakeBeerHandlers(router, middleware, service)

	fileServer := http.FileServer(http.Dir("./web/static"))
	router.PathPrefix("/static/").Handler(middleware.With(
		negroni.Wrap(http.StripPrefix("/static/", fileServer)),
	)).Methods("GET", "OPTIONS")

	http.Handle("/", router)

	server := &http.Server{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":4000",
		Handler:      http.DefaultServeMux,
		ErrorLog:     log.New(os.Stderr, "logger: ", log.Lshortfile),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func hello(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beers, _ := service.GetAll()
		for _, beer := range beers {
			fmt.Println(beer)
		}
	})
}
