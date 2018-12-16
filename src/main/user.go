package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() {
	fmt.Println("DEDANS")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		fmt.Println(err)
	}
	r := db.Close()
	if r != nil {
		fmt.Println(r)
	}
}
