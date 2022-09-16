package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"simplebank/api"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"time"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Establish new connection to database
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer func() {
		ch := make(chan bool, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go func() {
			conn.Close()
			ch <- true
		}()
		select {
		case <-ctx.Done():
			log.Println("Database forced to disconnect")
			return
		case <-ch:
			log.Println("Disconnected from database")
		}
	}()

	store := db.NewStore(conn)

	// Start the server
	server := api.NewServer(store)
	server.Start(config.ServerAddress)
}
