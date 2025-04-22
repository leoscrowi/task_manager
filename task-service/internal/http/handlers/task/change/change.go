package change

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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// swagger:model
type Request struct {
	// example: b063de04-6fd7-41cd-8f4c-8d113e786be8
	Id string `json:"id" validate:"id_valid,required"`

	// example: New Task Title
	Title string `json:"title"`

	// example: This is a new task description.
	Description string `json:"description"`

	// enum: TODO, IN_PROGRESS, DONE
	// example: TODO
	TaskStatus string `json:"task_status" validate:"task_status_valid"`

	// enum: DAILY, WEEKLY, MONTHLY, YEARLY, NEVER
	// example: DAILY
	RepeatTask string `json:"repeat_task" validate:"repeat_task_valid"`
}

type Response struct {
	response.Response
	Id string `json:"id" validate:"id_valid,required"`
}

type TaskChanger interface {
	UpdateTaskById(id uuid.UUID, updates domain.Task) error
	GetTaskById(id uuid.UUID) (domain.Task, error)
}

// @Summary Update task by uuid
// @Description Update task by its UUID
// @Tags Task
// @Accept json
// @Produce json
// @Param request body Request true "Request"
// @Success 200 {object} Response "Task updated successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 500 {object} response.Response "Failed to update task"
// @Router /task [patch]
func New(log *slog.Logger, taskChanger TaskChanger, rdb *redis.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.change.New"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		req := Request{
			Id: chi.URLParam(r, "id"),
		}

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Failed to decode request", sl.Error(err))
			render.JSON(w, r, response.ErrorClient("Invalid request"))
			return
		}

		validate := validator.New()
		validate.RegisterValidation("id_valid", validators.IsValidId)
		validate.RegisterValidation("task_status_valid", validators.IsValidTaskStatus)
		validate.RegisterValidation("repeat_task_valid", validators.IsValidRepeatTask)

		if err := validate.Struct(req); err != nil {
			log.Error("Invalid request", sl.Error(err))
			render.JSON(w, r, response.ErrorClient("Invalid request"))
			return
		}

		updates := domain.Task{
			Title:       req.Title,
			Description: req.Description,
			TaskStatus:  domain.TaskStatus(req.TaskStatus),
			RepeatTask:  domain.TaskRepeatType(req.RepeatTask),
		}

		err := taskChanger.UpdateTaskById(uuid.MustParse(req.Id), updates)
		if err != nil {
			log.Error("Failed to update task", sl.Error(err))
			render.JSON(w, r, response.Error("Failed to update task"))
			return
		}

		updatedTask, err := taskChanger.GetTaskById(uuid.MustParse(req.Id))
		if err != nil {
			log.Error("Failed to load updated task for Redis", sl.Error(err))
		} else {
			taskJSON, err := json.Marshal(updatedTask)
			if err != nil {
				log.Error("Failed to marshal updated task", sl.Error(err))
			} else {
				cacheKey := "task:" + req.Id
				if err := rdb.Set(r.Context(), cacheKey, string(taskJSON), 5*time.Minute); err != nil {
					log.Error("Failed to update task in Redis", sl.Error(err))
				} else {
					log.Info("Task updated in Redis cache", slog.String("TaskId", req.Id))
				}
			}
		}

		log.Info("Task updated successfully", slog.String("TaskId", req.Id))

		render.JSON(w, r, Response{
			Response: response.StatusOK(),
			Id:       req.Id,
		})
	}
}
