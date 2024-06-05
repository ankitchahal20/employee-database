package main

import (
	"assignment/internal/config"
	"assignment/internal/db"
	"assignment/internal/server"
	"assignment/internal/service"
	"assignment/internal/utils"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Initializing the Log client
	utils.InitLogClient()

	// Initializing the GlobalConfig
	err := config.InitGlobalConfig()
	if err != nil {
		log.Fatalf("Unable to initialize global config")
	}

	// Establishing the connection to DB.
	postgres, err := db.New()
	if err != nil {
		log.Fatal("Unable to connect to DB : ", err)
	}

	// Initializing the client for employee records service
	//service.NewEmployeeService(postgres)
	_ = service.NewEmployeeService(postgres)

	// Starting the server
	server.Start()
}
