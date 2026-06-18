package job

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
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

type MidtransNotification struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}

func (ctrl *Controller) WebhookHandler(c echo.Context) error {
	var payload MidtransNotification

	// 1. Parse the incoming JSON payload using Echo's binder
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request"})
	}

	// 2. Construct the string to be hashed
	// Formula: order_id + status_code + gross_amount + server_key
	hashString := payload.OrderID + payload.StatusCode + payload.GrossAmount + os.Getenv("MIDTRANS_SERVER_KEY")

	// 3. Generate the SHA512 hash
	hasher := sha512.New()
	hasher.Write([]byte(hashString))
	generatedSignature := hex.EncodeToString(hasher.Sum(nil))

	// 4. Verify the signature
	if generatedSignature != payload.SignatureKey {
		log.Println("⚠️ Invalid signature detected! Potential spoofing.")
		// Return 403 Forbidden
		return c.JSON(http.StatusForbidden, map[string]string{"message": "Invalid signature"})
	}

	// 5. Process the transaction based on the status
	log.Printf("✅ Signature verified for Order ID: %s\n", payload.OrderID)

	switch payload.TransactionStatus {
	case "capture":
		if payload.FraudStatus == "challenge" {
			log.Println("Payment challenged by FDS (Fraud Detection System).")
			// Update DB: status = 'challenge'
		} else if payload.FraudStatus == "accept" {
			log.Println("Payment accepted.")
			// Update DB: status = 'success'
		}
	case "settlement":
		log.Println("Payment settled successfully!")
		// Update DB: status = 'success'
	case "cancel", "deny", "expire":
		log.Printf("Payment failed/expired. Status: %s\n", payload.TransactionStatus)
		// Update DB: status = 'failed'
	case "pending":
		log.Println("Payment is pending (awaiting customer action).")
		// Update DB: status = 'pending'
	}

	// 6. Acknowledge receipt to Midtrans IMMEDIATELY
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Webhook received",
	})
}
