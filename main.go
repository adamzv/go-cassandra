package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func main() {
	fmt.Println("Hello, World!")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)
	log.Fatal(http.ListenAndServe(":80", router))
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}