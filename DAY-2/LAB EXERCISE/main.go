package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"task-app/models"
	"task-app/repository"
)

var repo *repository.TaskRepo

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		tasks, err := repo.GetAll()
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var t models.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err := repo.Save(&t); err != nil {
			http.Error(w, "Could not save task", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	// 1. Initialize SQLite connection
	db, err := sql.Open("sqlite", "tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	// 2. Create schema
	db.Exec("CREATE TABLE IF NOT EXISTS tasks (id INTEGER PRIMARY KEY AUTOINCREMENT, title TEXT, done BOOLEAN)")
	
	// 3. Dependency Injection: Pass DB to Repository
	repo = repository.NewRepo(db)

	// 4. Set up routes and start server
	http.HandleFunc("/tasks", handleTasks)
	fmt.Println("🚀 Modular Persistent API running at http://localhost:8080/tasks")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
