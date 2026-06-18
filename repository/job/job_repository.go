package job

import (
	"context"

	"github.com/nuninnih/service_marketplace/service/job"
	"gorm.io/gorm"
)

type GormRepository struct {
	*gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db.Table("jobs"),
	}
}

func (r *GormRepository) GetAllJobs(input string) (jobs []job.Job, err error) {
	err = r.DB.WithContext(context.Background()).Preload("Client").
		Where("title ILIKE ? OR description ILIKE ?", "%"+input+"%", "%"+input+"%").
		Order("created_at desc").Find(&jobs).Error
	return jobs, err
}

func (r *GormRepository) GetAllJobByUser(userId int) (jobs []job.Job, err error) {
	err = r.DB.WithContext(context.Background()).
		Preload("Client").
		Where("client_id = ?", userId).Find(&jobs).Error
	return jobs, err
}

func (r *GormRepository) GetJobById(jobId int) (job job.Job, err error) {
	err = r.DB.WithContext(context.Background()).
		Preload("Client").
		Where("id = ?", jobId).First(&job).Error
	return job, err
}

func (r *GormRepository) CreateJob(input job.Job) (job job.Job, err error) {
	err = r.WithContext(context.Background()).Create(&input).Error
	return input, err
}

func (r *GormRepository) UpdateJob(input job.Job) (job job.Job, err error) {
	err = r.WithContext(context.Background()).Save(&input).Error
	return input, err
}

func (r *GormRepository) DeleteJob(jobId int) (err error) {
	return r.WithContext(context.Background()).Where("id = ?", jobId).Delete(&job.Job{}).Error
}
