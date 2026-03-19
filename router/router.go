package router

import (
	"database/sql"
	"net/http"
	"todo-app-with-auth/handlers"
)

func Router(db *sql.DB) http.Handler {
	mux := http.NewServeMux()
	todoHandler := handlers.TodoHandler{DB: db}

	fileServer := http.FileServer(http.Dir("./static"))

	mux.Handle("GET /", fileServer)
	mux.HandleFunc("GET /api/health", handlers.HealthHandler)
	mux.HandleFunc("GET /api/todo/{id}", todoHandler.GetTodo)
	mux.HandleFunc("GET /api/todos", todoHandler.GetTodos)
	mux.HandleFunc("POST /api/todo/create", todoHandler.CreateTodo)
	mux.HandleFunc("PUT /api/todo/update/{id}", todoHandler.UpdateTodo)
	mux.HandleFunc("DELETE /api/todo/delete/{id}", todoHandler.DeleteTodo)

	return mux
}
