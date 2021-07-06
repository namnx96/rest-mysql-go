package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"net/http"
	"rest-mysql-go/db"
	postHandler "rest-mysql-go/handler"
)

const (
	Port = ":8080"
)

func main() {
	if err := db.Open(); err != nil {
		// handle error
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/posts", postHandler.GetPosts).Methods("GET")
	router.HandleFunc("/api/v1/posts", postHandler.CreatePost).Methods("POST")
	router.HandleFunc("/api/v1/posts/{id}", postHandler.GetPostById).Methods("GET")
	router.HandleFunc("/api/v1/posts", postHandler.UpdatePost).Methods("PUT")
	router.HandleFunc("/api/v1/posts/{id}", postHandler.DeletePost).Methods("DELETE")

	_ = http.ListenAndServe(Port, router)
}
