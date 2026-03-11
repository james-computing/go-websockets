package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Hello")
	setupAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setupAPI() {
	manager := NewManager()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.serveWS)
}
