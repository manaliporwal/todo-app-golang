package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
)

func initRouter() {
	router := mux.NewRouter()

	router.HandleFunc("/home", Home)
	router.HandleFunc("/signup", Signup).Methods("POST")
	router.HandleFunc("/signin", Signin).Methods("POST")
	router.HandleFunc("/addTask", CreateTask).Methods("POST")
	router.HandleFunc("/getTasks", GetTasks).Methods("GET")
	router.HandleFunc("/updateTask/{id}/", UpdateTask).Methods("PUT")
	router.HandleFunc("/completeTask/{id}/", CompleteTask).Methods("PUT")
	router.HandleFunc("/deleteTask/{id}/", DeleteTask).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
