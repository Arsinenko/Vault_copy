package main

// база данных postgresql
import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Entites struct {
	Id   int
	Name string
	Age  int
}

func Connect() {
	var db, err = sql.Open("postgres", "user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully connected!")
}
func main() {
	Connect()
}
