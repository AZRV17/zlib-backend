package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
)

type ReviewService struct {
	repository     repository.ReviewRepo
	bookRepository repository.BookRepo
}

func NewReviewService(repo repository.ReviewRepo) *ReviewService {
	return &ReviewService{repository: repo}
}

func (r ReviewService) GetReviewByID(id uint) (*domain.Review, error) {
	return r.repository.GetReviewByID(id)
}

func (r ReviewService) GetReviews() ([]*domain.Review, error) {
	return r.repository.GetReviews()
}

func (r ReviewService) CreateReview(reviewInput *CreateReviewInput) error {
	review := &domain.Review{
		UserID:  reviewInput.UserID,
		BookID:  reviewInput.BookID,
		Rating:  reviewInput.Rating,
		Message: reviewInput.Message,
	}

	err := r.updateBookRating(review.BookID, review.Rating)
	if err != nil {
		return err
	}

	return r.repository.CreateReview(review)
}

func (r ReviewService) UpdateReview(reviewInput *UpdateReviewInput) error {
	review, err := r.repository.GetReviewByID(reviewInput.ID)
	if err != nil {
		return err
	}

	review.UserID = reviewInput.UserID
	review.BookID = reviewInput.BookID
	review.Rating = reviewInput.Rating
	review.Message = reviewInput.Message

	return r.repository.UpdateReview(review)
}

func (r ReviewService) DeleteReview(id uint) error {
	return r.repository.DeleteReview(id)
}

func (r ReviewService) GetReviewsByBookID(id uint) ([]*domain.Review, error) {
	return r.repository.GetReviewsByBookID(id)
}

func (r ReviewService) updateBookRating(bookID uint, rating float32) error {
	book, err := r.bookRepository.GetBookByID(bookID)
	if err != nil {
		return err
	}

	bookReviews, err := r.repository.GetReviewsByBookID(bookID)
	if err != nil {
		return err
	}

	for _, r := range bookReviews {
		rating += r.Rating
	}

	rating /= float32(len(bookReviews) + 1)

	book.Rating = rating

	return r.bookRepository.UpdateBook(book)
}
