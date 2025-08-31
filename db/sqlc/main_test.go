package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	testDb, err := pgxpool.New(context.Background(), "postgresql://postgres:secret@localhost:5432/taskmanager?sslmode=disable")

	if err != nil {
		log.Fatal("Could not connect to db", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
