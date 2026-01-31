package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func (w http.ResponseWriter, r *http.Request) {
		fmt.Println("hrllo")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})
	if err := http.ListenAndServe(":9100", mux); err != nil {
		log.Fatalln("cannot start server")
	}
}
