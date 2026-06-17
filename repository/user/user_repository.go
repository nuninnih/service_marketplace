package user

import (
	"context"

	"github.com/nuninnih/service_marketplace/service/user"
	"gorm.io/gorm"
)

type GormRepository struct {
	*gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db.Table("users"),
	}
}

func (r *GormRepository) GetByEmail(email string) (u user.User, err error) {
	r.DB.WithContext(context.Background()).Where("email = ?", email).First(&u)
	return u, err
}

func (r *GormRepository) Create(user user.User) (err error) {
	return r.DB.WithContext(context.Background()).Create(&user).Error
}

func (r *GormRepository) GetById(id int) (user user.User, err error) {
	r.DB.WithContext(context.Background()).Where("id = ?", id).First(&user)
	return user, err
}
