package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type responseSucces struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"Data,omitempty"`
}

type responseError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

var tasks []Task

func main() {
	route := gin.Default()

	route.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "TODO APP",
		})
	})
	route.GET("/task", getTask)
	route.GET("/task/:id", getTaskById)

	route.Run(":8080")
}

func getTaskById(c *gin.Context) {
	id := c.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "id needs to be number",
		})
		return
	}
	for _, task := range tasks {
		if task.ID == taskId {
			c.JSON(http.StatusOK, responseSucces{
				Message: "Sucesfully get data",
				Data:    task,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, responseError{
		Message: "task not found",
	})
}

func getTask(c *gin.Context) {
	c.JSON(http.StatusOK, responseSucces{
		Message: "success",
		Data:    tasks,
	})
}

func addTask(c *gin.Context) {
	var newTask Task
	err := c.ShouldBind(newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "id needs to be a number",
		})
		return
	}

	for _, task := range tasks {
		if task.ID == newTask.ID {
			c.JSON(http.StatusBadRequest, responseError{
				Message: "Id already exist",
			})
			return
		}
	}

	tasks = append(tasks, newTask)

	c.JSON(http.StatusOK, responseSucces{
		Status:  "OK",
		Message: "Sucesfully Added New Task",
		Data:    tasks,
	})
}

func deleteAllTask(c *gin.Context) {
	tasks = []Task{}

	c.JSON(http.StatusOK, responseSucces{
		Status:  "Succes",
		Message: "All Task Deleted",
		Data:    tasks,
	})
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "id needs to be a number",
		})
		return
	}

	for index, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:index], tasks[index+1:]...)
			c.JSON(http.StatusOK, responseSucces{
				Message: "Task Deleted",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, responseError{
		Message: "Task Not Found",
	})
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Message: "id needs to be a number",
		})
		return
	}

	var newTask Task

	for index, task := range tasks {
		if task.ID == taskId {
			tasks[index].Title = newTask.Title
			tasks[index].Description = newTask.Description
			tasks[index].Status = newTask.Status
			c.JSON(http.StatusOK, responseSucces{
				Message: "Task Updated",
				Data:    task,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, responseSucces{
		Message: "Task Not Found",
	})
}
