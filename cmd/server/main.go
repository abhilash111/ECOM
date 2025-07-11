package main

import (
	"log"

	"github.com/abhilash111/ecom/config"
	"github.com/abhilash111/ecom/internal/router"
	"github.com/abhilash111/ecom/pkg/db"
	"github.com/go-sql-driver/mysql"
)

func main() {
	conn, err := db.NewMySqlDB(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	r := router.SetupRouter(conn)

	log.Println("Starting server on port", config.Envs.Port)
	if err := r.Run(":" + config.Envs.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
