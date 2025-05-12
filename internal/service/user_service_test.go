package service

import (
	"errors"
	"testing"

	"github.com/AZRV17/zlib-backend/internal/repository"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetUserByID(id uint) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) GetUsers() ([]*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) SignInByLogin(login, password string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) SignInByEmail(email, password string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) DeleteUser(id uint) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) UpdateUser(user *repository.UpdateUserDTOInput) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) UpdateUserRole(id uint, role domain.Role) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) SetResetPasswordToken(userID uint, token string, expiry string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) GetUserByResetToken(token string) (*domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) UpdatePassword(userID uint, password string) error {
	//TODO implement me
	panic("implement me")
}

func (m *MockUserRepo) GetUserByLogin(login string) (*domain.User, error) {
	args := m.Called(login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) SignUp(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func TestUserService_SignInByLogin(t *testing.T) {
	mockRepo := new(MockUserRepo)
	userService := NewUserService(mockRepo, nil)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	mockUser := &domain.User{
		Login:    "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetUserByLogin", "testuser").Return(mockUser, nil)

	t.Run(
		"valid login", func(t *testing.T) {
			user, err := userService.SignInByLogin("testuser", "password")
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, "testuser", user.Login)
		},
	)

	t.Run(
		"invalid password", func(t *testing.T) {
			_, err := userService.SignInByLogin("testuser", "wrongpassword")
			assert.Error(t, err)
			assert.Equal(t, "invalid password", err.Error())
		},
	)

	t.Run(
		"user not found", func(t *testing.T) {
			mockRepo.On("GetUserByLogin", "unknownuser").
				Return(nil, errors.New("user not found"))
			_, err := userService.SignInByLogin("unknownuser", "password")
			assert.Error(t, err)
			assert.Equal(t, "user not found", err.Error())
		},
	)
}

func TestUserService_SignUp(t *testing.T) {
	t.Run(
		"successful signup", func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			userService := NewUserService(mockRepo, nil)

			mockInput := &SignUpUserInput{
				Login:    "newuser",
				Password: "password",
				Email:    "newuser@example.com",
			}

			mockRepo.On("SignUp", mock.Anything).Return(nil)

			err := userService.SignUp(mockInput)
			assert.NoError(t, err)
		},
	)

	t.Run(
		"signup error", func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			userService := NewUserService(mockRepo, nil)

			mockInput := &SignUpUserInput{
				Login:    "newuser",
				Password: "password",
				Email:    "newuser@example.com",
			}

			mockRepo.On("SignUp", mock.Anything).Return(errors.New("signup failed"))

			err := userService.SignUp(mockInput)
			assert.Error(t, err)
			assert.Equal(t, "signup failed", err.Error())
		},
	)
}
