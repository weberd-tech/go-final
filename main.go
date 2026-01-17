package main

import (
	"fmt"
	"log"
	"net/http"

	"todo/pkg/api"
	"todo/pkg/db"
)

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal(err)
	}

	api.Init()

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	port := ":7540"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
