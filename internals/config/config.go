package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/obimadu/ipc3-stage-2/internals/db"
)

func Config() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Printf("Could not load .env file %s\n", err.Error())
	}

	// init db
	db.InitDB()
}
