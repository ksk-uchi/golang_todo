package main

import (
	"log"
	"todo-app/di"
)

func main() {
	app, cleanup, err := di.InitializeApp()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	app.Router.Setup(app.Engine)

	if err := app.Engine.Start(":8080"); err != nil {
		app.Engine.Logger.Error("shutting down the server")
	}
}
