package main

import (
	"fmt"
	"permissionist"
)

func main(){
	db, err := permissionist.InitDB("blah")
	if err != nil {
		fmt.Println(err)
	}
	P := permissionist.Permissionist{
		DB: db,
	}
	fmt.Println(P)
}
