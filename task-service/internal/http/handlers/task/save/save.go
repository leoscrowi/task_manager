package save

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"task-service/domain"
	"task-service/internal/http/handlers/validators"
	"task-service/internal/lib/api/response"
	"task-service/internal/lib/logger/sl"
	"task-service/internal/repo/redis"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// swagger:model
type Request struct {
	// example: Sample Task
	Title string `json:"title,omitempty"`

	// example: This is a sample task description.
	Description string `json:"description,omitempty"`

	// enum: DAILY, WEEKLY, MONTHLY, YEARLY, NEVER
	// example: DAILY
	RepeatTask string `json:"repeat_task,omitempty" validate:"repeat_task_valid"`
}

type Response struct {
	response.Response

	// example: b063de04-6fd7-41cd-8f4c-8d113e786be8
	TaskId string `json:"task_id,omitempty"`
}

type TaskSaver interface {
	SaveTask(entity domain.Task) error
}

// @Summary Create task
// @Description Create and save task
// @Tags Task
// @Accept json
// @Produce json
// @Param request body Request true "Request"
// @Success 201 {object} Response "Task created successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 500 {object} response.Response "Failed to save task"
// @Router /task [post]
func New(log *slog.Logger, taskSaver TaskSaver, rdb *redis.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.save.New"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("Failed to decode request", sl.Error(err))
			render.JSON(w, r, response.ErrorClient("Failed to decode request"))
			return
		}

		log.Info("Request decoded to JSON", slog.Any("request", req))

		validate := validator.New()
		validate.RegisterValidation("repeat_task_valid", validators.IsValidRepeatTask)

		if err := validate.Struct(req); err != nil {
			log.Error("Invalid request", sl.Error(err))
			render.JSON(w, r, response.ErrorClient("Invalid request"))
			return
		}

		task, err := CreateTask(req)
		if err != nil {
			log.Error("Invalid request", sl.Error(err))
			render.JSON(w, r, response.ErrorClient("Invalid request"))
			return
		}

		err = taskSaver.SaveTask(task)
		if err != nil {
			log.Error("Failed to save task", sl.Error(err))
			render.JSON(w, r, response.Error("Failed to save task"))
			return
		}

		ctx := r.Context()
		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Error("Failed to marshal task", sl.Error(err))
		} else {
			if err := rdb.Set(ctx, task.Id.String(), string(taskJSON), 5*time.Minute); err != nil {
				log.Error("Failed to set task in Redis", sl.Error(err))
			} else {
				log.Info("Task cached in Redis", slog.String("TaskId", task.Id.String()))
			}
		}

		log.Info("Task created successfully", slog.String("TaskId", task.Id.String()))

		render.JSON(w, r, Response{
			Response: response.StatusCreated(),
			TaskId:   task.Id.String(),
		})
		w.WriteHeader(http.StatusCreated)
	}
}

func CreateTask(req Request) (domain.Task, error) {
	if req.RepeatTask == "" {
		req.RepeatTask = "NEVER"
	}

	task := domain.Task{
		Id:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		TaskStatus:  domain.TODO,
		CreatedAt:   time.Now(),
		RepeatTask:  domain.TaskRepeatType(req.RepeatTask),
	}

	return task, nil
}
