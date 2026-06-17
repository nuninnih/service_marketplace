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

func (r *GormRepository) GetAllProjectByUser(userId int) (proposals []proposal.Proposal, err error) {
	err = r.DB.WithContext(context.Background()).Where("client_id = ?", userId).Find(&proposals).Error
	return proposals, err
}

func (r *GormRepository) GetProposalById(projectId int) (proposal proposal.Proposal, err error) {
	err = r.DB.WithContext(context.Background()).Where("id = ?", projectId).First(&proposal).Error
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

func (r *GormRepository) DeleteProposal(projectId int) (err error) {
	return r.WithContext(context.Background()).Where("id = ?", projectId).Delete(&proposal.Proposal{}).Error
}
