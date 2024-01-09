package main

import (
	"encoding/json"
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
		c.JSON(http.StatusInternalServerError, responseError{
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

func addTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{
			Status:  "Error",
			Message: "Failed to decode json",
		})
		return
	}

	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(responseSucces{
		Status:  "OK",
		Message: "Sucesfully Added New Task",
		Data:    tasks,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func deleteAllTask(w http.ResponseWriter, r *http.Request) {
	tasks = []Task{}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseSucces{
		Status:  "Succes",
		Message: "All Task Deleted",
		Data:    tasks,
	})
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{
			Status:  "Error",
			Message: "Failed to decode json",
		})
		return
	}

	for index, task := range tasks {
		if task.ID == newTask.ID {
			tasks = append(tasks[:index], tasks[index+1:]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseSucces{
				Status:  "Succes",
				Message: "Task Deleted",
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(responseError{
		Status:  "error",
		Message: "Task Not Found",
	})
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(responseError{Message: "Metod Not Allowed", Status: "error"})
		return
	}

	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(responseError{
			Status:  "Error",
			Message: "Failed to decode json",
		})
		return
	}

	for index, task := range tasks {
		if task.ID == newTask.ID {
			tasks[index].Title = newTask.Title
			tasks[index].Description = newTask.Description
			tasks[index].Status = newTask.Status
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseSucces{
				Status:  "success",
				Message: "Task Updated",
				Data:    tasks,
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(responseSucces{
		Status:  "error",
		Message: "Task Not Found",
	})
}
