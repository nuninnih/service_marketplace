package project

import (
	"context"

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
	err = r.DB.WithContext(context.Background()).Where("client_id = ?", userId).Find(&projects).Error
	return projects, err
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
