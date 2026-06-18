package proposal

import (
	"time"
)

type Proposal struct {
	ID           int
	JobId        int
	FreelancerId int
	CoverLetter  string
	BidAmount    int
	Status       string
	CreatedAt    time.Time
}

type JobProposalUser struct {
	ID           int
	JobID        int
	FreelancerID int
	CoverLetter  string
	BidAmount    float64
	Status       string
	CreatedAt    time.Time

	JobTitle       string
	JobDescription string

	FreelancerName  string
	FreelancerEmail string
}
