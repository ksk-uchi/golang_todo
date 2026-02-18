package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"todo-app/ent"

	_ "github.com/go-sql-driver/mysql"
)

type Todo struct {
	Title       string
	Description string
}

func main() {
	client, err := ent.Open("mysql", "user:password@tcp(localhost:3306)/todo_db?parseTime=True")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		if err := client.Close(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	var todos []Todo
	for i := 0; i < 1000; i++ {
		todos = append(todos, Todo{
			Title:       "title " + strconv.Itoa(i),
			Description: "description " + strconv.Itoa(i),
		})
	}
	createdTodos, err := client.Todo.MapCreateBulk(todos, func(c *ent.TodoCreate, i int) {
		c.
			SetTitle(todos[i].Title).
			SetDescription(todos[i].Description).
			SetUserID(1)
	}).Save(context.Background())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(len(createdTodos))
}
