package get

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
}

type Response struct {
	response.Response
	Id          string `json:"id" validate:"id_valid,required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TaskStatus  string `json:"task_status"`
	CreatedAt   string `json:"created_at"`
	RepeatTask  string `json:"repeat_task" validate:"repeat_task_valid"`
}

type TaskGetter interface {
	GetTaskById(id uuid.UUID) (domain.Task, error)
}

// @Summary Get task by uuid
// @Description Get task by its UUID
// @Tags Task
// @Accept json
// @Produce json
// @Param request body Request true "Request"
// @Success 200 {object} Response "Task retrieved successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 500 {object} response.Response "Failed to save task"
// @Router /task [post]
func New(log *slog.Logger, taskGetter TaskGetter, rdb *redis.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.get.New"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		req := Request{
			Id: chi.URLParam(r, "id"),
		}

		validate := validator.New()
		validate.RegisterValidation("id_valid", validators.IsValidId)

		if err := validate.Struct(req); err != nil {
			log.Error("Invalid request", sl.Error(err))
			render.JSON(w, r, response.ErrorClient("Invalid request"))
			return
		}

		ctx := r.Context()
		log.Info(req.Id)
		cached, err := rdb.Get(ctx, req.Id)
		if err != nil {
			log.Info("Failed to get task from Redis", sl.Error(err))
		}
		if cached != "" {
			var task domain.Task
			if err := json.Unmarshal([]byte(cached), &task); err == nil {
				log.Info("Task retrieved from Redis", slog.String("TaskId", task.Id.String()))
				render.JSON(w, r, Response{
					Response:    response.StatusOK(),
					Id:          task.Id.String(),
					Title:       task.Title,
					Description: task.Description,
					TaskStatus:  string(task.TaskStatus),
					CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
					RepeatTask:  string(task.RepeatTask),
				})
				return
			}
		}
		log.Info("Task not found in Redis, fetching from database")

		task, err := taskGetter.GetTaskById(uuid.MustParse(req.Id))
		if err != nil {
			log.Error("Failed to get task", sl.Error(err))
			render.JSON(w, r, response.Error("Failed to get task"))
			return
		}

		taskJSON, err := json.Marshal(task)
		if err != nil {
			log.Error("Failed to marshal task", sl.Error(err))
		} else {
			if err := rdb.Set(ctx, req.Id, string(taskJSON), 5*time.Minute); err != nil {
				log.Error("Failed to set task in Redis", sl.Error(err))
			}
		}

		log.Info("Task get", slog.String("TaskId", task.Id.String()))
		render.JSON(w, r, Response{
			Response:    response.StatusOK(),
			Id:          task.Id.String(),
			Title:       task.Title,
			Description: task.Description,
			TaskStatus:  string(task.TaskStatus),
			CreatedAt:   task.CreatedAt.Format("2006-01-02 15:04:05"),
			RepeatTask:  string(task.RepeatTask),
		})
	}
}
