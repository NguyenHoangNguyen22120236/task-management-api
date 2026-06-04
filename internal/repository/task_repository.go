package repository

import (
	"errors"

	"github.com/NguyenNH36/task-management-api/internal/api"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Create(task api.Task) (api.Task, error)
	FindAll() ([]api.Task, error)
	FindByID(id string) (api.Task, error)
	Update(id string, task api.Task) (api.Task, error)
	Delete(id string) error
}