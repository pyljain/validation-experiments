package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/api/v1/client", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("No client found"))
	})

	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		log.Printf("Server error has occured %s", err)
	}
}
