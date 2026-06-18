package user

import "time"

type User struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	Profile   string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
