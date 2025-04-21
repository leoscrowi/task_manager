package postgresql

import (
	"database/sql"
	"fmt"
	"task-service/domain"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

// TODO: изменить название интерфейса
func NewDb(config string) (*Repository, error) {
	db, err := sql.Open("postgres", config)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) SaveTask(entity domain.Task) error {
	const op = "repo.postgresql.Save"

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("Failed to begin transaction: %w", err)
	}

	query := `
        INSERT INTO tasks (id, user_id, title, description, status, created_at, repeatable, parent_task_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	var parentTaskId interface{}
	if entity.ParentTaskId == uuid.Nil {
		parentTaskId = nil
	} else {
		parentTaskId = entity.ParentTaskId
	}

	_, err = tx.Exec(query,
		entity.Id,
		entity.UserId,
		entity.Title,
		entity.Description,
		entity.TaskStatus,
		entity.CreatedAt,
		entity.RepeatTask,
		parentTaskId,
	)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to save task: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *Repository) DeleteTaskById(id uuid.UUID) error {
	const op = "repo.postgresql.DeleteTaskById"

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: Failed to begin transaction: %w", op, err)
	}

	query := `DELETE FROM tasks WHERE id = $1`
	_, err = tx.Exec(query, id)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: failed to delete task: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return nil
}

func (r *Repository) GetTaskById(id uuid.UUID) (domain.Task, error) {
	const op = "repo.postgresql.GetTaskById"

	query := `
		SELECT id, user_id, title, description, status, created_at, repeatable, parent_task_id
		FROM tasks
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	var task domain.Task
	err := row.Scan(
		&task.Id,
		&task.UserId,
		&task.Title,
		&task.Description,
		&task.TaskStatus,
		&task.CreatedAt,
		&task.RepeatTask,
		&task.ParentTaskId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Task{}, fmt.Errorf("%s: task not found: %w", op, err)
		}
		return domain.Task{}, fmt.Errorf("%s: failed to get task by id: %w", op, err)
	}

	return task, nil
}
