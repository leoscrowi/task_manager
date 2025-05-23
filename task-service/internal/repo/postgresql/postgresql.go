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
        INSERT INTO tasks (id, title, description, status, created_at, repeatable)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	_, err = tx.Exec(query,
		entity.Id,
		entity.Title,
		entity.Description,
		entity.TaskStatus,
		entity.CreatedAt,
		entity.RepeatTask,
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
		SELECT id, title, description, status, created_at, repeatable
		FROM tasks
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	var task domain.Task
	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.TaskStatus,
		&task.CreatedAt,
		&task.RepeatTask,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Task{}, fmt.Errorf("%s: task not found: %w", op, err)
		}
		return domain.Task{}, fmt.Errorf("%s: failed to get task by id: %w", op, err)
	}

	return task, nil
}

func (r *Repository) UpdateTaskById(id uuid.UUID, updates domain.Task) error {
	const op = "repo.postgresql.UpdateTaskById"

	query := `
        UPDATE tasks
        SET 
            title = COALESCE($1, title),
            description = COALESCE($2, description),
            status = COALESCE($3, status),
            repeatable = COALESCE($4, repeatable)
        WHERE id = $5
    `

	result, err := r.db.Exec(query,
		sql.NullString{String: updates.Title, Valid: updates.Title != ""},
		sql.NullString{String: updates.Description, Valid: updates.Description != ""},
		sql.NullString{String: string(updates.TaskStatus), Valid: updates.TaskStatus != ""},
		sql.NullString{String: string(updates.RepeatTask), Valid: updates.RepeatTask != ""},
		id,
	)

	if err != nil {
		return fmt.Errorf("%s: failed to update task: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: failed to update task: %w", op, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: task with id %s not found", op, id)
	}

	return nil
}
