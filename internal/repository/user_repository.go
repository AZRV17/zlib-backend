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

	tx := u.DB.Begin()

	if err := tx.Find(&users).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return users, nil
}

func (u UserRepository) SignInByLogin(login, password string) (*domain.User, error) {
	var user domain.User

	tx := u.DB.Begin()

	if err := tx.Raw("SELECT * FROM sign_in(?, ?)", login, password).
		Scan(&user).Error; err != nil {
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
	tx := u.DB.Begin()

	if err := tx.Delete(&domain.User{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
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

func (u UserRepository) UpdateUserRole(id uint, role domain.Role) error {
	tx := u.DB.Begin()

	if err := tx.Where("id = ?", id).Update("role", role).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u UserRepository) SetResetPasswordToken(userID uint, token string, expiry string) error {
	tx := u.DB.Begin()

	if err := tx.Model(&domain.User{}).Where("id = ?", userID).
		Updates(
			map[string]interface{}{
				"reset_password_token": token,
				"reset_token_expiry":   expiry,
			},
		).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (u UserRepository) GetUserByResetToken(token string) (*domain.User, error) {
	var user domain.User

	tx := u.DB.Begin()

	if err := tx.Where("reset_password_token = ?", token).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u UserRepository) UpdatePassword(userID uint, password string) error {
	tx := u.DB.Begin()

	if err := tx.Model(&domain.User{}).Where("id = ?", userID).
		Updates(
			map[string]interface{}{
				"password":             password,
				"reset_password_token": "",
				"reset_token_expiry":   nil,
			},
		).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}
