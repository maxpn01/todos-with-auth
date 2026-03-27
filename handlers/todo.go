package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"todo-app-with-auth/models"
	"todo-app-with-auth/postgres"
)

func todoIDFromRequest(r *http.Request) (int64, error) {
	id := r.PathValue("id")
	if id == "" {
		return 0, errors.New("missing todo id")
	}

	parsedID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse todo id: %w", err)
	}

	return parsedID, nil
}

func (h Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	todoID, err := postgres.CreateTodo(h.DB, todo)
	if err != nil {
		log.Printf("create todo: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to create todo")
		return
	}

	writeJSON(w, http.StatusCreated, response{
		ID:      todoID,
		Message: "todo created successfully",
	})
}

func (h Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := todoIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	rowsAffected, err := postgres.UpdateTodo(h.DB, id, todo)
	if err != nil {
		log.Printf("update todo %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "failed to update todo")
		return
	}

	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}

	writeJSON(w, http.StatusOK, response{
		ID:      id,
		Message: "todo updated successfully",
	})
}

func (h Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := todoIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	rowsAffected, err := postgres.DeleteTodo(h.DB, id)
	if err != nil {
		log.Printf("delete todo %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "failed to delete todo")
		return
	}

	if rowsAffected == 0 {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}

	writeJSON(w, http.StatusOK, response{
		ID:      id,
		Message: "todo deleted successfully",
	})
}

func (h Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
	id, err := todoIDFromRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid todo id")
		return
	}

	todo, err := postgres.GetTodo(h.DB, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "todo not found")
			return
		}

		log.Printf("get todo %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "failed to get todo")
		return
	}

	writeJSON(w, http.StatusOK, todo)
}

func (h Handler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := postgres.GetTodos(h.DB)
	if err != nil {
		log.Printf("get todos: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to get todos")
		return
	}

	writeJSON(w, http.StatusOK, todos)
}
