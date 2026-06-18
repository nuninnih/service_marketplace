package job

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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

type allJob struct {
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
		ctrl.logger.Error("CTRL GET JOBS", slog.Any("Get all job", err))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"message": "Failed Processing Request"})
	}

	var response = []allJob{}

	for _, j := range jobs {
		response = append(response, allJob{
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

type createJobRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Budget      float64 `json:"budget" validate:"required"`
}

type createJobResponse struct {
	ID          int
	Title       string
	Description string
	Budget      float64
	Status      string
	CreatedAt   time.Time
}

func (ctrl *Controller) CreateJob(c echo.Context) error {
	id := c.Get("id").(int)

	request := new(createJobRequest)
	if err := c.Bind(request); err != nil {
		fmt.Println(err)
		ctrl.logger.Error("CTRL CREATE JOB", slog.Any("Bind", err))
		return common.CompleteErrorResponse(c, http.StatusBadRequest, "Invalid Specification")
	}

	if err := validator.New().Struct(request); err != nil {
		fmt.Println(err)
		ctrl.logger.Error("CTRL CREATE JOB", slog.Any("Validate", err))
		return common.CompleteErrorResponse(c, http.StatusBadRequest, common.ValidationErrors(err))
	}

	createdJob, err := ctrl.jobSvc.CreateJob(job.Job{
		ClientId:    id,
		Title:       request.Title,
		Description: request.Description,
		Budget:      request.Budget,
	})

	if err != nil {
		fmt.Println(err)

		if strings.Contains(err.Error(), "Exist") {
			return common.CompleteErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return common.CompleteErrorResponse(c, http.StatusInternalServerError, "Failed Processing Request")
	}

	var response = createJobResponse{
		ID:          createdJob.ID,
		Title:       createdJob.Title,
		Description: createdJob.Description,
		Budget:      createdJob.Budget,
		Status:      createdJob.Status,
		CreatedAt:   createdJob.CreatedAt,
	}

	return common.CompleteSuccessResponse(c, http.StatusCreated, response)
}
