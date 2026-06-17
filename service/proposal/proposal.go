package proposal

import "time"

type Proposal struct {
	ID           int
	JobId        int
	FreelancerId int
	CoverLetter  string
	BidAmount    int
	Status       string
	CreatedAt    time.Time
}
