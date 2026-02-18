package providers

import (
	"fmt"
	"todo-app/ent"

	_ "github.com/go-sql-driver/mysql"
)

func NewEntClient() (*ent.Client, func(), error) {
	client, err := ent.Open("mysql", "user:password@tcp(localhost:3306)/todo_db?parseTime=True")
	if err != nil {
		return nil, nil, fmt.Errorf("failed opening mysql: %w", err)
	}

	cleanup := func() {
		_ = client.Close()
	}

	return client, cleanup, nil
}
