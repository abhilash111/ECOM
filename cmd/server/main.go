package main

import (
	"database/sql"
	"log"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/router"
	"github.com/abhilash111/ecom/pkg/db"
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
	initStorage(db)
	server := router.NewApiServer("0.0.0.0:8080", db)
	if err := server.Start(); err != nil {
		log.Println("Error starting server", err)
	}
}

func initStorage(db *sql.DB) {
	err :=
		db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
