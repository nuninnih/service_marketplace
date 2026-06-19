package payments

import (
	"log/slog"
	"net/http"
)

type MidtransConfig struct {
	Client *http.Client
	ApiKey string
	Host   string
}

type MidtransRepository struct {
	logger         *slog.Logger
	MidtransConfig MidtransConfig
}

func NewMidtransRepository(logger *slog.Logger, config MidtransConfig) *MidtransRepository {
	return &MidtransRepository{
		logger,
		config,
	}
}
