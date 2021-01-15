package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/adamzv/go-cassandra/Cassandra"
	"github.com/adamzv/go-cassandra/Messages"
	"github.com/adamzv/go-cassandra/Users"
	"github.com/gorilla/mux"
)

type heartbeatResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
}

func main() {
	fmt.Println("Hello, World!")

	CassandraSession := Cassandra.Session
	defer CassandraSession.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", heartbeat)

	router.HandleFunc("/users", Users.Get)
	router.HandleFunc("/users/{user_uuid}", Users.GetOne)
	router.HandleFunc("/users/new", Users.Post)

	router.HandleFunc("/messages", Messages.Get)
	router.HandleFunc("/messages/{message_uuid}", Messages.GetOne)
	router.HandleFunc("/messages/new", Messages.Post)

	log.Fatal(http.ListenAndServe(":80", router))
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(heartbeatResponse{Status: "OK", Code: 200})
}
