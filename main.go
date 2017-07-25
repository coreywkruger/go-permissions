package main

import (
	"fmt"
	// "permissionist"
)

func main(){
	db, err := InitDB(DbConfig{"db", "db", "db", "db", "5432"})
	if err != nil {
		fmt.Println(err)
	}
	P := Permissionist{
		DB: db,
	}
	fmt.Println(P)
}
