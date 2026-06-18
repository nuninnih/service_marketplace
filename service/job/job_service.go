package job

import "log/slog"

type service struct {
	logger *slog.Logger
	repo   Repository
}

type Service interface {
	GetAllJobs(input string) (jobs []Job, err error)
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
