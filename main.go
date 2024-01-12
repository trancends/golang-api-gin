package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type responseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type responseError struct {
	Error string `json:"error"`
}

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
}

var tasks []Task

var (
	username string = "enigma"
	password string = "rahasia"
)

func BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()

		if !ok || user != username || pass != password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseError{
				Error: "Unauthorized",
			})
			return
		}

		c.Next()
	}
}

func main() {
	route := gin.Default()

	route.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "This is TODO APP",
		})
	})

	route.Use(BasicAuthMiddleware())
	route.GET("/task", getTask)
	route.GET("/task/:id", getTaskById)
	route.POST("/task/add", addTask)

	route.PUT("/task/:id", updateTask)
	route.DELETE("/task/delete/all", deleteAllTask)
	route.DELETE("/task/delete/:id", deleteTask)

	err := route.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getTaskById(c *gin.Context) {
	id := c.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Error: "id needs to be number",
		})
		return
	}
	for _, task := range tasks {
		if task.ID == taskId {
			c.JSON(http.StatusOK, responseSuccess{
				Message: "Successfully get data",
				Data:    task,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, responseError{
		Error: "task not found",
	})
}

func getTask(c *gin.Context) {
	c.JSON(http.StatusOK, responseSuccess{
		Message: "success",
		Data:    tasks,
	})
}

func addTask(c *gin.Context) {
	var newTask Task
	err := c.ShouldBind(&newTask)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Error: "id needs to be a number",
		})
		return
	}

	for _, task := range tasks {
		if task.ID == newTask.ID {
			c.JSON(http.StatusBadRequest, responseError{
				Error: "id already exist",
			})
			return
		}
	}

	tasks = append(tasks, newTask)

	c.JSON(http.StatusOK, responseSuccess{
		Message: "Sucesfully Added New Task",
		Data:    tasks,
	})
}

func deleteAllTask(c *gin.Context) {
	tasks = []Task{}

	c.JSON(http.StatusOK, responseSuccess{
		Message: "All Task Deleted",
		Data:    tasks,
	})
}

func deleteTask(c *gin.Context) {
	id := c.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Error: "id needs to be a number",
		})
		return
	}

	for index, task := range tasks {
		if task.ID == taskId {
			tasks = append(tasks[:index], tasks[index+1:]...)
			c.JSON(http.StatusOK, responseSuccess{
				Message: "Task Deleted",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, responseError{
		Error: "Task Not Found",
	})
}

func updateTask(c *gin.Context) {
	id := c.Param("id")
	taskId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, responseError{
			Error: "id needs to be a number",
		})
		return
	}

	var newTask Task
	c.ShouldBind(&newTask)

	for index, task := range tasks {
		if task.ID == taskId {
			tasks[index].Title = newTask.Title
			tasks[index].Description = newTask.Description
			tasks[index].Status = newTask.Status
			c.JSON(http.StatusOK, responseSuccess{
				Message: "Task Updated",
				Data:    task,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, responseSuccess{
		Message: "Task Not Found",
	})
}
