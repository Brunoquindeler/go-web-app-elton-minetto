package beer

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const CreateDBQuery = `
CREATE TABLE IF NOT EXISTS beer (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    type INTEGER NOT NULL,
    style INTEGER NOT NULL
);
`

type UseCase interface {
	GetAll() ([]*Beer, error)
	Get(ID int64) (*Beer, error)
	Store(beer *Beer) (int, error)
	Update(beer *Beer) error
	Remove(ID int64) error
}

type Service struct {
	DB *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}

func (s *Service) GetAll() ([]*Beer, error) {
	var beers []*Beer

	rows, err := s.DB.Query("SELECT id, name, type, style FROM beer")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var beer Beer

		err = rows.Scan(&beer.ID, &beer.Name, &beer.Type, &beer.Style)
		if err != nil {
			return nil, err
		}

		beers = append(beers, &beer)
	}

	return beers, nil
}

func (s *Service) Get(ID int64) (*Beer, error) {
	var beer Beer

	statement, err := s.DB.Prepare("SELECT id, name, type, style FROM beer WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	err = statement.QueryRow(ID).Scan(&beer.ID, &beer.Name, &beer.Type, &beer.Style)
	if err != nil {
		return nil, err
	}

	return &beer, nil
}

func (s *Service) Store(beer *Beer) (int, error) {
	transaction, err := s.DB.Begin()
	if err != nil {
		return 0, err
	}

	statement, err := transaction.Prepare("INSERT INTO beer(name, type, style) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(beer.Name, beer.Type, beer.Style)
	if err != nil {
		transaction.Rollback()
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	transaction.Commit()
	return int(id), nil
}

func (s *Service) Update(beer *Beer) error {
	if beer.ID == 0 {
		return fmt.Errorf("invalid ID")
	}

	transaction, err := s.DB.Begin()
	if err != nil {
		return err
	}

	statement, err := transaction.Prepare("UPDATE beer SET name=?, type=?, style=? WHERE id=?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(beer.Name, beer.Type, beer.Style, beer.ID)
	if err != nil {
		transaction.Rollback()
		return err
	}

	transaction.Commit()
	return nil
}

func (s *Service) Remove(ID int64) error {
	if ID == 0 {
		return fmt.Errorf("invalid ID")
	}

	transaction, err := s.DB.Begin()
	if err != nil {
		return err
	}

	statement, err := transaction.Prepare("DELETE FROM beer WHERE id=?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(ID)
	if err != nil {
		transaction.Rollback()
		return err
	}

	transaction.Commit()
	return nil
}
