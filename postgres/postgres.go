package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"todo-app-with-auth/models"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const todoSchema = `
CREATE TABLE IF NOT EXISTS todos (
	id BIGSERIAL PRIMARY KEY,
	text TEXT NOT NULL,
	is_completed BOOLEAN NOT NULL DEFAULT FALSE
);`

func New() (*sql.DB, error) {
	_ = godotenv.Load(".env")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	if _, err := db.Exec(todoSchema); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize todos schema: %w", err)
	}

	log.Println("Connected to postgres")
	return db, nil
}

func CreateTodo(db *sql.DB, todo models.Todo) (int64, error) {
	const sqlStatement = `INSERT INTO todos(text, is_completed) VALUES ($1, $2) RETURNING id`

	var id int64
	if err := db.QueryRow(sqlStatement, todo.Text, todo.IsCompleted).Scan(&id); err != nil {
		return 0, fmt.Errorf("insert todo: %w", err)
	}

	return id, nil
}

func UpdateTodo(db *sql.DB, id int64, todo models.Todo) (int64, error) {
	const sqlStatement = `UPDATE todos SET text=$2, is_completed=$3 WHERE id=$1`

	res, err := db.Exec(sqlStatement, id, todo.Text, todo.IsCompleted)
	if err != nil {
		return 0, fmt.Errorf("update todo: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("count updated rows: %w", err)
	}

	return rowsAffected, nil
}

func DeleteTodo(db *sql.DB, id int64) (int64, error) {
	const sqlStatement = `DELETE FROM todos WHERE id=$1`

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		return 0, fmt.Errorf("delete todo: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("count deleted rows: %w", err)
	}

	return rowsAffected, nil
}

func GetTodo(db *sql.DB, id int64) (models.Todo, error) {
	var todo models.Todo
	const sqlStatement = `SELECT id, text, is_completed FROM todos WHERE id=$1`

	if err := db.QueryRow(sqlStatement, id).Scan(&todo.ID, &todo.Text, &todo.IsCompleted); err != nil {
		return models.Todo{}, err
	}

	return todo, nil
}

func GetTodos(db *sql.DB) ([]models.Todo, error) {
	const sqlStatement = `SELECT id, text, is_completed FROM todos`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Text, &todo.IsCompleted); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate todos: %w", err)
	}

	return todos, nil
}
