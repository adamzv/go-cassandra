package Messages

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adamzv/go-cassandra/Cassandra"
	"github.com/adamzv/go-cassandra/Users"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func Get(w http.ResponseWriter, r *http.Request) {
	var messageList []Message
	var enrichedMessages []Message
	var userList []gocql.UUID
	m := map[string]interface{}{}

	query := "SELECT id, user_id, message FROM messages"
	iterable := Cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		userID := m["userID"].(gocql.UUID)
		messageList = append(messageList, Message{
			ID:      m["id"].(gocql.UUID),
			UserID:  userID,
			Message: m["message"].(string),
		})
		userList = append(userList, userID)
		m = map[string]interface{}{}
	}

	names := Users.Enrich(userList)
	for _, message := range messageList {
		message.UserFullName = names[message.UserID.String()]
		enrichedMessages = append(enrichedMessages, message)
	}
	fmt.Println("message list after enrichment", enrichedMessages)

	json.NewEncoder(w).Encode(AllMessagesResponse{Messages: enrichedMessages})
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	var message Message
	var errs []string
	var found bool = false

	vars := mux.Vars(r)
	id := vars["message_uuid"]

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		m := map[string]interface{}{}
		query := "SELECT id, user_id, message FROM messages WHERE id = ? LIMIT 1"
		iterable := Cassandra.Session.Query(query, uuid).Consistency(gocql.One).Iter()
		for iterable.MapScan(m) {
			found = true
			userID := m["userID"].(gocql.UUID)
			names := Users.Enrich([]gocql.UUID{userID})
			fmt.Println("names", names)
			message = Message{
				ID:           userID,
				UserID:       m["userID"].(gocql.UUID),
				UserFullName: names[userID.String()],
				Message:      m["message"].(string),
			}
		}
		if !found {
			errs = append(errs, "Message not found")
		}
	}
	if found {
		json.NewEncoder(w).Encode(GetMessageResponse{Message: message})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
