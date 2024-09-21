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

func (r ReviewRepository) GetReviewByID(id uint) (*domain.Review, error) {
	var review domain.Review

	if err := r.DB.First(&review, id).Error; err != nil {
		return nil, err
	}

	return &review, nil
}

func (r ReviewRepository) GetReviews() ([]*domain.Review, error) {
	var reviews []*domain.Review

	if err := r.DB.Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}

func (r ReviewRepository) CreateReview(review *domain.Review) error {
	if err := r.DB.Create(review).Error; err != nil {
		return err
	}

	return nil
}

func (r ReviewRepository) UpdateReview(review *domain.Review) error {
	if err := r.DB.Save(review).Error; err != nil {
		return err
	}

	return nil
}

func (r ReviewRepository) DeleteReview(id uint) error {
	if err := r.DB.Delete(&domain.Review{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r ReviewRepository) GetReviewsByBookID(id uint) ([]*domain.Review, error) {
	var reviews []*domain.Review

	if err := r.DB.Where("book_id = ?", id).Find(&reviews).Error; err != nil {
		return nil, err
	}

	return reviews, nil
}
