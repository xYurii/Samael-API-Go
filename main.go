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

	// my token: a3e6hznQXpbwyTfzfldH86EjhKouor

	router.Initiliaze()
}
