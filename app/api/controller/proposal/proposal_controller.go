package proposal

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nuninnih/service_marketplace/app/api/common"
	"github.com/nuninnih/service_marketplace/service/proposal"
)

type Controller struct {
	v       *validator.Validate
	logger  *slog.Logger
	propSvc proposal.Service
}

func NewController(
	logger *slog.Logger,
	s proposal.Service,
) *Controller {
	v := validator.New()
	return &Controller{
		v:       v,
		logger:  logger,
		propSvc: s,
	}
}

type job struct {
	Title       string
	Description string
}

type freelancer struct {
	Name  string
	Email string
}

type proposalResponse struct {
	ID          int
	Job         job
	Freelancer  freelancer
	CoverLetter string
	BidAmount   float64
	Status      string
	CreatedAt   time.Time
}

func (ctrl *Controller) GetJobProposalPerUser(c echo.Context) error {
	id := c.Get("id").(int)

	jobId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("CTRL PROP PER USER", slog.Any("Get Params", err))
		return common.CompleteErrorResponse(c, http.StatusBadRequest, "ID Should Be Number")
	}

	proposals, err := ctrl.propSvc.GetJobProposalPerUser(id, jobId)
	if err != nil {
		fmt.Println(err)

		if strings.Contains(err.Error(), "Not Found") {
			return common.CompleteErrorResponse(c, http.StatusNotFound, err.Error())
		}

		if strings.Contains(err.Error(), "Forbidden") {
			return common.CompleteErrorResponse(c, http.StatusForbidden, err.Error())
		}

		return common.CompleteErrorResponse(c, http.StatusInternalServerError, "Failed Processing Request")
	}

	var response = []proposalResponse{}

	for _, p := range proposals {
		response = append(response, proposalResponse{
			ID: p.ID,
			Job: job{
				Title:       p.JobTitle,
				Description: p.JobDescription,
			},
			Freelancer: freelancer{
				Name:  p.FreelancerName,
				Email: p.FreelancerEmail,
			},
			CoverLetter: p.CoverLetter,
			BidAmount:   p.BidAmount,
			Status:      p.Status,
			CreatedAt:   p.CreatedAt,
		})
	}
	return common.CompleteSuccessResponse(c, http.StatusOK, response)
}
func (ctrl *Controller) CreateProposal(c echo.Context) error {
	return common.CompleteSuccessResponse(c, http.StatusCreated, "response")
}
func (ctrl *Controller) UpdateStatusProposal(c echo.Context) error {
	return common.CompleteSuccessResponse(c, http.StatusOK, "response")
}
