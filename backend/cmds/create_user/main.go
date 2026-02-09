package main

import (
	"context"
	"fmt"

	"todo-app/ent"

	_ "github.com/go-sql-driver/mysql"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func main() {
	client, err := ent.Open("mysql", "user:password@tcp(localhost:3306)/todo_db?parseTime=True")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	password, err := hashPassword("password")
	if err != nil {
		fmt.Println(err)
		return
	}
	user, err := client.User.Create().
		SetName("keisuke").
		SetEmail("ksk.uchi+github@gmail.com").
		SetPassword(password).
		Save(context.Background())

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(user.String())
}
