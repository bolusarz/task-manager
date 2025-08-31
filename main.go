package main

import (
	"context"
	"fmt"
	"log"

	"github.com/bolusarz/task-manager/api"
	db "github.com/bolusarz/task-manager/db/sqlc"
	"github.com/bolusarz/task-manager/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal(fmt.Errorf("unable to load config %v", err))
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)

	server, err := api.NewServer(store, config)

	if err != nil {
		log.Fatal(err)
	}

	server.StartServer(config.HTTPServerAddress)
}
