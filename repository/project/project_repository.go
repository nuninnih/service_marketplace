package project

import (
	"context"
	"time"

	"github.com/nuninnih/service_marketplace/service/project"
	"gorm.io/gorm"
)

type GormRepository struct {
	*gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db.Table("projects"),
	}
}

func (r *GormRepository) GetAllProjectByUser(userId int) (projects []project.Project, err error) {
	err = r.DB.WithContext(context.Background()).Where("freelancer_id = ?", userId).Find(&projects).Error
	return projects, err
}

func (r *GormRepository) GetAllProjectDetail(projectId int) (proj project.ProjectDetail, err error) {
	err = r.DB.WithContext(context.Background()).
		Table("projects p").
		Select(`
		p.*,
		j.title AS job_title,
		c.name AS client_name,
		f.name AS freelancer_name,
		pr.bid_amount AS amount
	`).
		Joins("JOIN jobs j ON j.id = p.job_id").
		Joins("JOIN users c ON c.id = p.client_id").
		Joins("JOIN users f ON f.id = p.freelancer_id").
		Joins("JOIN proposals pr ON pr.id = p.proposal_id").
		Where("p.id = ?", projectId).
		Scan(&proj).Error
	return proj, err
}

func (r *GormRepository) GetProjectById(projectId int) (project project.Project, err error) {
	err = r.DB.WithContext(context.Background()).Where("id = ?", projectId).First(&project).Error
	return project, err
}

func (r *GormRepository) CreateProject(input project.Project) (project project.Project, err error) {
	err = r.WithContext(context.Background()).Create(&input).Error
	return input, err
}

func (r *GormRepository) UpdateProject(input project.Project) (project project.Project, err error) {
	err = r.WithContext(context.Background()).Save(&input).Error
	return input, err
}

func (r *GormRepository) DeleteProject(projectId int) (err error) {
	return r.WithContext(context.Background()).Where("id = ?", projectId).Delete(&project.Project{}).Error
}

func (r *GormRepository) PatchProject(projectId int, status string) (err error) {
	now := time.Now()
	return r.DB.WithContext(context.Background()).
		Model(&project.Project{}).
		Where("id = ?", projectId).
		Updates(map[string]interface{}{
			"status":       status,
			"submitted_at": now,
		}).Error
}

func (r *GormRepository) CompleteProject(projectId int) (err error) {
	now := time.Now()
	return r.DB.WithContext(context.Background()).
		Model(&project.Project{}).
		Where("id = ?", projectId).
		Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": now,
		}).Error
}
