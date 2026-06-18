package user

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	errSvc "github.com/nuninnih/service_marketplace/service"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	logger *slog.Logger
	repo   Repository
}

type Service interface {
	Register(user User) (u User, err error)
	generateToken(jwtSign string, id int, role string) (signedToken string, err error)
	Login(email string, password string) (accessToken string, err error)
	GetUser(id int) (user User, err error)
	GetAllFreelancer(desc string) (user []User, err error)
}

func NewService(
	logger *slog.Logger,
	repo Repository,
) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) Register(user User) (u User, err error) {
	getUser, err := s.repo.GetByEmail(user.Email)
	if err != nil {
		s.logger.Error("SVC REGISTER", slog.Any("Get User", err))
		return
	}

	if getUser.Email != "" {
		s.logger.Error("SVC REGISTER", slog.Any("Email Exists", err))
		err = errSvc.ErrEmailExists
		return
	}

	encPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("SVC REGISTER", slog.Any("Generate Token", err))
		err = errSvc.ErrGenerateToken
		return
	}
	fmt.Printf("REGISTER HASH: %s\n", string(encPassword))
	user.ID = getUser.ID
	user.Password = string(encPassword)

	err = s.repo.Create(user)
	if err != nil {
		s.logger.Error("SVC REGISTER", slog.Any("Create", err))
		return
	}

	return user, nil
}

func (s *service) generateToken(jwtSign string, id int, role string) (signedToken string, err error) {
	type jwtClaims struct {
		ID   int    `json:"id"`
		Role string `json:"role"`
		jwt.RegisteredClaims
	}

	timeNow := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(time.Hour * 24)),
		},
	})

	signedToken, err = token.SignedString([]byte(jwtSign))
	if err != nil {
		s.logger.Error("SVC GENERATE TOKEN", slog.Any("Sign Token", err))
		return "", errSvc.ErrInternalServer
	}

	return signedToken, nil
}

func (s *service) Login(email string, password string) (accessToken string, err error) {
	getUser, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", errSvc.ErrInternalServer
	}

	if getUser.ID == 0 {
		s.logger.Error("SVC LOGIN", slog.Any("User Not Found", err))
		return "", errSvc.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(getUser.Password), []byte(password)); err != nil {
		s.logger.Error("SVC LOGIN", slog.Any("Compare Password", err.Error()))
		return "", errSvc.ErrInvalidEmailPassword
	}

	token, err := s.generateToken(os.Getenv("JWT_SECRET"), getUser.ID, getUser.Role)
	if err != nil {
		s.logger.Error("SVC LOGIN", slog.Any("Generate Token", err.Error()))
		return "", errSvc.ErrGenerateToken
	}

	return token, nil
}

func (s *service) GetUser(id int) (user User, err error) {
	getUser, err := s.repo.GetById(id)
	if err != nil {
		fmt.Println(err)
		s.logger.Error("SVC GET USER", slog.Any("Get User", err))
		return
	}

	if getUser.ID == 0 {
		s.logger.Error("SVC GET USER", slog.Any("User Not Found", err))
		return User{}, errSvc.ErrUserNotFound
	}

	return getUser, err
}

func (s *service) GetAllFreelancer(desc string) (user []User, err error) {
	return s.repo.GetAllFreelancer(desc)
}
