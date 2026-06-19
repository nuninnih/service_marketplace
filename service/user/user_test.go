package user_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	errSvc "github.com/nuninnih/service_marketplace/service"
	user "github.com/nuninnih/service_marketplace/service/user"
	mock_user "github.com/nuninnih/service_marketplace/service/user/mock"
)

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

func TestRegister(t *testing.T) {
	tests := []struct {
		name      string
		inputUser user.User
		mockUser  func(m *mock_user.MockRepository)
		wantErr   bool
	}{
		{
			name: "error get by email",
			inputUser: user.User{
				Email:    "test@mail.com",
				Password: "123456",
			},
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().
					GetByEmail("test@mail.com").
					Return(user.User{}, errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "email already exists",
			inputUser: user.User{
				Email:    "test@mail.com",
				Password: "123456",
			},
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().
					GetByEmail("test@mail.com").
					Return(user.User{
						ID:    1,
						Email: "test@mail.com",
					}, nil)
			},
			wantErr: true,
		},
		{
			name: "error create user",
			inputUser: user.User{
				Name:     "Nunin",
				Email:    "test@mail.com",
				Password: "123456",
			},
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().
					GetByEmail("test@mail.com").
					Return(user.User{}, nil)

				m.EXPECT().
					Create(gomock.Any()).
					Return(errors.New("insert failed"))
			},
			wantErr: true,
		},
		{
			name: "success",
			inputUser: user.User{
				Name:     "Nunin",
				Email:    "test@mail.com",
				Password: "123456",
			},
			mockUser: func(m *mock_user.MockRepository) {
				m.EXPECT().
					GetByEmail("test@mail.com").
					Return(user.User{}, nil)

				m.EXPECT().
					Create(gomock.Any()).
					Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_user.NewMockRepository(ctrl)

			tt.mockUser(mockRepo)

			svc := user.NewService(
				logger,
				mockRepo,
			)

			result, err := svc.Register(tt.inputUser)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, user.User{}, result)
			} else {
				assert.NoError(t, err)

				assert.Equal(t,
					tt.inputUser.Email,
					result.Email,
				)

				assert.NotEmpty(t,
					result.Password,
				)

				assert.NotEqual(t,
					tt.inputUser.Password,
					result.Password,
				)
			}
		})
	}
}

func TestRegisterEmailExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)

	mockRepo.
		EXPECT().
		GetByEmail("test@mail.com").
		Return(user.User{
			ID:    1,
			Email: "test@mail.com",
		}, nil)

	svc := user.NewService(
		logger,
		mockRepo,
	)

	result, err := svc.Register(user.User{
		Name:     "Nunin",
		Email:    "test@mail.com",
		Password: "123456",
	})

	assert.Error(t, err)
	assert.Equal(t, errSvc.ErrEmailExists, err)
	assert.Equal(t, user.User{}, result)
}

func TestRegisterGetByEmailError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)

	expectedErr := errors.New("database error")

	mockRepo.
		EXPECT().
		GetByEmail("test@mail.com").
		Return(user.User{}, expectedErr)

	svc := user.NewService(
		logger,
		mockRepo,
	)

	result, err := svc.Register(user.User{
		Email:    "test@mail.com",
		Password: "123456",
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, user.User{}, result)
}

func TestRegisterCreateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)

	expectedErr := errors.New("insert failed")

	mockRepo.
		EXPECT().
		GetByEmail("test@mail.com").
		Return(user.User{}, nil)

	mockRepo.
		EXPECT().
		Create(gomock.Any()).
		Return(expectedErr)

	svc := user.NewService(
		logger,
		mockRepo,
	)

	result, err := svc.Register(user.User{
		Name:     "Nunin",
		Email:    "test@mail.com",
		Password: "123456",
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, user.User{}, result)
}

func TestRegisterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockRepository(ctrl)

	mockRepo.
		EXPECT().
		GetByEmail("test@mail.com").
		Return(user.User{}, nil)

	mockRepo.
		EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(u user.User) error {

			assert.Equal(t,
				"test@mail.com",
				u.Email,
			)

			assert.NotEqual(t,
				"123456",
				u.Password,
			)

			err := bcrypt.CompareHashAndPassword(
				[]byte(u.Password),
				[]byte("123456"),
			)

			assert.NoError(t, err)

			return nil
		})

	svc := user.NewService(
		logger,
		mockRepo,
	)

	result, err := svc.Register(user.User{
		Name:     "Nunin",
		Email:    "test@mail.com",
		Password: "123456",
	})

	assert.NoError(t, err)
	assert.Equal(t, "test@mail.com", result.Email)
	assert.NotEmpty(t, result.Password)
}
