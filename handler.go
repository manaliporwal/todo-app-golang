package main

import (
	"fmt"
	"time"
	"errors"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Status = 0 //Incomplete
// Status = 1 //Complete

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password", db:"password"`
	Username string `json:"username", db:"username"`
}

type Task struct {
	Id int `json: "id", db: "id"`
	Name string `json:"name", db:"name"`
	Description string `json:"description", db:"description"`
	Priority int  `json: "priority", db: "priority"`
	DueDate string `json:"due_date", db: "due_date"`
	Status int `json:"status", db:"status"`
}


var Home = func(w http.ResponseWriter, r *http.Request) {
		userName, err := validateSessionID(w, r)
		if err != nil {
			return
		}
		fmt.Printf("Username is %s", userName)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Welcome to my ToDo App!!")
}

var Signup =  func(w http.ResponseWriter, r *http.Request) {
		// Parse and decode the request body into a new `Credentials` instance
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
				// If there is something wrong with the request body, return a 400 status
				w.WriteHeader(http.StatusBadRequest)
				return
		}
		// Salt and hash the password using the bcrypt algorithm
		// The second argument is the cost of hashing, which we arbitrarily set as 8 (
		// this value can be more or less, depending on the computing power you wish to utilize)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)

		// Next, insert the username, along with the hashed password into the database
		if _, err = db.Query("insert into users values ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
				// If there is any issue with inserting into the database, return a 500 error
				w.WriteHeader(http.StatusInternalServerError)
				return
		}
		// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
}

func Signin(w http.ResponseWriter, r *http.Request){
		// Parse and decode the request body into a new `Credentials` instance
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
				// If there is something wrong with the request body, return a 400 status
				w.WriteHeader(http.StatusBadRequest)
				return
		}
		// Get the existing entry present in the database for the given username
		result := db.QueryRow("select password from users where username=$1", creds.Username)
		if err != nil {
				// If there is an issue with the database, return a 500 error
				w.WriteHeader(http.StatusInternalServerError)
				return
		}
		fmt.Printf("Username %s, Password %s\n", creds.Username, creds.Password)
		// We create another instance of `Credentials` to store the credentials we get from the database
		storedCreds := &Credentials{}
		// Store the obtained password in `storedCreds`
		err = result.Scan(&storedCreds.Password)
		if err != nil {
				// If an entry with the username does not exist, send an "Unauthorized"(401) status
				if err == sql.ErrNoRows {
						w.WriteHeader(http.StatusUnauthorized)
						return
				}

				// If the error is of any other type, send a 500 status
				w.WriteHeader(http.StatusInternalServerError)
				return
		}
		// Compare the stored hashed password, with the hashed version of the password that was received
		fmt.Printf("Input password %s, stored Password %s\n", creds.Password, storedCreds.Password)

		if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
				fmt.Printf("Passwords dont match\n")
				// If the two passwords don't match, return a 401 status
				w.WriteHeader(http.StatusUnauthorized)
				return
		}

		var sessionToken string
		if sessionToken, err = addSession(creds.Username, "120"); err != nil {
			// If there is an error in generating and setting up cache, return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("Session Token %s, Password %s\n", sessionToken, storedCreds.Password)
		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds, the same as the cache
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(120 * time.Second),
		})
		// If we reach this point, that means the users password was correct, and that they are authorized
		// The default 200 status is sent
}

var GetTasks = func(w http.ResponseWriter, r *http.Request) {
		userName, err := validateSessionID(w, r)
		if err != nil {
				return
		}
		fmt.Printf("Username = %s\n", userName)
		var tasks []Task

		SqlStatement := `
		SELECT * FROM tasks where username like $1
		`

		rows, err := db.Query(SqlStatement, userName)
		if err != nil {
				panic(err)
		}

		for rows.Next() {
				var Id int
				var Name string
				var Description string
				var Priority int
				var Status int
				var DueDate string
				rows.Scan(&Id, &Name, &Description, &Priority, &Status, &DueDate)

				tasks = append(tasks, Task{
						Id: Id,
						Name: Name,
						Description: Description,
						Priority: Priority,
						Status: Status,
						DueDate: DueDate,
				})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

}

var CreateTask = func(w http.ResponseWriter, r *http.Request) {
	userName, err := validateSessionID(w, r)
	if err != nil {
		return
	}
	var task Task

	print("Got post request, %s", r.Body)
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sqlStatement := `
		INSERT INTO tasks (name, description, priority, due_date, status, username)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
		`

	id := 0
	err = db.QueryRow(sqlStatement, task.Name, task.Description, task.Priority, task.DueDate, task.Status, userName).Scan(&id)
	if err != nil {
		panic(err)
	}

	task.Id = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

var UpdateTask = func(w http.ResponseWriter, r *http.Request) {
		userName, err := validateSessionID(w, r)
		if err != nil {
				return
		}
		params := mux.Vars(r)
		var todo Task

		r.ParseForm() // must be called to access r.FormValue()

		SqlStatement := `
				UPDATE tasks
				SET name = $1, description = $2
				WHERE id = $3 and username like $4
				RETURNING *
				`

		err = db.QueryRow(
				SqlStatement,
				r.FormValue("name"),
				r.FormValue("description"),
				params["id"],
				userName,
				).Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Priority, &todo.DueDate, &todo.Status)

		if err != nil {
				panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&todo)
}

var CompleteTask = func(w http.ResponseWriter, r *http.Request) {
		userName, err := validateSessionID(w, r)
		if err != nil {
				return
		}
		params := mux.Vars(r)
		var todo Task

		SqlStatement := `
				UPDATE tasks
				SET status = 1
				WHERE id = $1 and username = $2
				RETURNING *
				`

		err = db.QueryRow(
				SqlStatement,
				params["id"], userName,
				).Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Priority, &todo.DueDate, &todo.Status)

		if err != nil {
				panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&todo)
}
var DeleteTask = func(w http.ResponseWriter, r *http.Request) {
		userName, err := validateSessionID(w, r)
		if err != nil {
				return
		}
		params := mux.Vars(r)
		var todo Task

		SqlStatement := `
				DELETE FROM tasks
				WHERE id = $1 and username = $2
				RETURNING *
				`

		err = db.QueryRow(
				SqlStatement,
				params["id"], userName,
				).Scan(&todo.Id, &todo.Name, &todo.Description, &todo.Priority, &todo.DueDate, &todo.Status)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&todo)
}

func validateSessionID(w http.ResponseWriter, r *http.Request) (string, error) {
		// We can obtain the session token from the requests cookies, which come with every request
		c, err := r.Cookie("session_token")
		if err != nil {
				if err == http.ErrNoCookie {
						// If the cookie is not set, return an unauthorized status
						w.WriteHeader(http.StatusUnauthorized)
						return "", err
				}
				fmt.Printf("Error finding cookie\n")
				// For any other type of error, return a bad request status
				w.WriteHeader(http.StatusBadRequest)
				return "", err
		}
		sessionToken := c.Value
		fmt.Printf("Found cookie %s\n", sessionToken)
		// We then get the name of the user from our cache, where we set the session token
		if !validateSession(sessionToken) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Printf("Failed to validate cookie %s\n", sessionToken)
			return "", errors.New("sessionToken not valid\n")
		}
		if userName, ok := getUserNameFromSession(sessionToken); ok {
			return userName, nil
		}
		return "", errors.New("sessionToken not valid\n") 
}
