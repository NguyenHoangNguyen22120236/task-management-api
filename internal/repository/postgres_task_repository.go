package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/NguyenNH36/task-management-api/internal/api"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTaskRepository struct {
	db *pgxpool.Pool
}

func NewPostgresTaskRepository(db *pgxpool.Pool) *PostgresTaskRepository {
	return &PostgresTaskRepository{
		db: db,
	}
}

func (r *PostgresTaskRepository) Create(task api.Task) (api.Task, error) {
	query := `
		INSERT INTO tasks (id, title, description, completed)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		task.Id,
		task.Title,
		task.Description,
		task.Completed,
	)
	if err != nil {
		return api.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (r *PostgresTaskRepository) FindAll() ([]api.Task, error) {
	query := `
		SELECT id, title, description, completed
		FROM tasks
		ORDER BY id
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]api.Task, 0)

	for rows.Next() {
		var task api.Task

		if err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Completed,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to read task rows: %w", err)
	}

	return tasks, nil
}

func (r *PostgresTaskRepository) FindByID(id string) (api.Task, error) {
	query := `
		SELECT id, title, description, completed
		FROM tasks
		WHERE id = $1
	`

	var task api.Task

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Completed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return api.Task{}, ErrTaskNotFound
		}

		return api.Task{}, fmt.Errorf("failed to find task by id: %w", err)
	}

	return task, nil
}

func (r *PostgresTaskRepository) Update(id string, task api.Task) (api.Task, error) {
	query := `
		UPDATE tasks
		SET title = $1,
		    description = $2,
		    completed = $3,
		    updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.db.Exec(
		context.Background(),
		query,
		task.Title,
		task.Description,
		task.Completed,
		id,
	)
	if err != nil {
		return api.Task{}, fmt.Errorf("failed to update task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return api.Task{}, ErrTaskNotFound
	}

	task.Id = id
	return task, nil
}

func (r *PostgresTaskRepository) Delete(id string) error {
	query := `
		DELETE FROM tasks
		WHERE id = $1
	`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}