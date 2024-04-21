package main

import (
	"fmt"

	"github.com/lMikadal/assessment-tax/postgres"
)

func main() {
	db, err := postgres.New()
	if err != nil {
		panic(err)
	}

	_ = db
	fmt.Printf("Connected to postgres database\n")
}
