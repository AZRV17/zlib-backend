package psql

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

// Connect Подключение к БД
func Connect(dsn string) error {
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB.Exec("CREATE TYPE role AS ENUM ('user', 'admin', 'librarian');")

	err = DB.AutoMigrate(
		&domain.User{},
		&domain.Book{},
		&domain.UniqueCode{},
		&domain.Author{},
		&domain.Genre{},
		&domain.Publisher{},
		&domain.Favorite{},
		&domain.Review{},
		&domain.Reservation{},
		&domain.Notification{},
		&domain.Log{},
	)
	if err != nil {
		return err
	}

	return nil
}

// Close Закрытие соединения
func Close() {
	db, err := DB.DB()
	if err != nil {
		log.Fatal("error getting db", err)
	}

	log.Println("closing db")

	err = db.Close()
	if err != nil {
		log.Fatal("error closing db: ", err)
	}

	log.Println("db closed")
}
