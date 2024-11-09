package main

import (
	"log"
	"net/http"

	"github.com/Safiramdhn/project-app-inventaris-golang-safira/routers"
)

func main() {
	r := routers.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	log.Println("Server started on port 8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v\n", err.Error())
	}

	http.ListenAndServe(":8080", r)
}
