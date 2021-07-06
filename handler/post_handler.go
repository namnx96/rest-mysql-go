package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"rest-mysql-go/db"
	"rest-mysql-go/model"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	var posts []model.Post

	result, err := db.DB.Query("select id, title from post")
	panicOnError(err)

	defer result.Close()
	for result.Next() {
		var post model.Post
		err := result.Scan(&post.ID, &post.Title)
		panicOnError(err)
		posts = append(posts, post)
	}
	//return empty array instead of null data
	if posts == nil {
		posts = make([]model.Post, 0)
	}
	responseWithJson(w, http.StatusOK, posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.DB.Prepare("INSERT INTO post(title) VALUES(?)")
	panicOnError(err)

	body, err := ioutil.ReadAll(r.Body)
	panicOnError(err)

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	title := keyVal["title"]
	result, err := stmt.Exec(title)
	panicOnError(err)

	responseWithJson(w, http.StatusCreated, result)
}

func findById(id string) *model.Post {
	result, err := db.DB.Query("SELECT id, title FROM post WHERE id = ?", id)
	panicOnError(err)

	defer result.Close()
	var post model.Post
	for result.Next() {
		err := result.Scan(&post.ID, &post.Title)
		panicOnError(err)
	}
	if post.ID == 0 {
		return nil
	} else {
		return &post
	}
}
func GetPostById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["id"]
	post := findById(postId)
	if post == nil {
		responseWithJson(w, http.StatusNotFound, map[string]string{"message": "Post with Id: " + postId + " not found"})
	} else {
		responseWithJson(w, http.StatusOK, post)
	}
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["title"]
	postId := keyVal["id"]

	post := findById(postId)
	if post == nil {
		responseWithJson(w, http.StatusNotFound, map[string]string{"message": "Post with Id: " + postId + " not found"})
		return
	}

	stmt, err := db.DB.Prepare("UPDATE post SET title = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec(newTitle, postId)
	if err != nil {
		responseWithJson(w, http.StatusNotFound,
			map[string]string{"message": "Post with Id: " + postId + " cannot be updated with error: " + err.Error()})
	} else {
		responseWithJson(w, http.StatusOK, "Post with ID = "+postId+" was updated")
	}

}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.DB.Prepare("DELETE FROM post WHERE id = ?")
	panicOnError(err)

	idToDelete := params["id"]

	post := findById(idToDelete)
	if post == nil {
		responseWithJson(w, http.StatusNotFound, map[string]string{"message": "Post with Id: " + idToDelete + " not found"})
		return
	}

	_, err = stmt.Exec(idToDelete)
	if err != nil {
		responseWithJson(w, http.StatusNotFound,
			map[string]string{"message": "Post with Id: " + idToDelete + " cannot be deleted with error: " + err.Error()})
	}
	responseWithJson(w, http.StatusOK, "Post with ID = "+idToDelete+" was deleted")
}

func panicOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func responseWithJson(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(object)
}
