package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)
func main() {

	r := mux.NewRouter()

	r.HandleFunc("/signup", SignUp).Methods("POST")
	r.HandleFunc("/verify/{uuid}", Verify)
	r.HandleFunc("/login", CoursesCtrl).Methods("POST")
	r.HandleFunc("/courses", GetCourses).Methods("GET")

	fmt.Println("Server Is Listening: ")

	http.ListenAndServe(":8080", r)
}
