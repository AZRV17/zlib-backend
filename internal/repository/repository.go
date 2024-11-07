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
	CreateAuthorBook(authorBook *domain.AuthorBook) error
	DeleteAuthorBook(id uint) error
	UpdateAuthorBook(authorBook *domain.AuthorBook) error
}

type BookRepo interface {
	GetBookByID(id uint) (*domain.Book, error)
	GetBooks() ([]*domain.Book, error)
	CreateBook(book *domain.Book) error
	UpdateBook(book *domain.Book) error
	DeleteBook(id uint) error
	GetBookByTitle(title string) (*domain.Book, error)
	GetBookByUniqueCode(code uint) (*domain.Book, error)
	GetGroupedBooksByTitle() ([]*domain.Book, error)
	GetBookUniqueCodes(id uint) ([]*domain.UniqueCode, error)
	CreateUniqueCode(uniqueCode *domain.UniqueCode) error
	DeleteUniqueCode(id uint) error
	UpdateUniqueCode(uniqueCode *domain.UniqueCode) error
	GetBookByIDWithTransactions(id uint, tx *gorm.DB) (*domain.Book, error)
	GetBookUniqueCodesWithTransactions(id uint, tx *gorm.DB) ([]*domain.UniqueCode, error)
	UpdateUniqueCodeWithTransactions(uniqueCode *domain.UniqueCode, tx *gorm.DB) error
	GetUniqueCodes() ([]*domain.UniqueCode, error)
	GetUniqueCodeByID(id uint) (*domain.UniqueCode, error)
}

type FavoriteRepo interface {
	GetFavoriteByID(id uint) (*domain.Favorite, error)
	GetFavorites() ([]*domain.Favorite, error)
	CreateFavorite(favorite *domain.Favorite) error
	DeleteFavorite(id uint) error
	GetFavoritesByUserID(id uint) ([]*domain.Favorite, error)
	DeleteFavoriteByUserIDAndBookID(userID uint, bookID uint) (*domain.Favorite, error)
}

type GenreRepo interface {
	GetGenreByID(id uint) (*domain.Genre, error)
	GetGenres() ([]*domain.Genre, error)
	CreateGenre(genre *domain.Genre) error
	UpdateGenre(genre *domain.Genre) error
	DeleteGenre(id uint) error
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
	CreateReservationWithTransactions(reservation *domain.Reservation, tx *gorm.DB) error
	GetUserReservations(id uint) ([]*domain.Reservation, error)
}

type ReviewRepo interface {
	GetReviewByID(id uint) (*domain.Review, error)
	GetReviews() ([]*domain.Review, error)
	CreateReview(review *domain.Review) error
	UpdateReview(review *domain.Review) error
	DeleteReview(id uint) error
	GetReviewsByBookID(id uint) ([]*domain.Review, error)
}

type UpdateUserDTOInput struct {
	ID             uint        `json:"id" gorm:"primaryKey,autoIncrement"`
	Login          string      `json:"login" gore:"unique"`
	FullName       string      `json:"full_name"`
	Password       string      `json:"password"`
	Role           domain.Role `json:"role" gorm:"type:role;default:'user'"`
	Email          string      `json:"email" gorm:"unique"`
	PhoneNumber    string      `json:"phoneNumber" gorm:"unique"`
	PassportNumber int         `json:"passportNumber" gorm:"unique"`
}

type UserRepo interface {
	GetUserByID(id uint) (*domain.User, error)
	GetUsers() ([]*domain.User, error)
	SignInByLogin(login, password string) (*domain.User, error)
	SignInByEmail(email, password string) (*domain.User, error)
	SignUp(user *domain.User) error
	DeleteUser(id uint) error
	UpdateUser(user *UpdateUserDTOInput) error
	GetUserByLogin(login string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUserRole(id uint, role domain.Role) error
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
		DB:               db,
		AuthorRepo:       NewAuthorRepository(db.Model(&domain.Author{})),
		BookRepo:         NewBookRepository(db),
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
