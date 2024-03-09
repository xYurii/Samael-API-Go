package main

import (
	"apisamael/database"
	"apisamael/entities"
	"apisamael/router"
	"context"
	"fmt"

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

	u := entities.User{
		ID:  "fodase",
		Tag: "sua mãe.",
	}

	a := database.User.GetUser(context.Background(), u)
	fmt.Println(a, a.ID)

	router.Initiliaze()
}
