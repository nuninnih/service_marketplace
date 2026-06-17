package project

import "time"

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
}
