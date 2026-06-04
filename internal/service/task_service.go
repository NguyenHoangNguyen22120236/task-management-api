package service

import (
	"strconv"
	"sync/atomic"

	"github.com/NguyenNH36/task-management-api/internal/api"
	"github.com/NguyenNH36/task-management-api/internal/repository"
)

type TaskService interface {
	CreateTask(req api.CreateTaskRequest) (api.Task, error)
	GetTasks() ([]api.Task, error)
	GetTaskByID(id string) (api.Task, error)
	UpdateTask(id string, req api.UpdateTaskRequest) (api.Task, error)
	DeleteTask(id string) error
}

type DefaultTaskService struct {
	repo      repository.TaskRepository
	idCounter uint64
}

func NewTaskService(repo repository.TaskRepository) *DefaultTaskService {
	return &DefaultTaskService{
		repo: repo,
	}
}

func (s *DefaultTaskService) CreateTask(req api.CreateTaskRequest) (api.Task, error) {
	id := atomic.AddUint64(&s.idCounter, 1)

	task := api.Task{
		Id:          strconv.FormatUint(id, 10),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
	}

	return s.repo.Create(task)
}

func (s *DefaultTaskService) GetTasks() ([]api.Task, error) {
	return s.repo.FindAll()
}

func (s *DefaultTaskService) GetTaskByID(id string) (api.Task, error) {
	return s.repo.FindByID(id)
}

func (s *DefaultTaskService) UpdateTask(id string, req api.UpdateTaskRequest) (api.Task, error) {
	task := api.Task{
		Id:          id,
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	}

	return s.repo.Update(id, task)
}

func (s *DefaultTaskService) DeleteTask(id string) error {
	return s.repo.Delete(id)
}