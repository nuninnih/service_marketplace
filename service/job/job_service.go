package job

import "log/slog"

type service struct {
	logger *slog.Logger
	repo   Repository
}

type Service interface {
	GetAllJobs(input string) (jobs []Job, err error)
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
