package save

import (
	"log/slog"
	"net/http"
	"task-service/domain"
	"task-service/internal/lib/api/response"
	"task-service/internal/lib/logger/sl"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Request struct {
	UserId       string `json:"user_id" validate:"id_valid,required"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	RepeatTask   string `json:"repeat_task" validate:"repeat_task_valid"`
	ParentTaskId string `json:"parent_task_id,omitempty" validate:"id_valid"`
}

type Response struct {
	response.Response
	TaskId string `json:"task_id,omitempty"`
}

type TaskSaver interface {
	SaveTask(entity domain.Task) error
}

func New(log *slog.Logger, taskSaver TaskSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.save.New"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request", sl.Error(err))
			render.JSON(w, r, response.Error)
			return
		}

		log.Info("Request decoded to JSON", slog.Any("request", req))

		validate := validator.New()
		validate.RegisterValidation("repeat_task_valid", IsValidRepeatTask)
		validate.RegisterValidation("id_valid", IsValidId)

		if err := validate.Struct(req); err != nil {
			log.Error("Invalid request", sl.Error(err))
			render.JSON(w, r, response.Error("Invalid request"))
			return
		}

		task, err := CreateTask(req)
		if err != nil {
			log.Error("Invalid request", sl.Error(err))
			render.JSON(w, r, response.Error("Invalid request"))
			return
		}

		err = taskSaver.SaveTask(task)
		if err != nil {
			log.Error("Failed to save task", sl.Error(err))
			render.JSON(w, r, response.Error("Failed to save task"))
			return
		}

		log.Info("Task added", slog.String("TaskId", task.Id.String()))

		render.JSON(w, r, Response{
			Response: response.OK(),
			TaskId:   task.Id.String(),
		})
	}
}

func CreateTask(req Request) (domain.Task, error) {
	task := domain.Task{
		Id:           uuid.New(),
		UserId:       uuid.MustParse(req.UserId),
		Title:        req.Title,
		Description:  req.Description,
		TaskStatus:   domain.TODO,
		CreatedAt:    time.Now(),
		RepeatTask:   domain.TaskReapeatType(req.RepeatTask),
		ParentTaskId: uuid.Nil,
	}

	return task, nil
}

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
