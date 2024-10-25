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

	tx := u.DB.Begin()

	if err := tx.First(&user, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
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

	tx := u.DB.Begin()

	if err := tx.Where("login = ? AND password = ?", login, password).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) SignInByEmail(email, password string) (*domain.User, error) {
	var user domain.User

	tx := u.DB.Begin()

	if err := tx.Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) SignUp(user *domain.User) error {
	tx := u.DB.Begin()

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
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

func (u UserRepository) UpdateUser(user *UpdateUserDTOInput) error {
	tx := u.DB.Begin()

	if err := tx.Where("id = ?", user.ID).Updates(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u UserRepository) GetUserByLogin(login string) (*domain.User, error) {
	var user domain.User

	tx := u.DB.Begin()

	if err := tx.Where("login = ?", login).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	tx := u.DB.Begin()

	if err := tx.Where("email = ?", email).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
