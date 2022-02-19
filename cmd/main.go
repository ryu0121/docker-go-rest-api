package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

const (
	DB_ACCESS_PATH = "root:password@tcp(db:3306)/docker_compose_db"
)

func main() {
	e := echo.New()
	e.GET("/tasks", Index)
	e.GET("/tasks/:id", Show)
	e.POST("/tasks", Create)
	e.PUT("/tasks/:id", Update)
	e.DELETE("/tasks/:id", Destroy)
	e.Logger.Fatal(e.Start(":8080"))
}
func Index(c echo.Context) (err error) {
	conn, err := sql.Open("mysql", DB_ACCESS_PATH)
	if err != nil {
		panic(err.Error)
	}

	rows, err := conn.Query("SELECT id, title, description FROM tasks")
	if err != nil {
		return err
	}
	defer rows.Close()

	var tasks Tasks
	for rows.Next() {
		var id int
		var title string
		var description string

		if err := rows.Scan(&id, &title, &description); err != nil {
			// 一つエラーが起きても処理を止めずに最後の行まで行う
			continue
		}

		task := Task{
			ID:          id,
			Title:       title,
			Description: description,
		}
		tasks = append(tasks, task)
	}

	// 全イテレータが問題なく終了したかを確認
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, tasks)
}
func Show(c echo.Context) (err error) {
	conn, err := sql.Open("mysql", DB_ACCESS_PATH)
	if err != nil {
		panic(err.Error)
	}

	id_parameter := c.Param("id")
	var id string
	var title string
	var description string
	row := conn.QueryRow("SELECT id, title, description FROM tasks WHERE id = ?", id_parameter)
	err = row.Scan(&id, &title, &description)
	if err != nil {
		fmt.Printf("task id %s is not found", id_parameter)
		return c.String(http.StatusBadRequest, "No Record Found")
	}

	id_parsed, _ := strconv.Atoi(id)

	task := Task{
		id_parsed,
		title,
		description,
	}

	return c.JSON(http.StatusOK, task)
}

func Create(c echo.Context) (err error) {
	title := c.FormValue("Title")
	description := c.FormValue("Description")

	conn, err := sql.Open("mysql", DB_ACCESS_PATH)
	if err != nil {
		panic(err.Error)
	}
	defer conn.Close()

	result, err := conn.Exec("INSERT INTO tasks (title, description) VALUES (?, ?)", title, description)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Some Parameters Wrong")
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Something Went Wrong")
	}

	task := Task{
		int(id64),
		title,
		description,
	}

	return c.JSON(http.StatusOK, task)
}

func Update(c echo.Context) (err error) {
	conn, err := sql.Open("mysql", DB_ACCESS_PATH)
	if err != nil {
		panic(err.Error)
	}

	id_parameter := c.Param("id")
	title := c.FormValue("Title")
	description := c.FormValue("Description")
	_, err = conn.Exec("UPDATE tasks SET title = ?, description = ? WHERE id = ?", title, description, id_parameter)
	if err != nil {
		fmt.Printf("task id %s is not found", id_parameter)
		return c.String(http.StatusBadRequest, "No Record Found")
	}

	id_parsed, _ := strconv.Atoi(id_parameter)

	task := Task{
		id_parsed,
		title,
		description,
	}

	return c.JSON(http.StatusOK, task)
}

func Destroy(c echo.Context) (err error) {
	conn, err := sql.Open("mysql", DB_ACCESS_PATH)
	if err != nil {
		panic(err.Error)
	}

	id_parameter := c.Param("id")
	_, err = conn.Exec("DELETE FROM tasks WHERE id = ?", id_parameter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "not deleted")
	}

	return c.JSON(http.StatusOK, "task deleted")
}

type Task struct {
	ID          int
	Title       string
	Description string
}

type Tasks []Task
