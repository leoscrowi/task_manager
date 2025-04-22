package validators

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

func IsValidRepeatTask(fl validator.FieldLevel) bool {
	repeat := fl.Field().String()
	switch repeat {
	case "DAILY", "WEEKLY", "MONTHLY", "YEARLY", "NEVER":
		return true
	default:
		return false
	}
}

func IsValidId(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	_, err := uuid.Parse(id)
	if err != nil {
		return false
	}
	return true
}

func IsValidTaskStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	switch status {
	case "TODO", "IN_PROGRESS", "DONE":
		return true
	default:
		return false
	}
}
