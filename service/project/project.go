package project

import (
	"time"
)

type Project struct {
	ID           int
	JobId        int
	ProposalId   int
	ClientId     int
	FreelancerId int
	Status       string
	SubmittedAt  *time.Time
	CompletedAt  *time.Time
	CreatedAt    time.Time
}

type ProjectDetail struct {
	ID           int
	JobID        int
	ProposalID   int
	ClientID     int
	FreelancerID int
	Status       string
	SubmittedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	JobTitle       string
	ClientName     string
	FreelancerName string
	Amount         float64
}
