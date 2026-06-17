package proposal

type Repository interface {
	GetAllProjectByUser(userId int) (proposals []Proposal, err error)
	GetProposalById(projectId int) (proposal Proposal, err error)
	CreateProposal(input Proposal) (proposal Proposal, err error)
	UpdateProposal(input Proposal) (proposal Proposal, err error)
	DeleteProposal(projectId int) (err error)
}
