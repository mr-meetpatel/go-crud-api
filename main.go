package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type CustomError struct {
	Message string `json:"message"`
}

var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to home page")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["id"]
	fmt.Println("Endpoint Hit: returnAllArticles")
	w.Header().Set("Content-Type", "application/json")
	for _, article := range Articles {
		if article.Id == articleId {
			json.NewEncoder(w).Encode(article)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(CustomError{Message: "Article not found"})
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")
	payload, _ := ioutil.ReadAll(r.Body)
	var newArticle Article
	json.Unmarshal(payload, &newArticle)
	Articles = append(Articles, newArticle)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newArticle)
}

func deleteArticleById(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["id"]
	fmt.Println("Endpoint Hit: deleteArticleById")
	w.Header().Set("Content-Type", "application/json")
	for index, article := range Articles {

		if article.Id == articleId {
			fmt.Println(Articles[:index], Articles[index+1:])
			w.WriteHeader(http.StatusNoContent)
			Articles = append(Articles[:index], Articles[index+1:]...)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(CustomError{Message: "Article not found"})

}
func updateArticleById(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["id"]
	fmt.Println("Endpoint Hit: updateArticle")
	w.Header().Set("Content-Type", "application/json")
	for index, article := range Articles {
		if article.Id == articleId {
			payload, _ := ioutil.ReadAll(r.Body)
			var updatedArticle Article
			json.Unmarshal(payload, &updatedArticle)
			Articles[index] = updatedArticle
			json.NewEncoder(w).Encode(updatedArticle)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	// json.NewEncoder(w).Encode(struct {
	// 	Message string `json:"message"`
	// }{
	// 	Message: "Article not found",
	// })
	json.NewEncoder(w).Encode(CustomError{Message: "Article not found"})
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
	myRouter.HandleFunc("/articles/{id}", returnSingleArticle).Methods("GET")
	myRouter.HandleFunc("/articles", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/articles/{id}", deleteArticleById).Methods("DELETE")
	myRouter.HandleFunc("/articles/{id}", updateArticleById).Methods("PUT")
	fmt.Printf("Server Start on port 8000")
	http.ListenAndServe(":8000", myRouter)
}

func main() {
	Articles = []Article{
		{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	handleRequests()
}
