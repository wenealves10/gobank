package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/wenealves10/gobank/api"
	db "github.com/wenealves10/gobank/db/sqlc"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://gobank:gobank1234@localhost:5434?sslmode=disable&database=gobank"
    serverAddress = "0.0.0.0:8080"

)

func main() {
    conn, err := sql.Open(dbDriver, dbSource)
    if err != nil {
        log.Fatal("cannot conn to db:", err)
    }

    store := db.NewStore(conn)
    server := api.NewServer(store)

    err = server.Start(serverAddress)
    if err != nil {
        log.Fatal("cannot start server", err)
    }
}