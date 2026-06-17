package job

import "time"

type Job struct {
	ID          int
	ClientId    int
	Title       string
	Description string
	Budget      int
	Status      string
	CreatedAt   time.Time
}
