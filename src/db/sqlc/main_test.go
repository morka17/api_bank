package db

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/shiny_bank?sslmode=disable"
)

var testQueries1 *Queries

func ConnectDB(t *testing.T) *sql.DB {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		t.Errorf("Expected no error, %v", err)
	}

	testQueries1 = New(conn)

	if conn == nil {
		t.Error("Database connection error")
	}

	return conn

}

func TestDB(t *testing.T) {
	ConnectDB(t)

}
