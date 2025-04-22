package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TODO        TaskStatus = "TODO"
	DONE        TaskStatus = "DONE"
	IN_PROGRESS TaskStatus = "IN_PROGRESS"
)

type TaskRepeatType string

const (
	DAILY   TaskRepeatType = "DAILY"
	WEEKLY  TaskRepeatType = "WEEKLY"
	MONTHLY TaskRepeatType = "MONTHLY"
	YEARLY  TaskRepeatType = "YEARLY"
	NEVER   TaskRepeatType = "NEVER"
)

type Task struct {
	Id          uuid.UUID
	Title       string
	Description string
	TaskStatus  TaskStatus
	CreatedAt   time.Time
	RepeatTask  TaskRepeatType
}
