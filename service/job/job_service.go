package job

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
	GetAllJobs(input string) (jobs []Job, err error)
	GetJobById(id int) (jobs Job, err error)
	CreateJob(input Job) (job Job, err error)
	GetMyJob(userId int) (job []Job, err error)
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

func (s *service) GetAllJobs(input string) (jobs []Job, err error) {
	return s.repo.GetAllJobs(input)
}

func (s *service) GetJobById(id int) (jobs Job, err error) {
	getJob, err := s.repo.GetJobById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("SVC JOB BY ID", slog.Any("Data Not Found", err))
			return Job{}, errSvc.ErrDataNotFound
		}

		return Job{}, err
	}

	return getJob, err
}

func (s *service) CreateJob(input Job) (job Job, err error) {
	input.Status = "open"
	created, err := s.repo.CreateJob(input)
	if err != nil {
		s.logger.Error("SVC CREATE JOB", slog.Any("Create JOB", err))
		return
	}

	return created, err
}

func (s *service) GetMyJob(userId int) (job []Job, err error) {
	return s.repo.GetAllJobByUser(userId)
}
