package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type RequestBody struct {
	Message string
}

type Response struct {
	ID      string
	Message string
}

func setupHandlers(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("POST /message", func(w http.ResponseWriter, r *http.Request) {
		var reqBody RequestBody
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		statement, err := db.Prepare("INSERT INTO messages (message) VALUES (?)")
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		result, err := statement.Exec(reqBody.Message)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Response{
			ID:      strconv.Itoa(int(id)),
			Message: reqBody.Message,
		})
	})

	mux.HandleFunc("GET /message/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}

		row := db.QueryRow("SELECT id, message FROM messages WHERE id = ?", id)
		var response Response
		err := row.Scan(&response.ID, &response.Message)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})
}

func main() {
	db, err := sql.Open("sqlite3", "messages.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	setupHandlers(mux, db)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
