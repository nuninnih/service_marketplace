package project

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	errSvc "github.com/nuninnih/service_marketplace/service"
	"gorm.io/gorm"
)

type service struct {
	logger *slog.Logger
	repo   Repository
}

type Service interface {
	UpdateStatusProject(userId, projectId int, status string) (err error)
	PayProject(userId, projectId int) (output interface{}, err error)
	CreateTransaction(projectID uint, amount int64) (*snap.Response, error)
	HandleWebhook(notificationPayload map[string]interface{}) error
}

func NewService(
	logger *slog.Logger,
	repo Repository,
) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) UpdateStatusProject(userId, projectId int, status string) (err error) {
	getProject, err := s.repo.GetProjectById(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("SVC UPDATE PROJ", slog.Any("Data Not Found", err))
			return errSvc.ErrDataNotFound
		}

		return err
	}

	if getProject.FreelancerId != userId {
		s.logger.Error("SVC UPDATE PROJ", slog.Any("Forbidden", err))
		return errSvc.ErrForbidden
	}

	err = s.repo.PatchProject(projectId, status)
	if err != nil {
		s.logger.Error("SVC UPDATE PROJ", slog.Any("Error Update status", err))
	}

	return nil
}

func (s *service) PayProject(userId, projectId int) (output interface{}, err error) {
	getProject, err := s.repo.GetAllProjectDetail(projectId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("SVC UPDATE PROJ", slog.Any("Data Not Found", err))
			return nil, errSvc.ErrDataNotFound
		}

		return nil, err
	}

	if getProject.ClientID != userId {
		s.logger.Error("SVC UPDATE PROJ", slog.Any("Forbidden", err))
		return nil, errSvc.ErrForbidden
	}

	orderID, err := s.CreateTransaction(uint(projectId), int64(getProject.Amount))
	fmt.Println(orderID)
	// err = s.repo.PatchProject(projectId, status)
	// if err != nil {
	// 	s.logger.Error("SVC UPDATE PROJ", slog.Any("Error Update status", err))
	// }

	return orderID, nil
}

func (s *service) CreateTransaction(projectID uint, amount int64) (*snap.Response, error) {
	// Inisialisasi Midtrans Client (Sebaiknya di-set di main.go/init)
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox

	orderID := fmt.Sprintf("ORDER-%d-%d", projectID, time.Now().Unix())

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
	}

	snapResp, err := snap.CreateTransaction(req)
	if err != nil {
		return nil, err
	}

	// Simpan orderID ini ke tabel project/transaksi kamu di DB lewat repo
	// s.repo.SaveOrderID(projectID, orderID)

	return snapResp, nil
}

// 2. Fungsi untuk handle Webhook dari Midtrans
func (s *service) HandleWebhook(notificationPayload map[string]interface{}) error {
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox

	// Gunakan CoreAPI client untuk memvalidasi notification payload
	c := coreapi.Client{}
	c.New(midtrans.ServerKey, midtrans.Sandbox)

	// Validasi & petakan payload ke struct Midtrans
	orderID, exists := notificationPayload["order_id"].(string)
	if !exists {
		return errors.New("invalid payload")
	}

	// Sangat direkomendasikan check langsung ke API Midtrans (menghindari penipuan/IP spoofing)
	transactionStatusResp, err := c.CheckTransaction(orderID)
	if err != nil {
		return err
	}

	// Logika penentuan status berdasarkan dokumentasi Midtrans
	if transactionStatusResp != nil {
		status := ""
		if transactionStatusResp.TransactionStatus == "capture" {
			if transactionStatusResp.FraudStatus == "accept" {
				status = "complete" // Pembayaran kartu kredit aman
			}
		} else if transactionStatusResp.TransactionStatus == "settlement" {
			status = "complete" // Pembayaran QRIS/Gopay/Transfer Bank sukses
		} else if transactionStatusResp.TransactionStatus == "deny" ||
			transactionStatusResp.TransactionStatus == "expire" ||
			transactionStatusResp.TransactionStatus == "cancel" {
			status = "failed"
		}

		parts := strings.Split(orderID, "-")

		projectID, err := strconv.Atoi(parts[1])
		if err != nil {
			return err
		}
		// Jika status berubah jadi complete, update ke DB
		if status != "" {
			err := s.repo.CompleteProject(projectID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
