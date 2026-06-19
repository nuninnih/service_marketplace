package proposal

import (
	"errors"
	"log/slog"

	errSvc "github.com/nuninnih/service_marketplace/service"
	"github.com/nuninnih/service_marketplace/service/job"
	"github.com/nuninnih/service_marketplace/service/project"
	"gorm.io/gorm"
)

type service struct {
	logger   *slog.Logger
	repo     Repository
	jobRepo  job.Repository
	projRepo project.Repository
}

type Service interface {
	GetJobProposalPerUser(userId, jobId int) (proposals []JobProposalUser, err error)
	CreateProposal(input Proposal) (proposal Proposal, err error)
	UpdateStatusProposal(userId, proposalId int, status string) (proj project.Project, err error)
}

func NewService(
	logger *slog.Logger,
	repo Repository,
	jobRepo job.Repository,
	projRepo project.Repository,
) Service {
	return &service{
		logger:   logger,
		repo:     repo,
		jobRepo:  jobRepo,
		projRepo: projRepo,
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

func (s *service) GetProposalById(proposalId int) (proposal Proposal, err error) {
	getProposal, err := s.repo.GetProposalById(proposalId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("SVC GET PROP BY ID", slog.Any("Data Not Found", err))
			return Proposal{}, errSvc.ErrDataNotFound
		}

		return Proposal{}, err
	}

	return getProposal, err
}

func (s *service) UpdateStatusProposal(userId, proposalId int, status string) (proj project.Project, err error) {
	getProposal, err := s.repo.GetProposalById(proposalId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("SVC UPDATE PROP", slog.Any("Data Not Found", err))
			return project.Project{}, errSvc.ErrDataNotFound
		}

		return project.Project{}, err
	}

	if getProposal.Status == "accepted" {
		return project.Project{}, errSvc.ErrAccepted
	}

	getJob, err := s.jobRepo.GetJobById(getProposal.JobID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error("SVC UPDATE PROP", slog.Any("Data Not Found", err))
			return project.Project{}, errSvc.ErrDataNotFound
		}

		return project.Project{}, err
	}

	if getJob.ClientId != userId {
		s.logger.Error("SVC UPDATE PROP", slog.Any("Forbidden", err))
		return project.Project{}, errSvc.ErrForbidden
	}

	// Proposal -> accepted
	// Job -> closed
	// Project -> created

	err = s.repo.PatchProposal(proposalId, status)
	if err != nil {
		s.logger.Error("SVC UPDATE PROP", slog.Any("Error Patch Status Prop", err))
		return
	}

	err = s.jobRepo.PatchJob(getJob.ID, "closed")
	if err != nil {
		s.logger.Error("SVC UPDATE PROP", slog.Any("Error Patch Status Job", err))
		return
	}

	outputProject, err := s.projRepo.CreateProject(project.Project{
		JobId:        getJob.ID,
		ProposalId:   proposalId,
		ClientId:     userId,
		FreelancerId: getProposal.FreelancerID,
		Status:       "in_progress",
	})

	if err != nil {
		s.logger.Error("SVC UPDATE PROP", slog.Any("Error create project", err))
		return
	}

	return outputProject, err
}
