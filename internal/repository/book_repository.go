package repository

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/gorm"
)

type BookRepository struct {
	DB *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (b BookRepository) GetBookByID(id uint) (*domain.Book, error) {
	var book domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).Preload("Author").Preload("Genre").Preload("Publisher").First(
		&book,
		id,
	).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetBooks() ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).
		Preload("Author").
		Preload("Genre").
		Preload("Publisher").
		Preload("UniqueCodes").
		Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookRepository) CreateBook(book *domain.Book) error {
	if err := b.DB.Model(&domain.Book{}).Create(book).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) UpdateBook(book *domain.Book) error {
	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).Where("id = ?", book.ID).Save(book).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (b BookRepository) DeleteBook(id uint) error {
	if err := b.DB.Model(&domain.Book{}).Delete(&domain.Book{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) GetBookByTitle(title string) (*domain.Book, error) {
	var book domain.Book

	if err := b.DB.Model(&domain.Book{}).Where("title = ?", title).First(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetBookByUniqueCode(code uint) (*domain.Book, error) {
	var book domain.Book

	if err := b.DB.Model(&domain.Book{}).Where("unique_code = ?", code).First(&book).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetGroupedBooksByTitle() ([]*domain.Book, error) {
	var books []*domain.Book

	if err := b.DB.Select(
		`
				array_agg(id) as id,
				title,
				author_id,
				genre_id,
				publisher_id,
				isbn,
				year_of_publication,
				picture,
				avg(rating) as rating,
				array_agg(unique_code) as unique_code,
				array_egg(created_at) as created_at,
				array_agg(updated_at) as updated_at,
				`,
	).Group("title").Error; err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookRepository) GetBookUniqueCodes(id uint) ([]*domain.UniqueCode, error) {
	var codes []*domain.UniqueCode

	tx := b.DB.Model(&domain.UniqueCode{}).Begin()

	if err := tx.Preload("Book").Preload("Book.Author").Where("book_id = ?", id).Find(&codes).Error; err != nil {
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func (b BookRepository) CreateUniqueCode(uniqueCode *domain.UniqueCode) error {
	tx := b.DB.Begin()

	if err := tx.Create(uniqueCode).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (b BookRepository) DeleteUniqueCode(id uint) error {
	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Delete(&domain.UniqueCode{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (b BookRepository) UpdateUniqueCode(uniqueCode *domain.UniqueCode) error {
	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Where("id = ?", uniqueCode.ID).Save(uniqueCode).Error; err != nil {
		tx.Rollback()
		return err
	}

	err := tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (b BookRepository) GetBookByIDWithTransactions(id uint, tx *gorm.DB) (*domain.Book, error) {
	var book domain.Book

	if err := tx.Model(&domain.Book{}).Preload("Author").Preload("Genre").Preload("Publisher").First(
		&book,
		id,
	).Error; err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookRepository) GetBookUniqueCodesWithTransactions(id uint, tx *gorm.DB) ([]*domain.UniqueCode, error) {
	var codes []*domain.UniqueCode

	if err := tx.Model(&domain.UniqueCode{}).Preload("Book").Preload("Book.Author").Where(
		"book_id = ?",
		id,
	).Find(&codes).Error; err != nil {
		return nil, err
	}

	return codes, nil
}

func (b BookRepository) UpdateUniqueCodeWithTransactions(uniqueCode *domain.UniqueCode, tx *gorm.DB) error {
	if err := tx.Model(&domain.UniqueCode{}).Where("id = ?", uniqueCode.ID).Save(uniqueCode).Error; err != nil {
		return err
	}

	return nil
}

func (b BookRepository) GetUniqueCodes() ([]*domain.UniqueCode, error) {
	var codes []*domain.UniqueCode

	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Preload("Book").Find(&codes).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return codes, nil
}

func (b BookRepository) GetUniqueCodeByID(id uint) (*domain.UniqueCode, error) {
	var code domain.UniqueCode

	tx := b.DB.Begin()

	if err := tx.Model(&domain.UniqueCode{}).Preload("Book").First(&code, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &code, nil
}

func (b BookRepository) GetBooksWithPagination(limit int, offset int) ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).Limit(limit).Offset(offset).Preload("Genre").Preload("Author").Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookRepository) FindBookByTitle(limit int, offset int, title string) ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).Where(
		"lower(title) like lower(?)",
		"%"+title+"%",
	).Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return books, nil
}

func (b BookRepository) ExportBooksToCSV() ([]byte, error) {
	books, err := b.GetBooks()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)

	headers := []string{"ID", "Название", "Автор", "Жанр", "Издательство", "Дата выхода"}

	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	for _, book := range books {
		row := []string{
			strconv.FormatUint(uint64(book.ID), 10),
			book.Title,
			book.Author.Name + " " + book.Author.Lastname,
			book.Genre.Name,
			book.Publisher.Name,
			book.YearOfPublication.Format("2006-01-02"),
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

func (b BookRepository) SortBooksByTitle() ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).Order("title").
		Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b BookRepository) GetAudiobookFilesByBookID(bookID uint) ([]*domain.AudiobookFile, error) {
	tx := b.DB.Begin()

	var files []*domain.AudiobookFile

	if err := tx.Model(&domain.AudiobookFile{}).Where("book_id = ?", bookID).Find(&files).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return files, nil
}

func (b BookRepository) GetAudiobookFileByID(id uint) (*domain.AudiobookFile, error) {
	tx := b.DB.Begin()

	var file domain.AudiobookFile

	if err := tx.Model(&domain.AudiobookFile{}).First(&file, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return &file, nil
}

func (b BookRepository) CreateAudiobookFile(file *domain.AudiobookFile) error {
	tx := b.DB.Begin()

	if err := tx.Create(file).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (b BookRepository) UpdateAudiobookFile(file *domain.AudiobookFile) error {
	tx := b.DB.Begin()

	if err := tx.Model(&domain.AudiobookFile{}).Where("id = ?", file.ID).Save(file).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (b BookRepository) DeleteAudiobookFile(id uint) error {
	tx := b.DB.Begin()

	if err := tx.Model(&domain.AudiobookFile{}).Delete(&domain.AudiobookFile{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (b BookRepository) FindBooks(limit int, offset int, query string) ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()

	// Попробуем преобразовать запрос в число для поиска по ISBN
	var isbn int
	if _, err := fmt.Sscanf(query, "%d", &isbn); err == nil {
		// Если запрос - число, то ищем по ISBN
		if err := tx.Model(&domain.Book{}).Where(
			"isbn = ?", isbn,
		).Preload("Author").Preload("Genre").Preload("Publisher").
			Limit(limit).Offset(offset).Find(&books).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Если нашли книги по ISBN, возвращаем их
		if len(books) > 0 {
			tx.Commit()
			return books, nil
		}
	}

	// Ищем книги по названию и по авторам (имя или фамилия содержит запрос)
	if err := tx.Model(&domain.Book{}).
		Joins("JOIN authors ON books.author_id = authors.id").
		Where(
			"lower(books.title) LIKE lower(?) OR lower(authors.name) LIKE lower(?) OR lower(authors.lastname) LIKE lower(?)",
			"%"+query+"%", "%"+query+"%", "%"+query+"%",
		).
		Preload("Author").Preload("Genre").Preload("Publisher").
		Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return books, nil
}

func (b BookRepository) CountBooksMatchingSearch(query string) (int, error) {
	var count int64
	tx := b.DB.Begin()

	if err := tx.Model(&domain.Book{}).
		Joins("JOIN authors ON books.author_id = authors.id").
		Where(
			"lower(books.title) LIKE lower(?) OR lower(authors.name) LIKE lower(?) OR lower(authors.lastname) LIKE lower(?)",
			"%"+query+"%", "%"+query+"%", "%"+query+"%",
		).
		Count(&count).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	err := tx.Commit().Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (b BookRepository) FindBooksWithFilters(
	limit int,
	offset int,
	query string,
	authorID uint,
	genreID uint,
	yearStart, yearEnd time.Time,
	sortBy string,
	sortOrder string,
) ([]*domain.Book, error) {
	var books []*domain.Book

	tx := b.DB.Begin()
	db := tx.Model(&domain.Book{}).Preload("Author").Preload("Genre").Preload("Publisher")

	// Применяем фильтры, если они указаны
	if query != "" {
		// Попробуем преобразовать запрос в число для поиска по ISBN
		var isbn int
		if _, err := fmt.Sscanf(query, "%d", &isbn); err == nil {
			db = db.Where("isbn = ?", isbn)
		} else {
			db = db.Joins("JOIN authors ON books.author_id = authors.id").
				Where(
					"lower(books.title) LIKE lower(?) OR lower(authors.name) LIKE lower(?) OR lower(authors.lastname) LIKE lower(?)",
					"%"+query+"%", "%"+query+"%", "%"+query+"%",
				)
		}
	}

	if authorID > 0 {
		db = db.Where("author_id = ?", authorID)
	}

	if genreID > 0 {
		db = db.Where("genre_id = ?", genreID)
	}

	// Фильтр по диапазону дат (если даты установлены)
	defaultTime := time.Time{}
	if yearStart != defaultTime {
		db = db.Where("year_of_publication >= ?", yearStart)
	}
	if yearEnd != defaultTime {
		db = db.Where("year_of_publication <= ?", yearEnd)
	}

	// Применяем сортировку
	if sortBy != "" {
		if sortOrder != "asc" && sortOrder != "desc" {
			sortOrder = "asc" // По умолчанию сортируем по возрастанию
		}

		switch sortBy {
		case "title":
			db = db.Order("books.title " + sortOrder)
		case "rating":
			db = db.Order("books.rating " + sortOrder)
		case "year":
			db = db.Order("books.year_of_publication " + sortOrder)
		default:
			db = db.Order("books.id " + sortOrder) // По умолчанию сортируем по ID
		}
	}

	// Применяем пагинацию
	if err := db.Limit(limit).Offset(offset).Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return books, nil
}

func (b BookRepository) CountBooksWithFilters(
	query string,
	authorID uint,
	genreID uint,
	yearStart, yearEnd time.Time,
) (int, error) {
	var count int64

	tx := b.DB.Begin()
	db := tx.Model(&domain.Book{})

	// Применяем фильтры, если они указаны
	if query != "" {
		var isbn int
		if _, err := fmt.Sscanf(query, "%d", &isbn); err == nil {
			db = db.Where("isbn = ?", isbn)
		} else {
			db = db.Joins("JOIN authors ON books.author_id = authors.id").
				Where(
					"lower(books.title) LIKE lower(?) OR lower(authors.name) LIKE lower(?) OR lower(authors.lastname) LIKE lower(?)",
					"%"+query+"%", "%"+query+"%", "%"+query+"%",
				)
		}
	}

	if authorID > 0 {
		db = db.Where("author_id = ?", authorID)
	}

	if genreID > 0 {
		db = db.Where("genre_id = ?", genreID)
	}

	// Фильтр по диапазону дат
	defaultTime := time.Time{}
	if yearStart != defaultTime {
		db = db.Where("year_of_publication >= ?", yearStart)
	}
	if yearEnd != defaultTime {
		db = db.Where("year_of_publication <= ?", yearEnd)
	}

	if err := db.Count(&count).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	err := tx.Commit().Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (b BookRepository) GetTopBooksByReservations(limit int, periodMonths int) ([]*domain.Book, error) {
	var books []*domain.Book

	// Вычисляем дату начала периода (например, 3 месяца назад)
	periodStart := time.Now().AddDate(0, -periodMonths, 0)

	tx := b.DB.Begin()

	// SQL запрос для получения топ-книг по количеству бронирований за указанный период
	if err := tx.Model(&domain.Book{}).
		Select("books.*, COUNT(reservations.id) as reservation_count").
		Joins("JOIN reservations ON books.id = reservations.book_id").
		Where("reservations.date_of_issue >= ?", periodStart).
		Group("books.id").
		Order("reservation_count DESC").
		Preload("Author").
		Preload("Genre").
		Preload("Publisher").
		Limit(limit).
		Find(&books).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return books, nil
}
