package project

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nuninnih/service_marketplace/app/api/common"
	"github.com/nuninnih/service_marketplace/service/project"
)

type Controller struct {
	v       *validator.Validate
	logger  *slog.Logger
	projSvc project.Service
}

func NewController(
	logger *slog.Logger,
	s project.Service,
) *Controller {
	v := validator.New()
	return &Controller{
		v:       v,
		logger:  logger,
		projSvc: s,
	}
}

func (ctrl *Controller) UpdateStatusProject(c echo.Context) error {
	userId := c.Get("id").(int)

	projectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		ctrl.logger.Error("CTRL UPDATE PROJ", slog.Any("Get Params", err))
		return common.CompleteErrorResponse(c, http.StatusBadRequest, "ID Should Be Number")
	}

	err = ctrl.projSvc.UpdateStatusProject(userId, projectId, "submitted")
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

	return common.CompleteSuccessResponse(c, http.StatusOK, "Status Updated")
}
