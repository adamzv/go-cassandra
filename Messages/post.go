package Messages

import (
	"encoding/json"
	"net/http"

	"github.com/adamzv/go-cassandra/Cassandra"
	"github.com/gocql/gocql"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var errStr, userIDStr, message string

	if userIDStr, errStr = processFormField(r, "userID"); len(errStr) != 0 {
		errs = append(errs, errStr)
	}
	userID, err := gocql.ParseUUID(userIDStr)
	if err != nil {
		errs = append(errs, "Parameter 'userID' not a UUID")
	}

	if message, errStr = processFormField(r, "message"); len(errStr) != 0 {
		errs = append(errs, errStr)
	}

	gocqlUUID := gocql.TimeUUID()

	var created bool = false
	if len(errs) == 0 {
		if err := Cassandra.Session.Query(`
		INSERT INTO message (id, userID, message) VALUES (?, ?, ?)`,
			gocqlUUID, userID, message).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
	}
	if created {
		json.NewEncoder(w).Encode(NewMessageResponse{ID: gocqlUUID})
	} else {
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
