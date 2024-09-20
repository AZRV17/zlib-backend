package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type AuthorRepo interface {
	GetAuthorByID(id int) (*domain.Author, error)
	GetAuthors() ([]*domain.Author, error)
	CreateAuthor(author *domain.Author) error
	UpdateAuthor(author *domain.Author) error
	DeleteAuthor(id int) error
	GetAuthorBooks(id int) ([]*domain.Book, error)
}

type BookRepo interface {
	GetBookByID(id int) (*domain.Book, error)
	GetBooks() ([]*domain.Book, error)
	CreateBook(book *domain.Book) error
	UpdateBook(book *domain.Book) error
	DeleteBook(id int) error
	GetBookByTitle(title string) (*domain.Book, error)
	GetBookByUniqueCode(code int) (*domain.Book, error)
}

type FavoriteRepo interface {
	GetFavoriteByID(id int) (*domain.Favorite, error)
	GetFavorites() ([]*domain.Favorite, error)
	CreateFavorite(favorite *domain.Favorite) error
	DeleteFavorite(id int) error
}

type GenreRepo interface {
	GetGenreByID(id int) (*domain.Genre, error)
	GetGenres() ([]*domain.Genre, error)
	CreateGenre(genre *domain.Genre) error
	UpdateGenre(genre *domain.Genre) error
}

type LogRepo interface {
	GetLogByID(id int) (*domain.Log, error)
	GetLogs() ([]*domain.Log, error)
	CreateLog(log *domain.Log) error
	UpdateLog(log *domain.Log) error
	DeleteLog(id int) error
	GetLogsByUserID(id int) ([]*domain.Log, error)
}

type NotificationRepo interface {
	GetNotificationByID(id int) (*domain.Notification, error)
	GetNotifications() ([]*domain.Notification, error)
	CreateNotification(notification *domain.Notification) error
	UpdateNotification(notification *domain.Notification) error
	DeleteNotification(id int) error
	GetNotificationsByUserID(id int) ([]*domain.Notification, error)
}

type PublisherRepo interface {
	GetPublisherByID(id int) (*domain.Publisher, error)
	GetPublishers() ([]*domain.Publisher, error)
	CreatePublisher(publisher *domain.Publisher) error
	UpdatePublisher(publisher *domain.Publisher) error
	DeletePublisher(id int) error
}

type ReservationRepo interface {
	GetReservationByID(id int) (*domain.Reservation, error)
	GetReservations() ([]*domain.Reservation, error)
	CreateReservation(reservation *domain.Reservation) error
	UpdateReservation(reservation *domain.Reservation) error
	DeleteReservation(id int) error
}

type ReviewRepo interface {
	GetReviewByID(id int) (*domain.Review, error)
	GetReviews() ([]*domain.Review, error)
	CreateReview(review *domain.Review) error
	UpdateReview(review *domain.Review) error
	DeleteReview(id int) error
	GetReviewsByBookID(id int) ([]*domain.Review, error)
}

type UserRepo interface {
	GetUserByID(id int) (*domain.User, error)
	GetUsers() ([]*domain.User, error)
	SignInByLogin(login, password string) (*domain.User, error)
	SignInByEmail(email, password string) (*domain.User, error)
	DeleteUser(id int) error
	UpdateUser(user *domain.User) error
}

type Repository struct {
	DB *gorm.DB
	AuthorRepo
	BookRepo
	FavoriteRepo
	GenreRepo
	LogRepo
	NotificationRepo
	PublisherRepo
	ReservationRepo
	ReviewRepo
	UserRepo
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		DB:         db,
		AuthorRepo: NewAuthorRepository(db.Model(&domain.Author{})),
	}
}
