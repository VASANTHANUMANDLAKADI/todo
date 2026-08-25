package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/VASANTHANUMANDLAKADI/todo/db"
	"github.com/VASANTHANUMANDLAKADI/todo/internal/model"
	"github.com/VASANTHANUMANDLAKADI/todo/internal/storage"
)

var IDcount = 0
var Todos []*model.Todo

type TodoHandler struct {
	storage *storage.TodoStorage
}

func NewTodoHandler(storage *storage.TodoStorage) *TodoHandler {
	return &TodoHandler{
		storage: storage,
	}
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTodoRequest


	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Println(req)

	
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	IDcount= IDcount+ 1
	todoResponse:= model.Todo{
		ID: IDcount,
		Title: req.Title,
		Description: req.Description,
		Completed: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	Todos = append(Todos, &todoResponse)

    json.NewEncoder(w).Encode(todoResponse)


	/*
	todo := model.Todo{
		Title:       req.Title,
		Description: req.Description,
	}

	createdTodo := h.storage.AddTodo(todo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(createdTodo)*/
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	//db connection
	db, err := db.ConnectDB()
	println(db,err)
	json.NewEncoder(w).Encode(Todos)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request){

	type UpdateTodo struct{
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	var req UpdateTodo

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err!= nil{
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err!= nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	
	for _, todo := range Todos {

		if todo.ID == id {
			todo.Title = req.Title
			todo.Completed = req.Completed
			todo.Description = req.Description
			todo.UpdatedAt = time.Now()

			json.NewEncoder(w).Encode(todo)
			return
		}
	}
	http.Error(w, "Todo not found", http.StatusNotFound)
}

func DeleteTask(w http.ResponseWriter, r *http.Request){

	idStr:= r.URL.Query().Get("id")

	id, err:= strconv.Atoi(idStr)
	if err!= nil{
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i:= range Todos {
		if Todos[i].ID == id {
			Todos = append(Todos[:i], Todos[i+1:]...)
			IDcount = len(Todos)
			json.NewEncoder(w).Encode(Todos)
			return
		}
	}
	http.Error(w, "Not Found", http.StatusNotFound)
}