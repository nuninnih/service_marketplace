package proposal

type Repository interface {
	GetAllProposalByUser(jobId int) (proposals []JobProposalUser, err error)
	GetProposalById(proposalId int) (proposal Proposal, err error)
	CreateProposal(input Proposal) (proposal Proposal, err error)
	UpdateProposal(input Proposal) (proposal Proposal, err error)
	DeleteProposal(proposalId int) (err error)
	PatchProposal(proposalId int, status string) (err error)
}
