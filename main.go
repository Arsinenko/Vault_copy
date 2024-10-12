package main

import (
	"Vault_copy/db_operations"
	"fmt"
)

func main() {
	db, err := db_operations.DB_connection()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	db_operations.CreateUser(db, "TestUser", "+79123456789", "test@mail.com", "testpass", nil, nil)
	db.Commit()
	err = db.Close()
	if err != nil {
		return
	}
}
