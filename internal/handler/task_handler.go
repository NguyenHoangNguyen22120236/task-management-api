package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/NguyenNH36/task-management-api/internal/api"
	"github.com/NguyenNH36/task-management-api/internal/repository"
	"github.com/NguyenNH36/task-management-api/internal/service"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, api.HealthResponse{
		Status: "ok",
	})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req api.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "invalid request body",
			Code:    "INVALID_REQUEST_BODY",
		})
		return
	}

	task, err := h.service.CreateTask(api.CreateTaskRequest{
		Title:       req.Title,
		Description: valueOrEmpty(req.Description),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "internal server error",
			Code:    "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusCreated, toAPITask(task))
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.service.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "internal server error",
			Code:    "INTERNAL_SERVER_ERROR",
		})
		return
	}

	apiTasks := make([]api.Task, 0, len(tasks))
	for _, task := range tasks {
		apiTasks = append(apiTasks, toAPITask(task))
	}

	c.JSON(http.StatusOK, apiTasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context, id string) {
	task, err := h.service.GetTaskByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse{
				Message: "task not found",
				Code:    "TASK_NOT_FOUND",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "internal server error",
			Code:    "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, toAPITask(task))
}

func (h *TaskHandler) UpdateTask(c *gin.Context, id string) {
	var req api.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{
			Message: "invalid request body",
			Code:    "INVALID_REQUEST_BODY",
		})
		return
	}

	task, err := h.service.UpdateTask(id, api.UpdateTaskRequest{
		Title:       valueOrEmpty(req.Title),
		Description: valueOrEmpty(req.Description),
		Completed:   valueOrFalse(req.Completed),
	})
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse{
				Message: "task not found",
				Code:    "TASK_NOT_FOUND",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "internal server error",
			Code:    "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.JSON(http.StatusOK, toAPITask(task))
}

func (h *TaskHandler) DeleteTask(c *gin.Context, id string) {
	err := h.service.DeleteTask(id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, api.ErrorResponse{
				Message: "task not found",
				Code:    "TASK_NOT_FOUND",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, api.ErrorResponse{
			Message: "internal server error",
			Code:    "INTERNAL_SERVER_ERROR",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func toAPITask(task api.Task) api.Task {
	return api.Task{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Completed:   task.Completed,
	}
}

func valueOrEmpty(value string) string {
	if value == "" {
		return ""
	}

	return value
}

func valueOrFalse(value bool) bool {
	if !value {
		return false
	}

	return value
}