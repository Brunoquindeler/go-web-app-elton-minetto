package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brunoquindeler/go-web-app-elton-minetto/core/beer"
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

	service := beer.NewService(db)

	router := mux.NewRouter()

	middleware := negroni.New(
		negroni.NewLogger(),
	)

	router.Handle("/v1/beer", middleware.With(
		negroni.Wrap(hello(service)),
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
