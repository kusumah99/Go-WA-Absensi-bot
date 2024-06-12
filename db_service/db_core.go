package dbservice

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func Init(dialect, address string) (*sql.DB, error) {
	db, err := sql.Open(dialect, address)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// container := NewWithDB(db, dialect, log)

	return db, nil
}

func GetDb() *sql.DB {
	return db
}

// untuk perintah create, insert, update dan delete jika butuh rollback
func ExecCommit(sqlString string) error {
	var tx *sql.Tx
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sqlString)

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err
}

func Exec(sqlString string) error {
	// _, err := db.Exec("CREATE TABLE IF NOT EXISTS absensizram (version INTEGER)")
	_, err := db.Exec(sqlString)
	if err != nil {
		return err
	}

	return nil
}

func QueryOneRow(sqlString string) *sql.Row {
	// row := db.QueryRow("SELECT version FROM absensizram LIMIT 1")
	row := db.QueryRow(sqlString)
	if row != nil {
		// _ = row.Scan(&version)
		return row
	}

	return nil
}

func Query(sqlString string) (*sql.Rows, error) {
	// rows, err := db.Query("SELECT version FROM absensizram")
	rows, err := db.Query(sqlString)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
