package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"github.com/AZRV17/zlib-backend/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository   repository.UserRepo
	emailServ    *EmailService
	tokenManager *auth.TokenManager
}

func NewUserService(repo repository.UserRepo, emailService *EmailService, tokenManager *auth.TokenManager) *UserService {
	return &UserService{
		repository:   repo,
		emailServ:    emailService,
		tokenManager: tokenManager,
	}
}

// SignInByLogin - метод для авторизации по логину с JWT
func (u UserService) SignInByLogin(login, password string) (*domain.User, *auth.Tokens, error) {
	user, err := u.repository.GetUserByLogin(login)
	if err != nil {
		return nil, nil, err
	}

	if !u.comparePasswords(user.Password, password) {
		return nil, nil, fmt.Errorf("неверный пароль")
	}

	// Генерация JWT токенов
	tokens, err := u.tokenManager.GenerateTokens(user.ID, string(user.Role))
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// SignInByEmail - метод для авторизации по email с JWT
func (u UserService) SignInByEmail(email, password string) (*domain.User, *auth.Tokens, error) {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return nil, nil, err
	}

	if !u.comparePasswords(user.Password, password) {
		return nil, nil, fmt.Errorf("неверный пароль")
	}

	// Генерация JWT токенов
	tokens, err := u.tokenManager.GenerateTokens(user.ID, string(user.Role))
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// RefreshTokens - обновление JWT токенов
func (u UserService) RefreshTokens(refreshToken string) (*auth.Tokens, error) {
	return u.tokenManager.RefreshTokens(refreshToken)
}

// ParseToken - разбор JWT токена
func (u UserService) ParseToken(token string) (*auth.TokenClaims, error) {
	return u.tokenManager.ParseToken(token)
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
		FullName:       userInput.FullName,
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

	log.Println(userInput)

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

func (u UserService) UpdateUserRole(id uint, role domain.Role) error {
	return u.repository.UpdateUserRole(id, role)
}

func (u UserService) DeleteUser(id uint) error {
	return u.repository.DeleteUser(id)
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

// RequestPasswordReset запрашивает сброс пароля и отправляет ссылку на email
func (u UserService) RequestPasswordReset(email string) error {
	user, err := u.repository.GetUserByEmail(email)
	if err != nil {
		return fmt.Errorf("пользователь с таким email не найден")
	}

	// Генерируем случайный токен
	token, err := u.generateResetToken()
	if err != nil {
		return err
	}

	// Устанавливаем срок действия токена (24 часа)
	tokenExpiry := time.Now().Add(24 * time.Hour).Format(time.RFC3339)

	// Сохраняем токен в базу данных
	if err := u.repository.SetResetPasswordToken(user.ID, token, tokenExpiry); err != nil {
		return err
	}

	// Отправляем email с ссылкой для сброса пароля
	if err := u.emailServ.SendPasswordResetEmail(user.Email, token); err != nil {
		return err
	}

	return nil
}

// ValidateResetToken проверяет валидность токена для сброса пароля
func (u UserService) ValidateResetToken(token string) (*domain.User, error) {
	user, err := u.repository.GetUserByResetToken(token)
	if err != nil {
		return nil, fmt.Errorf("неверный токен сброса пароля")
	}

	// Проверяем, не истек ли срок действия токена
	if user.ResetTokenExpiry.Before(time.Now()) {
		return nil, fmt.Errorf("срок действия токена истек")
	}

	return user, nil
}

// ResetPassword сбрасывает пароль пользователя
func (u UserService) ResetPassword(token, newPassword string) error {
	user, err := u.ValidateResetToken(token)
	if err != nil {
		return err
	}

	// Хешируем новый пароль
	hashedPass, err := u.hashPassword(newPassword)
	if err != nil {
		return err
	}

	// Обновляем пароль пользователя
	if err := u.repository.UpdatePassword(user.ID, hashedPass); err != nil {
		return err
	}

	return nil
}

// generateResetToken генерирует случайный токен для сброса пароля
func (u UserService) generateResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
