package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "✅ Users service response")
	})

	log.Println("Users service running on :9001")
	http.ListenAndServe(":9001", nil)
}
