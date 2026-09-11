package main

import (
	"log"
	"net/http"

	"github.com/VASANTHANUMANDLAKADI/todo/internal/handler"
    "github.com/VASANTHANUMANDLAKADI/todo/internal/storage"
)

func main() {
	// Create storage
	todoStorage := storage.NewTodoStorage()

	// Create handler
	todoHandler := handler.NewTodoHandler(todoStorage)

	// Register routes
	http.HandleFunc("/info", handler.InfoHandler)
	http.HandleFunc("/create/todo", todoHandler.CreateTodo)
	http.HandleFunc("/get/todos", handler.GetAllTodos)
	http.HandleFunc("/update/todo", todoHandler.UpdateTodo)
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

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}