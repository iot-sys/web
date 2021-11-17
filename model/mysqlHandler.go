package model

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlHandler struct {
	db *sql.DB
}

func (s *mysqlHandler) GetTodos() []*Todo {
	todos := []*Todo{}
	rows, err := s.db.Query("SELECT id, name, completed, createdAt FROM todos")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var todo Todo
		rows.Scan(&todo.ID, &todo.Name, &todo.Completed, &todo.CreatedAt)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *mysqlHandler) AddTodo(name string) *Todo {
	stmt, err := s.db.Prepare("INSERT INTO todos (name, completed, createdAt) VALUES (?, ?, now())")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(name, false)
	if err != nil {
		panic(err)
	}
	id, _ := rst.LastInsertId()
	var todo Todo
	todo.ID = int(id)
	todo.Name = name
	todo.Completed = false
	todo.CreatedAt = time.Now()
	return &todo
}

func (s *mysqlHandler) RemoveTodo(id int) bool {
	stmt, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}
	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

func (s *mysqlHandler) CompleteTodo(id int, complete bool) bool {
	stmt, err := s.db.Prepare("UPDATE todos SET completed=? WHERE id=?")
	if err != nil {
		panic(err)
	}
	rst, err := stmt.Exec(complete, id)
	if err != nil {
		panic(err)
	}
	cnt, _ := rst.RowsAffected()
	return cnt > 0
}

func (s *mysqlHandler) Close() {
	s.db.Close()
}

func newMysqlHandler(connstring string) DBHandler {
	//database, err := sql.Open("sqlite3", filepath)
	//database, err := sql.Open("mysql", "iotsys:111111@tcp(192.168.0.212:3306)/iotdb")
	database, err := sql.Open("mysql", connstring)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection OK")
	// statement, _ := database.Prepare(
	// 	`CREATE TABLE IF NOT EXISTS todos (
	// 		id        INT(10)  NOT NULL AUTO_INCREMENT PRIMARY KEY,
	// 		name      VARCHAR2(512),
	// 		completed BOOLEAN,
	// 		createdAt DATETIME
	// 	)`)
	// statement.Exec()
	_, err = database.Exec(
		`CREATE TABLE IF NOT EXISTS todos (
			id        BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name      VARCHAR(512),
			completed BOOLEAN,
			createdAt DATETIME
		)`)
	if err != nil {
		panic(err)
	}
	return &mysqlHandler{db: database}
}
