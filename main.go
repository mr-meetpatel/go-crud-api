package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Article struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var db *sql.DB

func createTable() {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS articles (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255),
			content TEXT
		)
	`)
	if err != nil {
		log.Fatal("Error creating articles table:", err)
	}
}

func initDB() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbURI := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)

	conn, err := sql.Open("postgres", dbURI)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}

	db = conn

	fmt.Println("Connected to the database")
	createTable()
}

func closeDB() {
	if db != nil {
		db.Close()
		fmt.Println("Disconnected from the database")
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to home page")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	rows, err := db.Query("SELECT * FROM articles")
	if err != nil {
		fmt.Println("Error querying articles:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	articles := make([]Article, 0)
	for rows.Next() {
		var article Article
		err := rows.Scan(&article.Id, &article.Title, &article.Content)
		if err != nil {
			fmt.Println("Error scanning article:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		articles = append(articles, article)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("Error iterating over rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["id"]
	fmt.Println("Endpoint Hit: returnSingleArticle")

	row := db.QueryRow("SELECT * FROM articles WHERE id = $1", articleId)
	var article Article
	err := row.Scan(&article.Id, &article.Title, &article.Content)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(struct {
				Message string `json:"message"`
			}{
				Message: "Article not found",
			})
		} else {
			fmt.Println("Error scanning article:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewArticle")
	payload, _ := ioutil.ReadAll(r.Body)
	var newArticle Article
	err := json.Unmarshal(payload, &newArticle)
	if err != nil {
		fmt.Println("Error while decode payload:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	query := "INSERT INTO articles (title, content) VALUES ($1, $2) RETURNING id, title, content"
	err = db.QueryRow(query, newArticle.Title, newArticle.Content).Scan(&newArticle.Id, &newArticle.Title, &newArticle.Content)
	if err != nil {
		fmt.Println("Error creating article:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newArticle)
}

func deleteArticleById(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["id"]
	fmt.Println("Endpoint Hit: deleteArticleById")

	result, err := db.Exec("DELETE FROM articles WHERE id = $1", articleId)
	if err != nil {
		fmt.Println("Error deleting article:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: "Article not found",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	articleId := mux.Vars(r)["id"]
	fmt.Println("Endpoint Hit: updateArticle")
	payload, _ := ioutil.ReadAll(r.Body)
	var updatedArticle Article
	err := json.Unmarshal(payload, &updatedArticle)
	if err != nil {
		fmt.Println("Error while decode payload:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := db.Exec("UPDATE articles SET title = $1, content = $2 WHERE id = $3", updatedArticle.Title, updatedArticle.Content, articleId)
	if err != nil {
		fmt.Println("Error updating article:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: "Article not found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedArticle)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage).Methods("GET")
	myRouter.HandleFunc("/articles", returnAllArticles).Methods("GET")
	myRouter.HandleFunc("/articles/{id}", returnSingleArticle).Methods("GET")
	myRouter.HandleFunc("/articles", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/articles/{id}", deleteArticleById).Methods("DELETE")
	myRouter.HandleFunc("/articles/{id}", updateArticle).Methods("PUT")

	fmt.Printf("Server Start on port 8000\n")
	http.ListenAndServe(":8000", myRouter)
}

func main() {
	initDB()
	defer closeDB()
	handleRequests()
}
