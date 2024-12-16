package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	DB *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{DB: db}
}

func (r *ReviewRepository) GetReviewByID(id uint) (*domain.Review, error) {
	var review domain.Review

	tx := r.DB.Begin()

	if err := tx.First(&review, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &review, nil
}

func (r *ReviewRepository) GetReviews() ([]*domain.Review, error) {
	var reviews []*domain.Review

	tx := r.DB.Begin()

	if err := tx.Find(&reviews).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return reviews, nil
}

func (r *ReviewRepository) CreateReview(review *domain.Review) error {
	tx := r.DB.Begin()

	if err := tx.Create(review).Error; err != nil {
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

func (r *ReviewRepository) UpdateReview(review *domain.Review) error {
	tx := r.DB.Begin()

	if err := tx.Where("id = ?", review.ID).Save(review).Error; err != nil {
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

func (r *ReviewRepository) DeleteReview(id uint) error {
	tx := r.DB.Begin()

	if err := tx.Delete(&domain.Review{}, id).Error; err != nil {
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

func (r ReviewRepository) GetReviewsByBookID(id uint) ([]*domain.Review, error) {
	var reviews []*domain.Review

	tx := r.DB.Begin()

	if err := tx.Preload("User").Where("book_id = ?", id).Find(&reviews).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return reviews, nil
}
