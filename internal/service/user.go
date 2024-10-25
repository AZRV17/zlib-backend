package service

import (
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repository.UserRepo
}

func NewUserService(repo repository.UserRepo) *UserService {
	return &UserService{repository: repo}
}

func (u UserService) SignInByLogin(login, password string) (*domain.User, error) {
	user, err := u.repository.GetUserByLogin(login)
	if err != nil {
		return nil, err
	}

	if !u.comparePasswords(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (u UserService) SignInByEmail(email, password string) (*domain.User, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !u.comparePasswords(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (u UserService) SignUp(userInput *SignUpUserInput) error {
	hashedPass, err := u.hashPassword(userInput.Password)
	if err != nil {
		return err
	}

	user := &domain.User{
		Login:          userInput.Login,
		Password:       hashedPass,
		Email:          userInput.Email,
		Role:           userInput.Role,
		PhoneNumber:    userInput.PhoneNumber,
		PassportNumber: userInput.PassportNumber,
	}

	if err := u.repository.SignUp(user); err != nil {
		return err
	}

	return nil
}

func (u UserService) GetUserByID(id uint) (*domain.User, error) {
	user, err := u.repository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u UserService) GetAllUsers() ([]*domain.User, error) {
	users, err := u.repository.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserService) UpdateUser(userInput *UpdateUserInput) error {
	user := &repository.UpdateUserDTOInput{
		ID:             userInput.ID,
		Login:          userInput.Login,
		FullName:       userInput.FullName,
		Password:       userInput.Password,
		Role:           userInput.Role,
		Email:          userInput.Email,
		PhoneNumber:    userInput.PhoneNumber,
		PassportNumber: userInput.PassportNumber,
	}

	if err := u.repository.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (u UserService) GetUserByLogin(login string) (*domain.User, error) {
	return u.repository.GetUserByLogin(login)
}

func (u UserService) GetUserByEmail(email string) (*domain.User, error) {
	return u.repository.GetUserByEmail(email)
}

func (u UserService) hashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func (u UserService) comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
