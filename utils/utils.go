package utils

import "database/sql"

func ClearDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM beer")
	tx.Commit()
	return err
}
