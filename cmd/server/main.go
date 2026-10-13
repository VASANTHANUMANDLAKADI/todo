package main

import (
	"context"
	"log"
	"net/http"

	"github.com/VASANTHANUMANDLAKADI/todo/db"
	"github.com/VASANTHANUMANDLAKADI/todo/internal/handler"
	)

func main() {

	ctx := context.Background()

	// Create PostgreSQL connection pool once
	pool, err := db.ConnectDB(ctx)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Close the pool when the application exits
	defer pool.Close()

	// Create Todo handler using the shared database pool
	todoHandler := handler.NewTodoHandler(pool)

	// Register routes
	http.HandleFunc("/info", handler.InfoHandler)
	http.HandleFunc("/create/todo", todoHandler.CreateTodo)
	http.HandleFunc("/get/todos", todoHandler.GetAllTodos)
	http.HandleFunc("/update/todo", todoHandler.UpdateTodo)
	http.HandleFunc("/completed/todo", todoHandler.CompleteTodo)
	http.HandleFunc("/delete/todo", todoHandler.DeleteTask)
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	response := "Hello "+ name

	//w.Header().Set("Content-Type", "application/json")

	//json.NewEncoder(w).Encode(response)
		w.Write([]byte(response))
	})

	port := ":8080"

	log.Println("Server running on port", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}