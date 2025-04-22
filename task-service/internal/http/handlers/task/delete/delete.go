package delete

import (
	"log/slog"
	"net/http"
	"task-service/internal/http/handlers/validators"
	"task-service/internal/lib/api/response"
	"task-service/internal/lib/logger/sl"
	"task-service/internal/repo/redis"

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
	Id string `json:"id" validate:"id_valid,required"`
}

type taskDeleter interface {
	DeleteTaskById(id uuid.UUID) error
}

// @Summary Delete task by uuid
// @Description Delete task by its UUID
// @Tags Task
// @Accept json
// @Produce json
// @Param request body Request true "Request"
// @Success 200 {object} Response "Task deleted successfully"
// @Failure 400 {object} response.Response "Invalid request"
// @Failure 500 {object} response.Response "Failed to delete task"
// @Router /task [delete]
func New(log *slog.Logger, taskDeleter taskDeleter, rdb *redis.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.delete.New"

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

		err := taskDeleter.DeleteTaskById(uuid.MustParse(req.Id))
		if err != nil {
			log.Error("Failed to delete task", sl.Error(err))
			render.JSON(w, r, response.Error("Failed to delete task"))
			return
		}

		err = rdb.Delete(r.Context(), req.Id)
		if err != nil {
			log.Error("Failed to delete task from Redis", sl.Error(err))
		}
		log.Info("Task deleted from Redis", slog.String("TaskId", req.Id))

		log.Info("Task deleted", slog.String("TaskId", req.Id))

		render.JSON(w, r, Response{
			Response: response.StatusOK(),
			Id:       req.Id,
		})
	}
}
