package main

import (
	"apisamael/database"
	"apisamael/router"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("cannot read the env file")
	}

	_, err = database.Connect(database.GetEnvConfig())
	if err != nil {
		panic("cannot connect to the database")
	}

	// my token: S50SSP78c7jKH2sQbUVUX1PLcG3m80

	router.Initiliaze()
}
