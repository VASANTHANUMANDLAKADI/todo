package storage

import (
	"sync"

	"time"

	"github.com/VASANTHANUMANDLAKADI/todo/internal/model"
)

type TodoStorage struct{
	mu sync.Mutex
	todos []model.Todo
	nextID int
}

func NewTodoStorage() *TodoStorage{
	return &TodoStorage{
		todos: []model.Todo{},
		nextID: 1,
	}
}

func (s *TodoStorage) AddTodo(todo model.Todo) model.Todo{
	s.mu.Lock()
	defer s.mu.Unlock()

	todo.ID = s.nextID

	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	s.todos = append(s.todos, todo)
	s.nextID++
	return todo
}

func (s *TodoStorage) GetTodoById(id int) (model.Todo, bool){
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, todo:= range s.todos{
		if todo.ID == id{
			return todo, true
		}
	}
	return model.Todo{}, false
}

func (s *TodoStorage) GetAllTodos() []model.Todo{
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.todos
}

func (s *TodoStorage) UpdateTodo(id int, updatedTodo model.Todo) (model.Todo, bool){
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, todo := range s.todos{
		if todo.ID == id {
			updatedTodo.ID = todo.ID
			updatedTodo.CreatedAt = todo.CreatedAt
			updatedTodo.UpdatedAt = time.Now()

			s.todos[i] = updatedTodo

			return updatedTodo, true
		}
	}
	return model.Todo{}, false
}

func (s *TodoStorage) DeleteTodo(id int) bool{
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, todo := range s.todos {
		if todo.ID == id {
			s.todos = append(s.todos[:i],s.todos[i+1:]...)
			return true
		}
	}
	return false
}