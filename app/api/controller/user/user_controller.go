package user

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/nuninnih/service_marketplace/app/api/common"
	"github.com/nuninnih/service_marketplace/service/user"
)

type Controller struct {
	v       *validator.Validate
	logger  *slog.Logger
	userSvc user.Service
}

func NewController(
	logger *slog.Logger,
	s user.Service,
) *Controller {
	v := validator.New()
	return &Controller{
		v:       v,
		logger:  logger,
		userSvc: s,
	}
}

type userRegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required"`
}

type userRegisterResponse struct {
	Name  string
	Email string
	Role  string
}

func (ctrl *Controller) Register(c echo.Context) error {
	request := new(userRegisterRequest)
	if err := c.Bind(request); err != nil {
		fmt.Println(err)
		ctrl.logger.Error("CTRL REGISTER", slog.Any("Bind", err))
		return common.CompleteErrorResponse(c, http.StatusBadRequest, "Invalid Specification")
	}

	if err := validator.New().Struct(request); err != nil {
		fmt.Println(err)
		ctrl.logger.Error("CTRL REGISTER", slog.Any("Validate", err))
		return common.CompleteErrorResponse(c, http.StatusBadRequest, common.ValidationErrors(err))
	}

	createdUser, err := ctrl.userSvc.Register(user.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
	})

	response := userRegisterResponse{
		Name:  createdUser.Name,
		Email: createdUser.Email,
		Role:  createdUser.Role,
	}

	if err != nil {
		fmt.Println(err)

		if strings.Contains(err.Error(), "Exist") {
			return common.CompleteErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return common.CompleteErrorResponse(c, http.StatusInternalServerError, "Failed Processing Request")
	}

	return common.CompleteSuccessResponse(c, http.StatusCreated, response)
}

type userLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (ctrl *Controller) Login(c echo.Context) error {
	request := new(userLoginRequest)
	err := c.Bind(request)
	if err != nil {
		fmt.Println(err)
		ctrl.logger.Error("CTRL LOGIN", slog.Any("Bind", err))
		return common.CompleteErrorResponse(c, http.StatusBadRequest, "Invalid Specification")
	}

	err = ctrl.v.Struct(request)
	if err != nil {
		fmt.Println(err)
		ctrl.logger.Error("CTRL LOGIN", slog.Any("Validate", err))
		return common.CompleteErrorResponse(c, http.StatusBadRequest, common.ValidationErrors(err))
	}

	accessToken, err := ctrl.userSvc.Login(request.Email, request.Password)
	if err != nil {
		return common.CompleteErrorResponse(c, http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Login Success", "access_token": accessToken})
}

type userResponse struct {
	ID    int
	Name  string
	Email string
}

func (ctrl *Controller) GetUser(c echo.Context) error {
	id := c.Get("id").(int)

	user, err := ctrl.userSvc.GetUser(id)
	if err != nil {
		return common.CompleteErrorResponse(c, http.StatusInternalServerError, "Failed Processing Request")
	}

	response := userResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return common.CompleteSuccessResponse(c, http.StatusOK, response)
}
