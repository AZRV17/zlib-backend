package repository

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type AuthorRepo interface {
	GetAuthorByID(id uint) (*domain.Author, error)
	GetAuthors() ([]*domain.Author, error)
	CreateAuthor(author *domain.Author) error
	UpdateAuthor(author *domain.Author) error
	DeleteAuthor(id uint) error
	GetAuthorBooks(id uint) ([]*domain.Book, error)
}

type AuthorBookRepo interface {
	GetAuthorBookByID(id uint) (*domain.AuthorBook, error)
	GetAuthorBooks() ([]*domain.AuthorBook, error)
	CreateAuthorBook(authorBook *domain.AuthorBook) error
	UpdateAuthorBook(authorBook *domain.AuthorBook) error
	DeleteAuthorBook(id uint) error
}

type BookRepo interface {
	GetBookByID(id uint) (*domain.Book, error)
	GetBooks() ([]*domain.Book, error)
	CreateBook(book *domain.Book) error
	UpdateBook(book *domain.Book) error
	DeleteBook(id uint) error
	GetBookByTitle(title string) (*domain.Book, error)
	GetBookByUniqueCode(code uint) (*domain.Book, error)
}

type FavoriteRepo interface {
	GetFavoriteByID(id uint) (*domain.Favorite, error)
	GetFavorites() ([]*domain.Favorite, error)
	CreateFavorite(favorite *domain.Favorite) error
	DeleteFavorite(id uint) error
}

type GenreRepo interface {
	GetGenreByID(id uint) (*domain.Genre, error)
	GetGenres() ([]*domain.Genre, error)
	CreateGenre(genre *domain.Genre) error
	UpdateGenre(genre *domain.Genre) error
}

type LogRepo interface {
	GetLogByID(id uint) (*domain.Log, error)
	GetLogs() ([]*domain.Log, error)
	CreateLog(log *domain.Log) error
	UpdateLog(log *domain.Log) error
	DeleteLog(id uint) error
	GetLogsByUserID(id uint) ([]*domain.Log, error)
}

type NotificationRepo interface {
	GetNotificationByID(id uint) (*domain.Notification, error)
	GetNotifications() ([]*domain.Notification, error)
	CreateNotification(notification *domain.Notification) error
	UpdateNotification(notification *domain.Notification) error
	DeleteNotification(id uint) error
	GetNotificationsByUserID(id uint) ([]*domain.Notification, error)
}

type PublisherRepo interface {
	GetPublisherByID(id uint) (*domain.Publisher, error)
	GetPublishers() ([]*domain.Publisher, error)
	CreatePublisher(publisher *domain.Publisher) error
	UpdatePublisher(publisher *domain.Publisher) error
	DeletePublisher(id uint) error
}

type ReservationRepo interface {
	GetReservationByID(id uint) (*domain.Reservation, error)
	GetReservations() ([]*domain.Reservation, error)
	CreateReservation(reservation *domain.Reservation) error
	UpdateReservation(reservation *domain.Reservation) error
	DeleteReservation(id uint) error
}

type ReviewRepo interface {
	GetReviewByID(id uint) (*domain.Review, error)
	GetReviews() ([]*domain.Review, error)
	CreateReview(review *domain.Review) error
	UpdateReview(review *domain.Review) error
	DeleteReview(id uint) error
	GetReviewsByBookID(id uint) ([]*domain.Review, error)
}

type UserRepo interface {
	GetUserByID(id uint) (*domain.User, error)
	GetUsers() ([]*domain.User, error)
	SignInByLogin(login, password string) (*domain.User, error)
	SignInByEmail(email, password string) (*domain.User, error)
	SignUp(user *domain.User) error
	DeleteUser(id uint) error
	UpdateUser(user *domain.User) error
}

type Repository struct {
	DB *gorm.DB
	AuthorRepo
	AuthorBookRepo
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
		DB:               db,
		AuthorRepo:       NewAuthorRepository(db.Model(&domain.Author{})),
		AuthorBookRepo:   NewAuthorBookRepository(db.Model(&domain.AuthorBook{})),
		BookRepo:         NewBookRepository(db.Model(&domain.Book{})),
		FavoriteRepo:     NewFavoriteRepository(db.Model(&domain.Favorite{})),
		GenreRepo:        NewGenreRepository(db.Model(&domain.Genre{})),
		LogRepo:          NewLogRepository(db.Model(&domain.Log{})),
		NotificationRepo: NewNotificationRepository(db.Model(&domain.Notification{})),
		PublisherRepo:    NewPublisherRepository(db.Model(&domain.Publisher{})),
		ReservationRepo:  NewReservationRepository(db.Model(&domain.Reservation{})),
		ReviewRepo:       NewReviewRepository(db.Model(&domain.Review{})),
		UserRepo:         NewUserRepository(db.Model(&domain.User{})),
	}
}
