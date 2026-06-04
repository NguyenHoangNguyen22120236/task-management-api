package service

import (
	"testing"

	"github.com/NguyenNH36/task-management-api/internal/mocks"
	"github.com/NguyenNH36/task-management-api/internal/api"
	"github.com/NguyenNH36/task-management-api/internal/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateTask_WithMockRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo)

	req := api.CreateTaskRequest{
		Title:       "Learn mockgen",
		Description: "Test service with mock repository",
	}

	mockRepo.
		EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(task api.Task) (api.Task, error) {
			assert.Equal(t, "1", task.Id)
			assert.Equal(t, req.Title, task.Title)
			assert.Equal(t, req.Description, task.Description)
			assert.False(t, task.Completed)

			return task, nil
		})

	task, err := taskService.CreateTask(req)

	assert.NoError(t, err)
	assert.Equal(t, "1", task.Id)
	assert.Equal(t, req.Title, task.Title)
	assert.Equal(t, req.Description, task.Description)
	assert.False(t, task.Completed)
}

func TestGetTaskByID_WhenTaskExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo)

	expectedTask := api.Task{
		Id:          "1",
		Title:       "Learn Go",
		Description: "Build API",
		Completed:   false,
	}

	mockRepo.
		EXPECT().
		FindByID("1").
		Return(expectedTask, nil)

	task, err := taskService.GetTaskByID("1")

	assert.NoError(t, err)
	assert.Equal(t, expectedTask, task)
}

func TestGetTaskByID_WhenTaskDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTaskRepository(ctrl)
	taskService := NewTaskService(mockRepo)

	mockRepo.
		EXPECT().
		FindByID("999").
		Return(api.Task{}, repository.ErrTaskNotFound)

	task, err := taskService.GetTaskByID("999")

	assert.ErrorIs(t, err, repository.ErrTaskNotFound)
	assert.Empty(t, task)
}