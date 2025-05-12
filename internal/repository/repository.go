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
	ExportAuthorsToCSV() ([]byte, error)
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
	GetBooksWithPagination(limit int, offset int) ([]*domain.Book, error)
	FindBookByTitle(limit int, offset int, title string) ([]*domain.Book, error)
	FindBooks(limit int, offset int, query string) ([]*domain.Book, error) // Новый метод
	ExportBooksToCSV() ([]byte, error)

	// Методы для работы с аудиофайлами книги
	GetAudiobookFilesByBookID(bookID uint) ([]*domain.AudiobookFile, error)
	GetAudiobookFileByID(id uint) (*domain.AudiobookFile, error)
	CreateAudiobookFile(file *domain.AudiobookFile) error
	UpdateAudiobookFile(file *domain.AudiobookFile) error
	DeleteAudiobookFile(id uint) error
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
	ExportGenresToCSV() ([]byte, error)
}

type LogRepo interface {
	GetLogByID(id uint) (*domain.Log, error)
	GetLogs() ([]*domain.Log, error)
	CreateLog(log *domain.Log) error
	UpdateLog(log *domain.Log) error
	DeleteLog(id uint) error
	GetLogsByUserID(id uint) ([]*domain.Log, error)
}

type PublisherRepo interface {
	GetPublisherByID(id uint) (*domain.Publisher, error)
	GetPublishers() ([]*domain.Publisher, error)
	CreatePublisher(publisher *domain.Publisher) error
	UpdatePublisher(publisher *domain.Publisher) error
	DeletePublisher(id uint) error
	ExportPublishersToCSV() ([]byte, error)
}

type ReservationRepo interface {
	GetReservationByID(id uint) (*domain.Reservation, error)
	GetReservations() ([]*domain.Reservation, error)
	CreateReservation(reservation *domain.Reservation) error
	UpdateReservation(reservation *domain.Reservation) error
	DeleteReservation(id uint) error
	CreateReservationWithTransactions(reservation *domain.Reservation, tx *gorm.DB) error
	GetUserReservations(id uint) ([]*domain.Reservation, error)
	ExportReservationsToCSV() ([]byte, error)
}

type ReviewRepo interface {
	GetReviewByID(id uint) (*domain.Review, error)
	GetReviews() ([]*domain.Review, error)
	CreateReview(review *domain.Review) error
	UpdateReview(review *domain.Review) error
	DeleteReview(id uint) error
	GetReviewsByBookID(id uint) ([]*domain.Review, error)
	CheckUserReviewExists(userID uint, bookID uint) (bool, error)
}

type UpdateUserDTOInput struct {
	ID             uint        `json:"id" gorm:"primaryKey,autoIncrement"`
	Login          string      `json:"login" gore:"unique"`
	FullName       string      `json:"full_name"`
	Password       string      `json:"password"`
	Role           domain.Role `json:"role" gorm:"type:role;default:'user'"`
	Email          string      `json:"email" gorm:"unique"`
	PhoneNumber    string      `json:"phone_number" gorm:"unique"`
	PassportNumber int         `json:"passport_number" gorm:"unique"`
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
	SetResetPasswordToken(userID uint, token string, expiry string) error
	GetUserByResetToken(token string) (*domain.User, error)
	UpdatePassword(userID uint, password string) error
}

type ChatRepo interface {
	SaveMessage(message *domain.Message) error
	CreateChat(chat *domain.Chat) error
	GetChatByID(chatID uint) (*domain.Chat, error)
	GetMessagesByChatID(chatID uint) ([]domain.Message, error)
	GetActiveChatsForLibrarian() ([]domain.Chat, error)
	GetChatsByUserID(userID uint) ([]domain.Chat, error)
	AssignLibrarianToChat(chatID, librarianID uint) error
	CloseChat(chatID uint) error
	MarkMessagesAsRead(chatID, userID uint) error
	GetLibrarianChats(librarianID uint) ([]domain.Chat, error)
	GetUnassignedChats() ([]domain.Chat, error)
}

type Repository struct {
	DB *gorm.DB
	AuthorRepo
	BookRepo
	FavoriteRepo
	GenreRepo
	LogRepo
	PublisherRepo
	ReservationRepo
	ReviewRepo
	UserRepo
	ChatRepo
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		DB:              db,
		AuthorRepo:      NewAuthorRepository(db.Model(&domain.Author{})),
		BookRepo:        NewBookRepository(db),
		FavoriteRepo:    NewFavoriteRepository(db.Model(&domain.Favorite{})),
		GenreRepo:       NewGenreRepository(db.Model(&domain.Genre{})),
		LogRepo:         NewLogRepository(db.Model(&domain.Log{})),
		PublisherRepo:   NewPublisherRepository(db.Model(&domain.Publisher{})),
		ReservationRepo: NewReservationRepository(db),
		ReviewRepo:      NewReviewRepository(db.Model(&domain.Review{})),
		UserRepo:        NewUserRepository(db.Model(&domain.User{})),
		ChatRepo:        NewChatRepository(db),
	}
}
