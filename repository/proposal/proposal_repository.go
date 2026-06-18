package proposal

import (
	"context"

	"github.com/nuninnih/service_marketplace/service/proposal"
	"gorm.io/gorm"
)

type GormRepository struct {
	*gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db.Table("proposals"),
	}
}

func (r *GormRepository) GetAllProposalByUser(jobId int) (proposals []proposal.JobProposalUser, err error) {
	err = r.DB.WithContext(context.Background()).
		Table("proposals p").
		Select(`
		p.id,
		p.job_id,
		p.freelancer_id,
		p.cover_letter,
		p.bid_amount,
		p.status,
		p.created_at,

		j.title AS job_title,
		j.description AS job_description,

		u.name AS freelancer_name,
		u.email AS freelancer_email
	`).
		Joins("JOIN jobs j ON j.id = p.job_id").
		Joins("JOIN users u ON u.id = p.freelancer_id").
		Where("p.job_id = ?", jobId).
		Order("p.created_at DESC").
		Scan(&proposals).Error
	return proposals, err
}

func (r *GormRepository) GetProposalById(proposalId int) (proposal proposal.Proposal, err error) {
	err = r.DB.WithContext(context.Background()).Where("id = ?", proposalId).First(&proposal).Error
	return proposal, err
}

func (r *GormRepository) CreateProposal(input proposal.Proposal) (proposal proposal.Proposal, err error) {
	err = r.WithContext(context.Background()).Create(&input).Error
	return input, err
}

func (r *GormRepository) UpdateProposal(input proposal.Proposal) (proposal proposal.Proposal, err error) {
	err = r.WithContext(context.Background()).Save(&input).Error
	return input, err
}

func (r *GormRepository) DeleteProposal(proposalId int) (err error) {
	return r.WithContext(context.Background()).Where("id = ?", proposalId).Delete(&proposal.Proposal{}).Error
}
