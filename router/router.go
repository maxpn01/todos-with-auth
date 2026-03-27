package router

import (
	"database/sql"
	"net/http"
	"todo-app-with-auth/handlers"
)

func Router(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	handler := handlers.Handler{DB: db}

	fileServer := http.FileServer(http.Dir("./static"))

	mux.Handle("GET /", fileServer)

	mux.HandleFunc("GET /api/todo/{id}", handler.GetTodo)
	mux.HandleFunc("GET /api/todos", handler.GetTodos)
	mux.HandleFunc("POST /api/todo/create", handler.CreateTodo)
	mux.HandleFunc("PUT /api/todo/update/{id}", handler.UpdateTodo)
	mux.HandleFunc("DELETE /api/todo/delete/{id}", handler.DeleteTodo)

	mux.HandleFunc("POST /api/auth/signup", handler.Signup)
	mux.HandleFunc("POST /api/auth/signin", handler.Signin)
	mux.HandleFunc("POST /api/auth/signout", handler.Signout)

	return mux
}
