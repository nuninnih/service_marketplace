package job

import (
	"time"

	"github.com/nuninnih/service_marketplace/service/user"
)

type Job struct {
	ID          int
	ClientId    int
	Title       string
	Description string
	Budget      float64
	Status      string
	CreatedAt   time.Time
	Client      user.User `gorm:"foreignKey:ClientId"`
}
