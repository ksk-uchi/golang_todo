package providers

import (
	"fmt"
	"os"
	"todo-app/ent"

	_ "github.com/go-sql-driver/mysql"
)

func NewEntClient() (*ent.Client, func(), error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "user:password@tcp(localhost:3306)/todo_db?parseTime=True"
	}
	client, err := ent.Open("mysql", dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed opening mysql: %w", err)
	}

	cleanup := func() {
		_ = client.Close()
	}

	return client, cleanup, nil
}
