package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type CreateAuthorInput struct {
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Biography string    `json:"biography"`
	Birthdate time.Time `json:"birthdate"`
}

type UpdateAuthorInput struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Biography string    `json:"biography"`
	Birthdate time.Time `json:"birthdate"`
}

type AuthorServ interface {
	GetAuthorByID(id uint) (*domain.Author, error)
	GetAuthors() ([]*domain.Author, error)
	CreateAuthor(authorInput *CreateAuthorInput) error
	UpdateAuthor(authorInput *UpdateAuthorInput) error
	DeleteAuthor(id uint) error
	GetAuthorBooks(id uint) ([]*domain.Book, error)
	CreateAuthorBook(authorBookInput *domain.AuthorBook) error
	DeleteAuthorBook(id uint) error
	ExportAuthorsToCSV() ([]byte, error)
}

type CreateBookInput struct {
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	AuthorID          uint      `json:"author_id"`
	Author            []uint    `json:"author,omitempty"`
	GenreID           uint      `json:"genre_id"`
	Genre             uint      `json:"genre"`
	PublisherID       uint      `json:"publisher_id"`
	Publisher         uint      `json:"publisher"`
	ISBN              int       `json:"isbn"`
	YearOfPublication time.Time `json:"year_of_publication"`
	Picture           string    `json:"picture"`
	Rating            float32   `json:"rating"`
}

type UpdateBookInput struct {
	ID                uint      `json:"id"`
	Title             string    `json:"title"`
	AuthorID          uint      `json:"author_id"`
	Author            []uint    `json:"author,omitempty"`
	GenreID           uint      `json:"genre_id"`
	Genre             uint      `json:"genre"`
	PublisherID       uint      `json:"publisher_id"`
	Publisher         uint      `json:"publisher"`
	ISBN              int       `json:"isbn"`
	YearOfPublication time.Time `json:"year_of_publication"`
	Picture           string    `json:"picture"`
	Rating            float32   `json:"rating"`
	IsAvailable       bool      `json:"is_available"`
}

type BookServ interface {
	GetBookByID(id uint) (*domain.Book, error)
	GetBooks() ([]*domain.Book, error)
	CreateBook(bookInput *CreateBookInput) error
	UpdateBook(bookInput *UpdateBookInput) error
	DeleteBook(id uint) error
	GetBookByTitle(title string) (*domain.Book, error)
	GetBookByUniqueCode(code uint) (*domain.Book, error)
	GetGroupedBooksByTitle() ([]*domain.Book, error)
	CreateUniqueCode(uniqueCode *domain.UniqueCode) error
	DeleteUniqueCode(id uint) error
	UpdateUniqueCode(uniqueCode *domain.UniqueCode) error
	GetBookUniqueCodes(id uint) ([]*domain.UniqueCode, error)
	ReserveBook(bookID, userID uint) (*domain.UniqueCode, error)
	GetUniqueCodes() ([]*domain.UniqueCode, error)
	GetUniqueCodeByID(id uint) (*domain.UniqueCode, error)
	GetBooksWithPagination(limit int, offset int) ([]*domain.Book, error)
	FindBookByTitle(limit int, offset int, title string) ([]*domain.Book, error)
}

type CreateFavoriteInput struct {
	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}

type FavoriteServ interface {
	GetFavoriteByID(id uint) (*domain.Favorite, error)
	GetFavorites() ([]*domain.Favorite, error)
	CreateFavorite(favoriteInput *CreateFavoriteInput) error
	DeleteFavorite(id uint) error
	GetFavoriteByUserID(userID uint) ([]*domain.Favorite, error)
	DeleteFavoriteByUserIDAndBookID(userID uint, bookID uint) (*domain.Favorite, error)
}

type CreateGenreInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateGenreInput struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GenreServ interface {
	GetGenreByID(id uint) (*domain.Genre, error)
	GetGenres() ([]*domain.Genre, error)
	CreateGenre(genreInput *CreateGenreInput) error
	UpdateGenre(genreInput *UpdateGenreInput) error
	DeleteGenre(id uint) error
}

type CreateLogInput struct {
	UserID  uint      `json:"user_id"`
	Action  string    `json:"action"`
	Date    time.Time `json:"date"`
	Details string    `json:"details"`
}

type UpdateLogInput struct {
	ID      uint      `json:"id"`
	UserID  uint      `json:"user_id"`
	Action  string    `json:"action"`
	Date    time.Time `json:"date"`
	Details string    `json:"details"`
}

type LogServ interface {
	GetLogByID(id uint) (*domain.Log, error)
	GetLogs() ([]*domain.Log, error)
	CreateLog(logInput *CreateLogInput) error
	UpdateLog(logInput *UpdateLogInput) error
	DeleteLog(id uint) error
	GetLogsByUserID(id uint) ([]*domain.Log, error)
	CreateLogWithCookie(cookie *http.Cookie, action string) error
}

type CreateNotificationInput struct {
	UserID  uint   `json:"user_id"`
	Message string `json:"message"`
}

type NotificationServ interface {
	GetNotificationByID(id uint) (*domain.Notification, error)
	GetNotifications() ([]*domain.Notification, error)
	CreateNotification(notificationInput *CreateNotificationInput) error
	DeleteNotification(id uint) error
}

type CreatePublisherInput struct {
	Name string `json:"name"`
}

type UpdatePublisherInput struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PublisherServ interface {
	GetPublisherByID(id uint) (*domain.Publisher, error)
	GetPublishers() ([]*domain.Publisher, error)
	CreatePublisher(publisherInput *CreatePublisherInput) error
	UpdatePublisher(publisherInput *UpdatePublisherInput) error
	DeletePublisher(id uint) error
}

type CreateReservationInput struct {
	UserID       uint      `json:"user_id"`
	BookID       uint      `json:"book_id"`
	Status       string    `json:"status,omitempty"`
	DateOfReturn time.Time `json:"date_of_return"`
	DateOfIssue  time.Time `json:"date_of_issue"`
}

type UpdateReservationInput struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	BookID       uint      `json:"book_id"`
	Status       string    `json:"status,omitempty"`
	DateOfReturn time.Time `json:"date_of_return"`
	DateOfIssue  time.Time `json:"date_of_issue"`
}

type ReservationServ interface {
	GetReservationByID(id uint) (*domain.Reservation, error)
	GetReservations() ([]*domain.Reservation, error)
	CreateReservation(reservationInput *CreateReservationInput) error
	UpdateReservation(reservationInput *UpdateReservationInput) error
	DeleteReservation(id uint) error
	GetReservationsByUserID(id uint) ([]*domain.Reservation, error)
}

type CreateReviewInput struct {
	UserID  uint    `json:"user_id"`
	BookID  uint    `json:"book_id"`
	Rating  float32 `json:"rating"`
	Message string  `json:"message"`
}

type UpdateReviewInput struct {
	ID      uint    `json:"id"`
	UserID  uint    `json:"user_id"`
	BookID  uint    `json:"book_id"`
	Review  string  `json:"review"`
	Rating  float32 `json:"rating"`
	Message string  `json:"message"`
}

type ReviewServ interface {
	GetReviewByID(id uint) (*domain.Review, error)
	GetReviews() ([]*domain.Review, error)
	CreateReview(reviewInput *CreateReviewInput) error
	UpdateReview(reviewInput *UpdateReviewInput) error
	DeleteReview(id uint) error
	GetReviewsByBookID(id uint) ([]*domain.Review, error)
}

type SignUpUserInput struct {
	Login          string      `json:"login"`
	FullName       string      `json:"full_name"`
	Password       string      `json:"password"`
	Email          string      `json:"email" binding:"required,email"`
	Role           domain.Role `json:"role"`
	PhoneNumber    string      `json:"phoneNumber"`
	PassportNumber int         `json:"passportNumber"`
}

type UpdateUserInput struct {
	ID             uint        `json:"id" gorm:"primaryKey,autoIncrement"`
	Login          string      `json:"login" gore:"unique"`
	FullName       string      `json:"full_name"`
	Password       string      `json:"password"`
	Role           domain.Role `json:"role" gorm:"type:role;default:'user'"`
	Email          string      `json:"email" gorm:"unique"`
	PhoneNumber    string      `json:"phoneNumber" gorm:"unique"`
	PassportNumber int         `json:"passportNumber" gorm:"unique"`
}

type UserServ interface {
	SignInByLogin(login, password string) (*domain.User, error)
	SignInByEmail(email, password string) (*domain.User, error)
	SignUp(userInput *SignUpUserInput) error
	GetUserByID(id uint) (*domain.User, error)
	GetAllUsers() ([]*domain.User, error)
	UpdateUser(userInput *UpdateUserInput) error
	GetUserByLogin(login string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdateUserRole(id uint, role domain.Role) error
	DeleteUser(id uint) error
	hashPassword(password string) (string, error)
	comparePasswords(hashedPassword, password string) bool
}

type Service struct {
	repo *repository.Repository
	AuthorServ
	BookServ
	FavoriteServ
	GenreServ
	LogServ
	NotificationServ
	PublisherServ
	ReservationServ
	ReviewServ
	UserServ
}

func NewService(repo *repository.Repository, db *gorm.DB) *Service {
	return &Service{
		repo:             repo,
		AuthorServ:       NewAuthorService(repo.AuthorRepo),
		BookServ:         NewBookService(repo.BookRepo, repo.ReservationRepo, db),
		FavoriteServ:     NewFavoriteService(repo.FavoriteRepo),
		GenreServ:        NewGenreService(repo.GenreRepo),
		LogServ:          NewLogService(repo.LogRepo),
		NotificationServ: NewNotificationService(repo.NotificationRepo),
		PublisherServ:    NewPublisherService(repo.PublisherRepo),
		ReservationServ:  NewReservationService(repo.ReservationRepo),
		ReviewServ:       NewReviewService(repo.ReviewRepo),
		UserServ:         NewUserService(repo.UserRepo),
	}
}
