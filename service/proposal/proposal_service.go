package proposal

import (
	"errors"
	"log/slog"

	errSvc "github.com/nuninnih/service_marketplace/service"
	"github.com/nuninnih/service_marketplace/service/job"
	"gorm.io/gorm"
)

type service struct {
	logger  *slog.Logger
	repo    Repository
	jobRepo job.Repository
}

type Service interface {
	GetJobProposalPerUser(userId, jobId int) (proposals []JobProposalUser, err error)
	CreateProposal(input Proposal) (proposal Proposal, err error)
	UpdateStatusProposal(proposalId int, status string) (proposal Proposal, err error)
}

func NewService(
	logger *slog.Logger,
	repo Repository,
	jobRepo job.Repository,
) Service {
	return &service{
		logger:  logger,
		repo:    repo,
		jobRepo: jobRepo,
	}
}

func (s *service) GetJobProposalPerUser(userId, jobId int) (proposals []JobProposalUser, err error) {
	getJob, err := s.jobRepo.GetJobById(jobId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("SVC PROP PER USER", slog.Any("Data Not Found", err))
			return nil, errSvc.ErrDataNotFound
		}

		return nil, err
	}

	if getJob.ClientId != userId {
		s.logger.Error("SVC PROP PER USER", slog.Any("Forbidden", err))
		return nil, errSvc.ErrForbidden
	}

	getProposal, err := s.repo.GetAllProposalByUser(jobId)
	if err != nil {
		s.logger.Error("SVC PROP PER USER", slog.Any("Get PER USER", err))
		return
	}

	return getProposal, err
}

func (s *service) CreateProposal(input Proposal) (proposal Proposal, err error) {
	getJob, err := s.jobRepo.GetJobById(input.JobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("SVC CREATE PROP", slog.Any("Data Not Found", err))
			return Proposal{}, errSvc.ErrDataNotFound
		}

		return Proposal{}, err
	}

	if getJob.Status == "closed" {
		s.logger.Error("SVC CREATE PROP", slog.Any("Job Closed", err))
		return Proposal{}, errSvc.ErrClosed
	}

	input.Status = "pending"
	createProposal, err := s.repo.CreateProposal(input)
	if err != nil {
		s.logger.Error("SVC CREATE PROP", slog.Any("Create prop", err))
		return
	}

	return createProposal, err
}

func (s *service) UpdateStatusProposal(proposalId int, status string) (proposal Proposal, err error) {
	return
}
