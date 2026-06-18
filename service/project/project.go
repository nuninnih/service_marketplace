package project

import (
	"time"

	"github.com/nuninnih/service_marketplace/service/job"
	"github.com/nuninnih/service_marketplace/service/proposal"
	"github.com/nuninnih/service_marketplace/service/user"
)

type Project struct {
	ID           int
	JobId        int
	ProposalId   int
	ClientId     int
	FreelancerId int
	Status       string
	SubmittedAt  time.Time
	CompletedAt  time.Time
	CreatedAt    time.Time

	Job        job.Job           `gorm:"foreignKey:JobId"`
	Proposal   proposal.Proposal `gorm:"foreignKey:ProposalId"`
	Client     user.User         `gorm:"foreignKey:ClientId"`
	Freelancer user.User         `gorm:"foreignKey:FreelancerId"`
}
