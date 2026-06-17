package user

type Repository interface {
	GetByEmail(email string) (u User, err error)
	Create(user User) (err error)
	GetById(id int) (user User, err error)
}
