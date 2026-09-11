package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

	// Decode request body
	var req model.CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(req)

	// Validate title
	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	// db connection
	db, err := db.ConnectDB()

	if err != nil {
		fmt.Println("error: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	println(db)

	// Insert query
	query := `
		INSERT INTO todos
			(title, description, completed)
		VALUES
			($1, $2, $3)
		RETURNING
			id,
			title,
			description,
			completed,
			COALESCE(created_at, CURRENT_TIMESTAMP),
			COALESCE(updated_at, CURRENT_TIMESTAMP)
	`

	ctx := r.Context()

	fmt.Println("before query", query)

	var todo model.Todo

	err = db.QueryRow(
		ctx,
		query,
		req.Title,
		req.Description,
		false,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to create todo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return created todo as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(todo)
}

func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	//db connection
	db, err := db.ConnectDB()

	if err != nil {
		fmt.Println("error: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	println(db)
	// select query to get all todo records`
	query := `
    SELECT
        id,
        title,
        description,
        completed,
        COALESCE(created_at, CURRENT_TIMESTAMP) AS created_at,
        COALESCE(updated_at, CURRENT_TIMESTAMP) AS updated_at
	FROM todos
    ORDER BY id ASC
`
	ctx := r.Context()

	//db.Ping(ctx)
	fmt.Println("before query", query)
	 rows, err := db.Query(ctx, query)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}
	 defer rows.Close()
	var todos []model.Todo

	// Read each row
	for rows.Next() {
		var todo model.Todo

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			http.Error(w, "Failed to read todo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		todos = append(todos, todo)
	}

	// Check for errors encountered while iterating
	// if err := rows.Err(); err != nil {
	// 	http.Error(w, "Failed to read rows", http.StatusInternalServerError)
	// 	return
	// }

	// Return database results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {

	type UpdateTodo struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	var req UpdateTodo

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Title) == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		fmt.Println("error:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update query
	query := `
    UPDATE todos
    SET
        title = $1,
        description = $2,
        completed = $3,
        updated_at = CURRENT_TIMESTAMP
    WHERE id = $4
    RETURNING
        id,
        title,
        description,
        completed,
        COALESCE(created_at, CURRENT_TIMESTAMP),
        COALESCE(updated_at, CURRENT_TIMESTAMP)
`

	ctx := r.Context()

	var todo model.Todo

	err = db.QueryRow(
		ctx,
		query,
		req.Title,
		req.Description,
		req.Completed,
		id,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
    fmt.Println("Update error:", err.Error())
    http.Error(w, "Failed to update todo: "+err.Error(), http.StatusInternalServerError)
    return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := db.ConnectDB()
	if err != nil {
		fmt.Println("error:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete query
	query := `
		DELETE FROM todos
		WHERE id = $1
		RETURNING
			id,
			title,
			description,
			completed,
			created_at,
			updated_at
	`

	ctx := r.Context()

	var todo model.Todo

	err = db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(todo)
}
