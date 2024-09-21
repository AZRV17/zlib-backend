package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (u UserRepository) GetUserByID(id uint) (*domain.User, error) {
	var user domain.User

	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) GetUsers() ([]*domain.User, error) {
	var users []*domain.User

	if err := u.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (u UserRepository) SignInByLogin(login, password string) (*domain.User, error) {
	var user domain.User

	if err := u.DB.Where("login = ? AND password = ?", login, password).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) SignInByEmail(email, password string) (*domain.User, error) {
	var user domain.User

	if err := u.DB.Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) SignUp(user *domain.User) error {
	if err := u.DB.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (u UserRepository) DeleteUser(id uint) error {
	if err := u.DB.Delete(&domain.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (u UserRepository) UpdateUser(user *domain.User) error {
	if err := u.DB.Save(user).Error; err != nil {
		return err
	}

	return nil
}
