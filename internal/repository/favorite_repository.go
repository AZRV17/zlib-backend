package repository

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type FavoriteRepository struct {
	DB *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{DB: db}
}

func (f FavoriteRepository) GetFavoriteByID(id uint) (*domain.Favorite, error) {
	var favorite domain.Favorite

	if err := f.DB.First(&favorite, id).Error; err != nil {
		return nil, err
	}

	return &favorite, nil
}

func (f FavoriteRepository) GetFavorites() ([]*domain.Favorite, error) {
	var favorites []*domain.Favorite

	if err := f.DB.Find(&favorites).Error; err != nil {
		return nil, err
	}

	return favorites, nil
}

func (f FavoriteRepository) CreateFavorite(favorite *domain.Favorite) error {
	tx := f.DB.Begin()

	if err := tx.Create(favorite).Error; err != nil {
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

func (f FavoriteRepository) DeleteFavorite(id uint) error {
	if err := f.DB.Delete(&domain.Favorite{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (f FavoriteRepository) GetFavoritesByUserID(id uint) ([]*domain.Favorite, error) {
	var favorites []*domain.Favorite

	tx := f.DB.Begin()

	if err := tx.Preload("Book").Preload("Book.Author").Preload("Book.Genre").Preload("Book.Publisher").Where(
		"user_id = ?",
		id,
	).Find(&favorites).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	log.Printf("Favorites: %v\n", favorites)
	log.Println(id)

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return favorites, nil
}

func (f FavoriteRepository) DeleteFavoriteByUserIDAndBookID(userID uint, bookID uint) (*domain.Favorite, error) {
	var favorite domain.Favorite

	tx := f.DB.Begin()

	if err := tx.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&favorite).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &favorite, nil
}

func (f FavoriteRepository) ExportFavoritesToCSV() ([]byte, error) {
	favorites, err := f.GetFavorites()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"ID", "Пользователь", "Книга"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, favorite := range favorites {
		row := []string{
			strconv.FormatUint(uint64(favorite.ID), 10),
			strconv.FormatUint(uint64(favorite.UserID), 10),
			strconv.FormatUint(uint64(favorite.BookID), 10),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("error flushing writer: %w", err)
	}

	return buf.Bytes(), nil
}
