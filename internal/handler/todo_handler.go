package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/VASANTHANUMANDLAKADI/todo/internal/model"
	"github.com/jackc/pgx/v5"
	//"github.com/jackc/pgx/v5/pgxpool"
)

var IDcount = 0
var Todos []*model.Todo

type DB interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TodoHandler struct {
	db DB
}

func (h *TodoHandler) DeleteTodo(rec *httptest.ResponseRecorder, req *http.Request) {
	panic("unimplemented")
}

func NewTodoHandler(db DB) *TodoHandler {
	return &TodoHandler{
		db: db,
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

	err = h.db.QueryRow(
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

func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {

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
	rows, err := h.db.Query(ctx, query)
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

	err = h.db.QueryRow(
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

func (h *TodoHandler) CompleteTodo(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Completed bool `json:"completed"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ctx := r.Context()

	query := `
		UPDATE todos
		SET
			completed = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING
			id,
			title,
			description,
			completed,
			COALESCE(created_at, CURRENT_TIMESTAMP),
			COALESCE(updated_at, CURRENT_TIMESTAMP)
	`

	var todo model.Todo

	err = h.db.QueryRow(
		ctx,
		query,
		request.Completed,
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

		if err == pgx.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}

		fmt.Println("Complete Todo error:", err.Error())
		http.Error(
			w,
			"Failed to update completed status: "+err.Error(),
			http.StatusInternalServerError,
		)
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

	err = h.db.QueryRow(
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
