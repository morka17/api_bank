package db

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/morka17/shiny_bank/v1/src/utils"
)

var testQueries1 *Queries

func ConnectDB(t *testing.T) *sql.DB {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		t.Errorf("Failed to load config %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

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
