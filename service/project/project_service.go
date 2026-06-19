package project

import (
	"errors"
	"log/slog"

	errSvc "github.com/nuninnih/service_marketplace/service"
	"gorm.io/gorm"
)

type service struct {
	logger *slog.Logger
	repo   Repository
}

type Service interface {
	UpdateStatusProject(userId, projectId int, status string) (err error)
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
