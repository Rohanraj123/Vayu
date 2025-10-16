package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "✅ Orders service response")
	})

	log.Println("Orders service running on :9002")
	http.ListenAndServe(":9002", nil)
}
