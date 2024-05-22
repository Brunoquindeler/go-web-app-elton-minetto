package beer_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/brunoquindeler/go-web-app-elton-minetto/core/beer"
)

func TestStore(t *testing.T) {
	originalBeer := &beer.Beer{
		Name:  "Test",
		Type:  beer.TypeLager,
		Style: beer.StylePale,
	}

	dbTest, err := sql.Open("sqlite3", "../../data/beer_test.db")
	if err != nil {
		t.Fatalf("erro conectando ao banco de dados: %s", err.Error())
	}
	defer dbTest.Close()

	if _, err := dbTest.Exec(beer.CreateDBQuery); err != nil {
		t.Fatalf("erro ao criar banco de dados: %s", err.Error())
	}

	service := beer.NewService(dbTest)

	id, err := service.Store(originalBeer)
	if err != nil {
		t.Fatalf("erro salvando no banco de dados: %s", err.Error())
	}

	id64 := int64(id)
	originalBeer.ID = id64

	retrievedBeer, err := service.Get(id64)
	if err != nil {
		t.Fatalf("erro buscando do banco de dados: %s", err.Error())
	}

	if !reflect.DeepEqual(originalBeer.ID, retrievedBeer.ID) ||
		!reflect.DeepEqual(originalBeer.Name, retrievedBeer.Name) ||
		!reflect.DeepEqual(originalBeer.Type, retrievedBeer.Type) ||
		!reflect.DeepEqual(originalBeer.Style, retrievedBeer.Style) {
		t.Fatalf("erro na comparação das cervejas.\nesperada %+v\nobtida %+v", originalBeer, retrievedBeer)
	}

	err = service.Update(originalBeer)
	if err != nil {
		t.Fatalf("erro atualizando no banco de dados: %s", err.Error())
	}

	beers, err := service.GetAll()
	if err != nil {
		t.Fatalf("erro buscando todas as cervejas do banco de dados: %s", err.Error())
	}

	if len(beers) != 1 {
		t.Fatalf("erro buscando todas as cervejas do banco de dados: %s", err.Error())
	}

	err = service.Remove(id64)
	if err != nil {
		t.Fatalf("erro removendo do banco de dados: %s", err.Error())
	}
}
