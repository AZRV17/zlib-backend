package service

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookService struct {
	bookRepo        repository.BookRepo
	reservationRepo repository.ReservationRepo
	db              *gorm.DB
}

func (b BookService) CountAllBooks() (int, error) {
	books, err := b.bookRepo.GetBooks()
	if err != nil {
		return 0, err
	}
	return len(books), nil
}

func (b BookService) CountBooksMatchingSearch(query string) (int, error) {
	return b.bookRepo.CountBooksMatchingSearch(query)
}

func (b BookService) CountBooksMatchingTitle(title string) (int, error) {
	// Получаем все книги с заданным заголовком
	books, err := b.bookRepo.FindBookByTitle(1000000, 0, title) // Используем большой лимит для получения всех книг
	if err != nil {
		return 0, err
	}
	return len(books), nil
}

func NewBookService(
	bookRepo repository.BookRepo,
	reservationRepo repository.ReservationRepo,
	db *gorm.DB,
) *BookService {
	return &BookService{bookRepo: bookRepo, reservationRepo: reservationRepo, db: db}
}

func (b BookService) GetBookByID(id uint) (*domain.Book, error) {
	book, err := b.bookRepo.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	if book.Picture == "" || strings.HasPrefix(book.Picture, "http") {
		return book, nil
	}

	book.Picture = "http://localhost:8080/" + book.Picture
	return book, err
}

func (b BookService) GetBooks() ([]*domain.Book, error) {
	books, err := b.bookRepo.GetBooks()
	if err != nil {
		return nil, err
	}

	for _, book := range books {
		if book.Picture == "" || strings.HasPrefix(book.Picture, "http") {
			continue
		}

		book.Picture = "http://localhost:8080/" + book.Picture
	}

	return books, nil
}

func (b BookService) CreateBook(bookInput *CreateBookInput) error {
	book := &domain.Book{
		Title:             bookInput.Title,
		Description:       bookInput.Description,
		AuthorID:          bookInput.AuthorID,
		GenreID:           bookInput.GenreID,
		PublisherID:       bookInput.PublisherID,
		ISBN:              bookInput.ISBN,
		YearOfPublication: bookInput.YearOfPublication,
		Picture:           bookInput.Picture,
		Rating:            bookInput.Rating,
		EpubFile:          bookInput.EpubFile,
	}

	return b.bookRepo.CreateBook(book)
}

func (b BookService) UpdateBook(bookInput *UpdateBookInput) error {
	book, err := b.bookRepo.GetBookByID(bookInput.ID)
	if err != nil {
		return err
	}

	log.Println("bookInput", bookInput)

	book.Title = bookInput.Title
	book.AuthorID = bookInput.AuthorID
	book.GenreID = bookInput.GenreID
	book.PublisherID = bookInput.PublisherID
	book.ISBN = bookInput.ISBN
	book.YearOfPublication = bookInput.YearOfPublication
	book.Picture = bookInput.Picture
	book.Rating = bookInput.Rating
	book.EpubFile = bookInput.EpubFile
	book.Description = bookInput.Description

	return b.bookRepo.UpdateBook(book)
}

func (b BookService) DeleteBook(id uint) error {
	return b.bookRepo.DeleteBook(id)
}

func (b BookService) GetBookByTitle(title string) (*domain.Book, error) {
	return b.bookRepo.GetBookByTitle(title)
}

func (b BookService) GetBookByUniqueCode(code uint) (*domain.Book, error) {
	return b.bookRepo.GetBookByUniqueCode(code)
}

func (b BookService) GetGroupedBooksByTitle() ([]*domain.Book, error) {
	return b.bookRepo.GetGroupedBooksByTitle()
}

func (b BookService) CreateUniqueCode(uniqueCode *domain.UniqueCode) error {
	return b.bookRepo.CreateUniqueCode(uniqueCode)
}

func (b BookService) DeleteUniqueCode(id uint) error {
	return b.bookRepo.DeleteUniqueCode(id)
}

func (b BookService) UpdateUniqueCode(uniqueCode *domain.UniqueCode) error {
	return b.bookRepo.UpdateUniqueCode(uniqueCode)
}

func (b BookService) GetUniqueCodes() ([]*domain.UniqueCode, error) {
	return b.bookRepo.GetUniqueCodes()
}

func (b BookService) GetBookUniqueCodes(id uint) ([]*domain.UniqueCode, error) {
	return b.bookRepo.GetBookUniqueCodes(id)
}

func (b BookService) ReserveBook(bookID, userID uint) (*domain.UniqueCode, error) {
	tx := b.db.Begin()

	book, err := b.bookRepo.GetBookByIDWithTransactions(bookID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	codes, err := b.bookRepo.GetBookUniqueCodesWithTransactions(book.ID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var code *domain.UniqueCode

	for _, c := range codes {
		if c.IsAvailable {
			code = c
		}
	}

	if code == nil {
		tx.Rollback()
		return nil, errors.New("no available books")
	}

	err = b.reservationRepo.CreateReservationWithTransactions(
		&domain.Reservation{
			BookID:       book.ID,
			UserID:       userID,
			UniqueCodeID: code.ID,
			Status:       "reserved",
			DateOfIssue:  time.Now(),
			DateOfReturn: time.Now().Add(7 * 24 * time.Hour),
		}, tx,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	code.IsAvailable = false

	err = b.bookRepo.UpdateUniqueCodeWithTransactions(code, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return code, nil
}

func (b *BookService) GetUniqueCodeByID(id uint) (*domain.UniqueCode, error) {
	return b.bookRepo.GetUniqueCodeByID(id)
}

func (b *BookService) GetBooksWithPagination(offset, limit int) ([]*domain.Book, error) {
	books, err := b.bookRepo.GetBooksWithPagination(offset, limit)
	if err != nil {
		return nil, err
	}

	for _, book := range books {
		if book.Picture == "" || strings.HasPrefix(book.Picture, "http") {
			continue
		}

		book.Picture = "http://localhost:8080/" + book.Picture
	}

	return books, nil
}

func (b *BookService) FindBookByTitle(limit int, offset int, title string) ([]*domain.Book, error) {
	return b.bookRepo.FindBookByTitle(limit, offset, title)
}

func (b *BookService) FindBooks(limit int, offset int, query string) ([]*domain.Book, error) {
	books, err := b.bookRepo.FindBooks(limit, offset, query)
	if err != nil {
		return nil, err
	}

	// Добавляем префикс к путям изображений, если они есть и не начинаются с http
	for _, book := range books {
		if book.Picture == "" || strings.HasPrefix(book.Picture, "http") {
			continue
		}

		book.Picture = "http://localhost:8080/" + book.Picture
	}

	return books, nil
}

func (b BookService) GetAudiobookFilesByBookID(bookID uint) ([]*domain.AudiobookFile, error) {
	return b.bookRepo.GetAudiobookFilesByBookID(bookID)
}

func (b BookService) GetAudiobookFileByID(id uint) (*domain.AudiobookFile, error) {
	return b.bookRepo.GetAudiobookFileByID(id)
}

func (b BookService) CreateAudiobookFile(file *domain.AudiobookFile, fileData []byte) error {
	fileName := fmt.Sprintf("%s.mp3", uuid.New().String())
	savePath := filepath.Join("uploads", "audio", fileName)

	if err := os.MkdirAll("uploads/audio", os.ModePerm); err != nil {
		return fmt.Errorf("ошибка при создании папки: %w", err)
	}

	if err := os.WriteFile(savePath, fileData, 0644); err != nil {
		return fmt.Errorf("ошибка при сохранении аудиофайла: %w", err)
	}

	file.FilePath = savePath

	return b.bookRepo.CreateAudiobookFile(file)
}

func (b BookService) UpdateAudiobookFile(file *domain.AudiobookFile, fileData []byte) error {
	oldFile, err := b.bookRepo.GetAudiobookFileByID(file.ID)
	if err != nil {
		return fmt.Errorf("ошибка при поиске аудиофайла: %w", err)
	}

	if len(fileData) > 0 {
		if oldFile.FilePath != "" {
			_ = os.Remove(oldFile.FilePath)
		}

		fileName := fmt.Sprintf("%s.mp3", uuid.New().String())
		savePath := filepath.Join("uploads", "audio", fileName)

		if err := os.WriteFile(savePath, fileData, 0644); err != nil {
			return fmt.Errorf("ошибка при сохранении нового аудиофайла: %w", err)
		}

		file.FilePath = savePath
	} else {
		file.FilePath = oldFile.FilePath
	}

	return b.bookRepo.UpdateAudiobookFile(file)
}

func (b BookService) DeleteAudiobookFile(id uint) error {
	file, err := b.bookRepo.GetAudiobookFileByID(id)
	if err != nil {
		return fmt.Errorf("ошибка при поиске аудиофайла: %w", err)
	}

	if file.FilePath != "" {
		if err := os.Remove(file.FilePath); err != nil {
			return fmt.Errorf("ошибка при удалении аудиофайла: %w", err)
		}
	}

	return b.bookRepo.DeleteAudiobookFile(id)
}

func (b BookService) UpdateAudiobookFileOrder(fileID uint, order int) error {
	file, err := b.bookRepo.GetAudiobookFileByID(fileID)
	if err != nil {
		return fmt.Errorf("ошибка при поиске аудиофайла: %w", err)
	}

	file.Order = order

	return b.bookRepo.UpdateAudiobookFile(file)
}

func (b BookService) ExportBooksToCSV() ([]byte, error) {
	return b.bookRepo.ExportBooksToCSV()
}

func (b BookService) ImportBooksFromCSV(data []byte) (int, error) {
	reader := csv.NewReader(bytes.NewReader(data))
	// Пропускаем заголовки
	if _, err := reader.Read(); err != nil {
		return 0, fmt.Errorf("ошибка при чтении заголовков CSV: %w", err)
	}

	importedCount := 0
	tx := b.db.Begin()

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			tx.Rollback()
			return importedCount, fmt.Errorf("ошибка при чтении записи CSV: %w", err)
		}

		// Проверяем, что у нас достаточно полей
		if len(record) < 6 {
			continue // Пропускаем некорректные записи
		}

		// Обработка полей CSV
		// [Title, AuthorName AuthorLastname, GenreName, PublisherName, YearOfPublication, ISBN]
		title := strings.TrimSpace(record[0])
		authorFullName := strings.TrimSpace(record[1])
		genreName := strings.TrimSpace(record[2])
		publisherName := strings.TrimSpace(record[3])
		yearStr := strings.TrimSpace(record[4])
		isbnStr := strings.TrimSpace(record[5])
		description := ""
		if len(record) > 6 {
			description = strings.TrimSpace(record[6])
		}

		// Разбиваем имя автора на имя и фамилию
		authorParts := strings.Split(authorFullName, " ")
		authorName := ""
		authorLastname := ""
		if len(authorParts) > 0 {
			authorName = authorParts[0]
			if len(authorParts) > 1 {
				authorLastname = strings.Join(authorParts[1:], " ")
			}
		}

		// Ищем или создаем автора
		var author domain.Author
		result := tx.Where("name = ? AND lastname = ?", authorName, authorLastname).First(&author)
		if result.Error != nil {
			// Автор не найден, создаем нового
			author = domain.Author{
				Name:     authorName,
				Lastname: authorLastname,
			}
			if err := tx.Create(&author).Error; err != nil {
				tx.Rollback()
				return importedCount, fmt.Errorf("ошибка при создании автора: %w", err)
			}
		}

		// Ищем или создаем жанр
		var genre domain.Genre
		result = tx.Where("name = ?", genreName).First(&genre)
		if result.Error != nil {
			// Жанр не найден, создаем новый
			genre = domain.Genre{
				Name: genreName,
			}
			if err := tx.Create(&genre).Error; err != nil {
				tx.Rollback()
				return importedCount, fmt.Errorf("ошибка при создании жанра: %w", err)
			}
		}

		// Ищем или создаем издателя
		var publisher domain.Publisher
		result = tx.Where("name = ?", publisherName).First(&publisher)
		if result.Error != nil {
			// Издатель не найден, создаем нового
			publisher = domain.Publisher{
				Name: publisherName,
			}
			if err := tx.Create(&publisher).Error; err != nil {
				tx.Rollback()
				return importedCount, fmt.Errorf("ошибка при создании издателя: %w", err)
			}
		}

		// Парсим ISBN
		isbn := 0
		if isbnStr != "" {
			var err error
			isbn, err = strconv.Atoi(isbnStr)
			if err != nil {
				// Если не удалось распарсить ISBN, просто пропускаем его
				isbn = 0
			}
		}

		// Парсим год публикации
		yearOfPublication, err := time.Parse("2006-01-02", yearStr)
		if err != nil {
			// Если не удалось распарсить дату, используем текущую
			yearOfPublication = time.Now()
		}

		// Создаем книгу
		book := domain.Book{
			Title:             title,
			Description:       description,
			AuthorID:          author.ID,
			GenreID:           genre.ID,
			PublisherID:       publisher.ID,
			ISBN:              isbn,
			YearOfPublication: yearOfPublication,
			Rating:            0, // Начальный рейтинг
		}

		// Проверяем, существует ли уже книга с таким названием и автором
		var existingBook domain.Book
		result = tx.Where("title = ? AND author_id = ?", title, author.ID).First(&existingBook)
		if result.Error == nil {
			// Книга уже существует, пропускаем
			continue
		}

		// Создаем книгу
		if err := tx.Create(&book).Error; err != nil {
			tx.Rollback()
			return importedCount, fmt.Errorf("ошибка при создании книги: %w", err)
		}

		importedCount++
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return importedCount, nil
}

func (b *BookService) FindBooksWithFilters(
	limit int,
	offset int,
	query string,
	authorID uint,
	genreID uint,
	yearStart, yearEnd time.Time,
	sortBy string,
	sortOrder string,
) ([]*domain.Book, error) {
	books, err := b.bookRepo.FindBooksWithFilters(limit, offset, query, authorID, genreID, yearStart, yearEnd, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	// Добавляем префикс к путям изображений, если они есть и не начинаются с http
	for _, book := range books {
		if book.Picture == "" || strings.HasPrefix(book.Picture, "http") {
			continue
		}

		book.Picture = "http://localhost:8080/" + book.Picture
	}

	return books, nil
}

func (b *BookService) CountBooksWithFilters(
	query string,
	authorID uint,
	genreID uint,
	yearStart, yearEnd time.Time,
) (int, error) {
	return b.bookRepo.CountBooksWithFilters(query, authorID, genreID, yearStart, yearEnd)
}

func (b BookService) GetTopBooksByReservations(limit int, periodMonths int) ([]*domain.Book, error) {
	books, err := b.bookRepo.GetTopBooksByReservations(limit, periodMonths)
	if err != nil {
		return nil, err
	}

	// Добавляем префикс к путям изображений
	for _, book := range books {
		if book.Picture == "" || strings.HasPrefix(book.Picture, "http") {
			continue
		}

		book.Picture = "http://localhost:8080/" + book.Picture
	}

	return books, nil
}
