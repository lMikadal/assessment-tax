package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Port: %s\n", os.Getenv("PORT"))
	fmt.Printf("Database_url: %s\n", os.Getenv("DATABASE_URL"))
	fmt.Printf("Admin_username: %s\n", os.Getenv("ADMIN_USERNAME"))
	fmt.Printf("Admin_password: %s\n", os.Getenv("ADMIN_PASSWORD"))
}
