package Users

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adamzv/go-cassandra/Cassandra"
	"github.com/gocql/gocql"
)

func Post(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var gocqlUUID gocql.UUID

	// FormToUser() -> Users/processing.go
	user, errs := FormToUser(r)

	var created bool = false

	if len(errs) == 0 {
		fmt.Println("Creating a new user")

		// generate a unique UUID for this user
		gocqlUUID = gocql.TimeUUID()

		// save the user to Cassandra
		if err := Cassandra.Session.Query(`
		INSERT INTO users (id, firstname, lastname, email, city, age) VALUES (?, ?, ?, ?, ?, ?)`,
			gocqlUUID, user.FirstName, user.LastName, user.Email, user.City, user.Age).Exec(); err != nil {
			errs = append(errs, err.Error())
		} else {
			created = true
		}
	}

	// depending on whether we created the user, return the
	// resource ID (UUID) in a JSON payload, or return our errors
	if created {
		fmt.Println("user_id ", gocqlUUID)
		json.NewEncoder(w).Encode(NewUserResponse{ID: gocqlUUID})
	} else {
		fmt.Println("errors ", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
