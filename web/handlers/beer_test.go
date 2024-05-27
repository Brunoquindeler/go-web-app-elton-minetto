package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brunoquindeler/go-web-app-elton-minetto/core/beer"
	"github.com/brunoquindeler/go-web-app-elton-minetto/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_getAllBeer(t *testing.T) {
	b1 := &beer.Beer{
		Name:  "Heineken",
		Type:  beer.TypeLager,
		Style: beer.StylePale,
	}
	b2 := &beer.Beer{
		Name:  "Skol",
		Type:  beer.TypeLager,
		Style: beer.StylePale,
	}

	dbTest, err := sql.Open("sqlite3", "../../data/beer_test.db")
	assert.Nil(t, err)
	defer dbTest.Close()

	assert.Nil(t, utils.ClearDB(dbTest))

	_, err = dbTest.Exec(beer.CreateDBQuery)
	assert.Nil(t, err)

	service := beer.NewService(dbTest)

	_, err = service.Store(b1)
	assert.Nil(t, err)

	_, err = service.Store(b2)
	assert.Nil(t, err)

	handler := getAllBeer(service)

	router := mux.NewRouter()

	router.Handle("/v1/beer", handler)

	req, err := http.NewRequest("GET", "/v1/beer", nil)
	assert.Nil(t, err)

	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result []*beer.Beer
	err = json.NewDecoder(rr.Body).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, b1.Name, result[0].Name)
	assert.Equal(t, b2.Name, result[1].Name)
}

type BeerServiceMock struct {
	beers []*beer.Beer
}

func (t *BeerServiceMock) GetAll() ([]*beer.Beer, error) {
	return t.beers, nil
}

func (t *BeerServiceMock) Get(ID int64) (*beer.Beer, error) {
	if t.beers[0].ID != ID {
		return nil, fmt.Errorf("cerveja com ID: %d não encontrada", ID)
	}

	return t.beers[0], nil
}

func (t *BeerServiceMock) Store(b *beer.Beer) (int, error) {
	t.beers = append(t.beers, b)

	return int(t.beers[0].ID), nil
}

func (t *BeerServiceMock) Update(b *beer.Beer) error {
	for i, beer := range t.beers {
		if beer.ID == b.ID {
			t.beers[i] = b
			return nil
		}
	}
	return errors.New("erro ao atualizar cerveja")
}

func (t *BeerServiceMock) Remove(ID int64) error {
	for i, beer := range t.beers {
		if beer.ID == ID {
			t.beers = append(t.beers[:i], t.beers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("cerveja com ID: %d não encontrada", ID)
}

func Test_getAllBeerWithMock(t *testing.T) {
	service := &BeerServiceMock{
		beers: []*beer.Beer{&beer.Beer{
			ID:    1,
			Name:  "Heineken",
			Type:  beer.TypeLager,
			Style: beer.StylePale,
		},
			&beer.Beer{
				ID:    2,
				Name:  "Skol",
				Type:  beer.TypeLager,
				Style: beer.StylePale,
			},
		},
	}

	handler := getAllBeer(service)
	r := mux.NewRouter()
	r.Handle("/v1/beer", handler)

	req, err := http.NewRequest("GET", "/v1/beer", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result []*beer.Beer
	err = json.NewDecoder(rr.Body).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, service.beers[0], result[0])
	assert.Equal(t, service.beers[1], result[1])
}

func Test_getBeerWithMock(t *testing.T) {
	service := &BeerServiceMock{
		beers: []*beer.Beer{&beer.Beer{
			ID:    1,
			Name:  "Heineken",
			Type:  beer.TypeLager,
			Style: beer.StylePale,
		}},
	}
	handler := getBeer(service)
	r := mux.NewRouter()
	r.Handle("/v1/beer/{id}", handler)

	req, err := http.NewRequest("GET", "/v1/beer/1", nil)
	assert.Nil(t, err)
	req.Header.Set("Accept", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	var result *beer.Beer
	err = json.NewDecoder(rr.Body).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, service.beers[0], result)
}

func Test_storeBeerWithMock(t *testing.T) {
	service := &BeerServiceMock{}
	handler := storeBeer(service)
	r := mux.NewRouter()
	r.Handle("/v1/beer", handler)

	b := &beer.Beer{
		ID:    1,
		Name:  "Heineken",
		Type:  beer.TypeLager,
		Style: beer.StylePale,
	}

	bJSON, err := json.Marshal(b)
	assert.Nil(t, err)

	req, err := http.NewRequest("POST", "/v1/beer", bytes.NewBuffer(bJSON))
	assert.Nil(t, err)
	req.Header.Set("Content-type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	assert.Equal(t, service.beers[0], b)
}

func Test_updateBeerWithMock(t *testing.T) {
	service := &BeerServiceMock{
		beers: []*beer.Beer{
			&beer.Beer{
				ID:    1,
				Name:  "Heineken",
				Type:  beer.TypeLager,
				Style: beer.StylePale,
			},
		},
	}
	handler := updateBeer(service)
	r := mux.NewRouter()
	r.Handle("/v1/beer/{id}", handler)

	b := &beer.Beer{
		ID:    1,
		Name:  "Novo Nome",
		Type:  beer.TypeLager,
		Style: beer.StylePale,
	}

	bJSON, err := json.Marshal(b)
	assert.Nil(t, err)

	req, err := http.NewRequest("PUT", "/v1/beer/1", bytes.NewBuffer(bJSON))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)

	assert.Equal(t, service.beers[0], b)
	assert.Equal(t, service.beers[0].Name, "Novo Nome")
}

func Test_removeBeerWithMock(t *testing.T) {
	service := &BeerServiceMock{
		beers: []*beer.Beer{
			&beer.Beer{
				ID:    1,
				Name:  "Heineken",
				Type:  beer.TypeLager,
				Style: beer.StylePale,
			},
		},
	}

	handler := removeBeer(service)
	r := mux.NewRouter()
	r.Handle("/v1/beer/{id}", handler)

	req, err := http.NewRequest("DELETE", "/v1/beer/1", nil)
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)

	assert.Equal(t, 0, len(service.beers))
}
