package get

import (
	"log/slog"
	"net/http"
	"task-service/domain"
	"task-service/internal/lib/api/response"
	"task-service/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type Request struct {
	Id string `json:"id" validate:"id_valid,required"`
}

type Response struct {
	response.Response
	Id           string `json:"id" validate:"id_valid,required"`
	UserId       string `json:"user_id" validate:"id_valid,required"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	TaskStatus   string `json:"task_status"`
	CreatedAt    string `json:"created_at"`
	RepeatTask   string `json:"repeat_task" validate:"repeat_task_valid"`
	ParentTaskId string `json:"parent_task_id,omitempty" validate:"id_valid"`
}

type TaskGetter interface {
	GetTaskById(id uuid.UUID) (domain.Task, error)
}

func New(log *slog.Logger, taskGetter TaskGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.task.get.New"

		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		req := Request{
			Id: chi.URLParam(r, "id"),
		}

		validate := validator.New()
		validate.RegisterValidation("id_valid", IsValidId)

		if err := validate.Struct(req); err != nil {
			log.Error("Invalid request", sl.Error(err))
			render.JSON(w, r, response.Error("Invalid request"))
			return
		}

		task, err := taskGetter.GetTaskById(uuid.MustParse(req.Id))
		if err != nil {
			log.Error("Failed to get task", sl.Error(err))
			render.JSON(w, r, response.Error("Failed to get task"))
			return
		}

		log.Info("Task get", slog.String("TaskId", task.Id.String()))

		var parentId string
		if task.ParentTaskId != uuid.Nil {
			parentId = task.ParentTaskId.String()
		}

		render.JSON(w, r, Response{
			Response:     response.OK(),
			Id:           task.Id.String(),
			UserId:       task.UserId.String(),
			Title:        task.Title,
			Description:  task.Description,
			TaskStatus:   string(task.TaskStatus),
			CreatedAt:    task.CreatedAt.Format("2006-01-02 15:04:05"),
			RepeatTask:   string(task.RepeatTask),
			ParentTaskId: parentId,
		})
	}
}

// TODO: вынести в отдельный пакет валидации
func IsValidId(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	_, err := uuid.Parse(id)
	if err != nil {
		return false
	}
	return true
}
