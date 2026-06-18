package job

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nuninnih/service_marketplace/app/api/common"
	"github.com/nuninnih/service_marketplace/service/job"
)

type Controller struct {
	v      *validator.Validate
	logger *slog.Logger
	jobSvc job.Service
}

func NewController(
	logger *slog.Logger,
	s job.Service,
) *Controller {
	v := validator.New()
	return &Controller{
		v:      v,
		logger: logger,
		jobSvc: s,
	}
}

type jobResponse struct {
	ID          int
	Client      string
	Title       string
	Description string
	Budget      float64
	Status      string
	CreatedAt   time.Time
}

func (ctrl *Controller) GetAllJobs(c echo.Context) error {
	desc := ""
	queryDesc := c.QueryParam("desc")
	if queryDesc != "" {
		desc = queryDesc
	}

	jobs, err := ctrl.jobSvc.GetAllJobs(desc)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed Processing Request"})
	}

	var response = []jobResponse{}

	for _, j := range jobs {
		response = append(response, jobResponse{
			ID:          j.ID,
			Title:       j.Title,
			Description: j.Description,
			Client:      j.Client.Name,
			Budget:      j.Budget,
			Status:      j.Status,
			CreatedAt:   j.CreatedAt,
		})
	}

	return common.CompleteSuccessResponse(c, http.StatusOK, response)
}
