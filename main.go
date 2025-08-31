package main

import (
	"context"
	"log"

	"github.com/bolusarz/task-manager/api"
	db "github.com/bolusarz/task-manager/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conn, err := pgxpool.New(context.Background(), "postgresql://postgres:secret@localhost:5432/taskmanager?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	server.StartServer()
}
