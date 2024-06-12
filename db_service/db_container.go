package dbservice

import (
	"database/sql"
	"fmt"
)

type DbContainer struct {
	db      *sql.DB
	dialect string
}

func aNew(dialect, address string) (*DbContainer, error) {
	db, err := sql.Open(dialect, address)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// container := NewWithDB(db, dialect, log)
	container := aNewWithDB(db, dialect)

	return container, nil
}

func aNewWithDB(db *sql.DB, dialect string) *DbContainer {
	return &DbContainer{
		db:      db,
		dialect: dialect,
		// log:     log,
	}
}

func query(c *DbContainer) error {
	var tx *sql.Tx
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}

	err = createTable(tx, c)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return err

}

func queryRow(c *DbContainer) (int, error) {
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS absensizram (version INTEGER)")
	if err != nil {
		return -1, err
	}

	version := 0

	row := c.db.QueryRow("SELECT version FROM absensizram LIMIT 1")
	if row != nil {
		_ = row.Scan(&version)
	}
	return version, nil
}

func createTable(tx *sql.Tx, _ *DbContainer) error {
	_, err := tx.Exec(`CREATE TABLE absensizram (
		mobile_number       TEXT PRIMARY KEY,
		salik_name          TEXT,
		datetime_riyadhah   DATETIME,
		response            TEXT`)

	return err
}

func createTableAnggota(tx *sql.Tx, _ *DbContainer) error {
	_, err := tx.Exec(`CREATE TABLE absensizram (
		mobile_number       TEXT PRIMARY KEY,
		salik_name          TEXT,
		datetime_riyadhah   DATETIME,
		response            TEXT`)

	return err
}
