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

type TaskReapeatType string

const (
	DAILY   TaskReapeatType = "DAILY"
	WEEKLY  TaskReapeatType = "WEEKLY"
	MONTHLY TaskReapeatType = "MONTHLY"
	YEARLY  TaskReapeatType = "YEARLY"
	NEVER   TaskReapeatType = "NEVER"
)

type Task struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	Title        string
	Description  string
	TaskStatus   TaskStatus
	CreatedAt    time.Time
	RepeatTask   TaskReapeatType
	ParentTaskId uuid.UUID
}
