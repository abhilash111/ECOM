package main

import (
	"database/sql"
	"log"

	"github.com/abhilash111/ecom/cmd/api"
	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	db, err := db.NewMySqlDB(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	inittStorage(db)
	server := api.NewApiServer("localhost:8080", db)
	if err := server.Start(); err != nil {
		log.Println("Error starting server", err)
	}
}

func inittStorage(db *sql.DB) {
	err :=
		db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
