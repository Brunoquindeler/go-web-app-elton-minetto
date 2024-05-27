package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/brunoquindeler/go-web-app-elton-minetto/core/beer"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

const errCorvertingToJSON = "erro convertendo em JSON"

func MakeBeerHandlers(router *mux.Router, middleware *negroni.Negroni, service beer.UseCase) {
	router.Handle("/v1/beer", middleware.With(
		negroni.Wrap(getAllBeer(service)),
	)).Methods("GET", "OPTIONS")

	router.Handle("/v1/beer/{id}", middleware.With(
		negroni.Wrap(getBeer(service)),
	)).Methods("GET", "OPTIONS")

	router.Handle("/v1/beer", middleware.With(
		negroni.Wrap(storeBeer(service)),
	)).Methods("POST", "OPTIONS")

	router.Handle("/v1/beer/{id}", middleware.With(
		negroni.Wrap(updateBeer(service)),
	)).Methods("PUT", "OPTIONS")

	router.Handle("/v1/beer/{id}", middleware.With(
		negroni.Wrap(removeBeer(service)),
	)).Methods("DELETE", "OPTIONS")
}

func getAllBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		beers, err := service.GetAll()
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(beers)
		if err != nil {
			writeResponseError(w, errCorvertingToJSON, http.StatusInternalServerError)
			return
		}
	})
}

func getBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}

		beer, err := service.Get(id)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(beer)
		if err != nil {
			writeResponseError(w, errCorvertingToJSON, http.StatusInternalServerError)
			return
		}
	})
}

func storeBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b beer.Beer
		err := json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: validar os dados
		_, err = service.Store(&b)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

func updateBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}

		b, err := service.Get(id)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewDecoder(r.Body).Decode(b)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}

		// TODO: validar os dados
		err = service.Update(b)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func removeBeer(service beer.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = service.Remove(id)
		if err != nil {
			writeResponseError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
