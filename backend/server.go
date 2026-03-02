package main

import (
	"fmt"
	"log"
	"os"
	"todo-app/di"

	"github.com/joho/godotenv"
)

func main() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "local"
	}

	if err := godotenv.Load(fmt.Sprintf("envs/%s.env", appEnv)); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	app, cleanup, err := di.InitializeApp()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	app.Router.Setup(app.Engine)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := app.Engine.Start(fmt.Sprintf(":%s", port)); err != nil {
		app.Engine.Logger.Error("shutting down the server")
	}
}
