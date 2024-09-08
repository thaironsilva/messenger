package main

import (
	"log"
	"net/http"

	"github.com/thaironsilva/messenger/api/router"
	"github.com/thaironsilva/messenger/config"
)

func main() {
	db := config.NewDB()

	defer db.Close()

	r := router.New(db)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
