package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/mr-meetpatel/go-crud-api/docs"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Article struct {
	Id      int    `json:"id"`
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
}
type ValidationError struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}
type CustomError struct {
	Message string `json:"message"`
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

// HomePage godoc
// @Summary  Welcome Message
// @Description Welcome Message
// @Tags Home
// @Produce plain
// @Success 200
// @Router / [get]
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to home page")
}

// ReturnAllArticles godoc
// @Summary  Return All Articles
// @Description Return All Articles aviable in Database
// @Tags Article
// @Produce  json
// @Success 200 {object} Article
// @Router /articles [get]
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

// ReturnSingleArticle godoc
// @Summary  Return single Article
// @Description Return single Article by articleId
// @Param id  path string  true  "Article ID"
// @Tags Article
// @Produce  json
// @Success 200 {object} Article
// @Failure 404 {object} CustomError
// @Router /articles/{id} [get]
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

// CreateNewArticle godoc
// @Summary Create a new Article
// @Description Create a new Article with the input paylod
// @Tags Article
// @Accept  json
// @Produce  json
// @Param article body Article true "Create article"
// @Success 201 {object} Article
// @Router /articles [post]
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
	validate := validator.New()
	err = validate.Struct(newArticle)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		validationErrorMessages := make([]ValidationError, 0)
		for _, e := range validationErrors {
			validationErrorMessages = append(validationErrorMessages, ValidationError{
				Key:   fmt.Sprintf("'%s'", e.Field()),
				Error: fmt.Sprintf("Field validation for '%s' failed on the 'required' tag", e.Field()),
			})
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErrorMessages)
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

// deleteArticleById godoc
// @Summary Delete an Article
// @Description Delete an Article by article id
// @Tags Article
// @Param id  path string  true  "Article ID"
// @Success 204
// @Router /articles/{id} [delete]
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

// UpdateArticleById godoc
// @Summary Update an Article
// @Description Update an Article by article id
// @Tags Article
// @Accept  json
// @Produce  json
// @Param article body Article true "Update article"
// @Param id  path string  true  "Article ID"
// @Success 200 {object} Article
// @Failure 404 {object} CustomError
// @Router /articles/{id} [put]
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
	// Swagger
	myRouter.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The path to your swagger.json file
	))
	fmt.Printf("Server Start on port 8000\n")
	http.ListenAndServe(":8000", myRouter)
}

// @title Articles API
// @version 1.0
// @description This is a sample API for managing Articles
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email email@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /
func main() {
	initDB()
	defer closeDB()
	handleRequests()
}
