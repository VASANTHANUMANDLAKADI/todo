package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/VASANTHANUMANDLAKADI/todo/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

func newMockTodoHandler(t *testing.T) (pgxmock.PgxPoolIface, *TodoHandler) {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatalf("failed to create mock database: %v", err)
	}

	t.Cleanup(func() {
		mock.Close()
	})

	handler := NewTodoHandler(mock)

	return mock, handler
}

func TestCreateTodo_Success(t *testing.T) {

	mock, handler := newMockTodoHandler(t)

	createdAt := time.Now()
	updatedAt := createdAt

	mock.ExpectQuery(`INSERT INTO todos`).
		WithArgs(
			"Learn Go",
			"Practice unit testing",
			false,
		).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"title",
				"description",
				"completed",
				"created_at",
				"updated_at",
			}).
				AddRow(
					1,
					"Learn Go",
					"Practice unit testing",
					false,
					createdAt,
					updatedAt,
				),
		)

	body := `{
		"title": "Learn Go",
		"description": "Practice unit testing"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/create/todo",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.CreateTodo(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			rec.Code,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"database expectations were not met: %v",
			err,
		)
	}
}

func TestCreateTodo_InvalidJSON(t *testing.T) {

	_, handler := newMockTodoHandler(t)

	req := httptest.NewRequest(
		http.MethodPost,
		"/create/todo",
		strings.NewReader(`invalid json`),
	)

	rec := httptest.NewRecorder()

	handler.CreateTodo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestCreateTodo_EmptyTitle(t *testing.T) {

	_, handler := newMockTodoHandler(t)

	body := `{
		"title": "  ",
		"description": "Something"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/create/todo",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.CreateTodo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}
}

func TestCreateTodo_DatabaseError(t *testing.T) {

	mock, handler := newMockTodoHandler(t)

	mock.ExpectQuery(`INSERT INTO todos`).
		WithArgs(
			"Learn Go",
			"Practice unit testing",
			false,
		).
		WillReturnError(errors.New("database error"))

	body := `{
        "title": "Learn Go",
        "description": "Practice unit testing"
    }`

	req := httptest.NewRequest(
		http.MethodPost,
		"/create/todo",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.CreateTodo(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"database expectations were not met: %v",
			err,
		)
	}
}

func TestGetAllTodos_Success(t *testing.T) {

	mock, handler := newMockTodoHandler(t)

	createdAt := time.Now()
	updatedAt := createdAt

	mock.ExpectQuery(`SELECT`).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"title",
				"description",
				"completed",
				"created_at",
				"updated_at",
			}).
				AddRow(
					1,
					"Learn Go",
					"Practise Go",
					false,
					createdAt,
					updatedAt,
				).
				AddRow(
					2,
					"Learn SQL",
					"Practise SQL",
					false,
					createdAt,
					updatedAt,
				),
		)

	req := httptest.NewRequest(
		http.MethodGet,
		"/get/todos",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetAllTodos(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}
	var todos []model.Todo
	err := json.NewDecoder(rec.Body).Decode(&todos)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)

		if len(todos) != 2 {
			t.Fatalf(
				"expected 2 todos, got %d",
				len(todos),
			)
		}

		if todos[0].Title != "Learn Go" {
			t.Fatalf(
				"expected first todo title %q, got %q",
				"Learn Go",
				todos[0].Title,
			)
		}

		if todos[1].Title != "Learn SQL" {
			t.Fatalf(
				"expected second todo title %q, got %q",
				"Learn SQL",
				todos[1].Title,
			)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf(
				"database expectations were not met: %v",
				err,
			)
		}
	}
}

func TestGetAllTodos_DatabaseError(t *testing.T) {

	mock, handler := newMockTodoHandler(t)

	mock.ExpectQuery(`SELECT`).
		WillReturnError(errors.New("database error"))

	req := httptest.NewRequest(
		http.MethodGet,
		"/get/todos",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetAllTodos(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"database expectations were not met: %v",
			err,
		)
	}
}

func TestGetAllTodos_ScanError(t *testing.T) {

	mock, handler := newMockTodoHandler(t)

	mock.ExpectQuery(`SELECT`).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"title",
				"description",
				"completed",
				"created_at",
				"updated_at",
			}).
				AddRow(
					"invalid-id",
					"Learn Go",
					"Practice Go",
					false,
					time.Now(),
					time.Now(),
				),
		)

	req := httptest.NewRequest(
		http.MethodGet,
		"/get/todos",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.GetAllTodos(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"database expectations were not met: %v",
			err,
		)
	}
}

func TestUpdateTodo_Success(t *testing.T) {

	mock, handler := newMockTodoHandler(t)

	createdAt := time.Now()
	updatedAt := createdAt

	mock.ExpectQuery(`UPDATE todos`).
		WithArgs(
			"Learn Go",
			"Practice unit testing",
			true,
			1,
		).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id",
				"title",
				"description",
				"completed",
				"created_at",
				"updated_at",
			}).
				AddRow(
					1,
					"Learn Go",
					"Practice unit testing",
					true,
					createdAt,
					updatedAt,
				),
		)

	body := `{
		"title": "Learn Go",
		"description": "Practice unit testing",
		"completed": true
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/update/todo?id=1",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.UpdateTodo(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf(
			"database expectations were not met: %v",
			err,
		)
	}
}

func TestUpdateTodo_InvalidJSON(t *testing.T) {

	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewTodoHandler(mock)

	reqBody := strings.NewReader(`{invalid json`)

	req := httptest.NewRequest(http.MethodPut, "/update/todo?id=1", reqBody)
	rec := httptest.NewRecorder()

	handler.UpdateTodo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}

}

func TestUpdateTodo_InvalidID(t *testing.T) {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatal(err)
	}

	defer mock.Close()

	handler := NewTodoHandler(mock)

	reqBody := strings.NewReader(`{
			"title": "Updated Todo",
			"description": "Updated description"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/update/todo?id=abc", reqBody)
	rec := httptest.NewRecorder()

	handler.UpdateTodo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestUpdateTodo_DataBaseError(t *testing.T) {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatal(err)
	}

	defer mock.Close()

	handler := NewTodoHandler(mock)

	mock.ExpectQuery("UPDATE todos").
		WithArgs("Updated Todo", "Updated description", false, 1).
		WillReturnError(errors.New("database error"))

	reqBody := strings.NewReader(`{
			"title": "Updated Todo",
			"description": "Updated description"
	}`)
	req := httptest.NewRequest(http.MethodPut, "/update/todo?id=1", reqBody)
	rec := httptest.NewRecorder()

	handler.UpdateTodo(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestUpdateTodo_EmptyTitle(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewTodoHandler(mock)

	reqBody := strings.NewReader(`{
		"title": "",
		"description": "Updated description"
	}`)

	req := httptest.NewRequest(http.MethodPut, "/update/todo?id=1", reqBody)
	rec := httptest.NewRecorder()

	handler.UpdateTodo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewTodoHandler(mock)

	rows := pgxmock.NewRows([]string{
		"id",
		"title",
		"description",
		"completed",
		"created_at",
		"updated_at",
	}).AddRow(
		1,
		"Learn Go",
		"Practice unit testing",
		false,
		time.Now(),
		time.Now(),
	)

	mock.ExpectQuery("DELETE FROM todos").
		WithArgs(1).
		WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodDelete, "/delete/todo?id=1", nil)

	rec := httptest.NewRecorder()

	handler.DeleteTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestDeleteTask_InvalidID(t *testing.T) {
	mock, err := pgxmock.NewPool()

	if err != nil {
		t.Fatal(err)
	}

	defer mock.Close()

	handler := NewTodoHandler(mock)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/delete/todo?id=abc",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.DeleteTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	defer mock.Close()
	handler := NewTodoHandler(mock)

	mock.ExpectQuery("DELETE from todos").
		WithArgs(999).
		WillReturnError(errors.New("todo not found"))

	req := httptest.NewRequest(http.MethodDelete, "/delete/todo?id=999", nil)
	rec := httptest.NewRecorder()

	handler.DeleteTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestDeleteTask_DatabaseError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()

	handler := NewTodoHandler(mock)

	mock.ExpectQuery("DELETE FROM todos").
		WithArgs(1).
		WillReturnError(errors.New("database error"))

	req := httptest.NewRequest(
		http.MethodDelete,
		"/delete/todo?id=1",
		nil,
	)

	rec := httptest.NewRecorder()

	handler.DeleteTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestCompleteTodo_Success(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}

	defer mock.Close()
	handler := NewTodoHandler(mock)

	rows := pgxmock.NewRows([]string{
		"id",
		"title",
		"description",
		"completed",
		"created_at",
		"updated_at",
	}).AddRow(
		1,
		"Learn Go",
		"Practice unit testing",
		true,
		time.Now(),
		time.Now(),
	)
	mock.ExpectQuery("UPDATE todos").
		WithArgs(true, 1).
		WillReturnRows(rows)

	reqBody := strings.NewReader(`{
		"completed" : true
		}`)

	req := httptest.NewRequest(http.MethodPut, "/completed/todo?id=1", reqBody)
	rec := httptest.NewRecorder()

	handler.CompleteTodo(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}

func TestCompleteTodo_MissingID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	handler := NewTodoHandler(mock)

	req := httptest.NewRequest(http.MethodPut, "/completed/todo", nil)
	rec := httptest.NewRecorder()

	handler.CompleteTodo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestCompleteTodo_InvalidJSON(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	handler := NewTodoHandler(mock)

	reqBody := strings.NewReader(`{
			"completed": true
	`)
	req := httptest.NewRequest(http.MethodPut, "/completed/todo?id=1", reqBody)
	rec := httptest.NewRecorder()

	handler.CompleteTodo(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestCompleteTodo_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	handler := NewTodoHandler(mock)

	mock.ExpectQuery("UPDATE todos").
		WithArgs(true, 1).
		WillReturnError(pgx.ErrNoRows)

	reqBody := strings.NewReader(`{
			"completed": true
	}`)

	req := httptest.NewRequest(http.MethodPut, "/completed/todo?id=1", reqBody)
	rec := httptest.NewRecorder()

	handler.CompleteTodo(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestCompleteTodo_DataBaseError(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mock.Close()
	handler := NewTodoHandler(mock)

	mock.ExpectQuery("UPDATE todos").
		WithArgs(true, 1).
		WillReturnError(errors.New("database error"))

	reqBody := strings.NewReader(`{
    "completed": true
	}`)

	req := httptest.NewRequest(http.MethodPut, "/completed/todo?id=1", reqBody)
	rec := httptest.NewRecorder()

	handler.CompleteTodo(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}

}
