package main

import (
	"log"
	"net/http"
	"os"

	"github.com/thaironsilva/messenger/api/router"
	"github.com/thaironsilva/messenger/config"
)

func main() {
	l := log.New(os.Stdout, "", log.LstdFlags)
	l.SetFlags(0)

	db := config.NewDB()

	defer db.Close()

	r := router.New(l, db)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	server.ListenAndServe()
}
