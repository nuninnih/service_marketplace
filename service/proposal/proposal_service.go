package proposal

import "log/slog"

type service struct {
	logger *slog.Logger
	repo   Repository
}

type Service interface {
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
