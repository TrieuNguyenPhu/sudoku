package main

import (
	"log"

	"sudoku/internal/http"
)

func main() {
	r := http.NewRouter()

	log.Println("Sudoku server is running at http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
