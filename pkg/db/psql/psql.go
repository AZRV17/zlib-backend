package psql

import (
	"bytes"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
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

// BackupDatabase создает полный бэкап базы данных в файл
func BackupDatabase(dsn string) error {
	backupFile := fmt.Sprintf("./backups/backup.sql")

	cmd := exec.Command("pg_dump", dsn, "-f", backupFile)
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		log.Println(cmd)
		log.Printf("Ошибка создания бэкапа: %v", err)
		return err
	}

	log.Printf("Бэкап базы данных создан: %s", backupFile)
	return nil
}

// DropAndCreateDatabase удаляет существующую базу данных и создает ее заново
func DropAndCreateDatabase(host, port, user, password, dbName string) error {
	// Устанавливаем переменные среды для подключения
	os.Setenv("PGHOST", host)
	os.Setenv("PGPORT", port)
	os.Setenv("PGUSER", user)
	os.Setenv("PGPASSWORD", password)

	// Подключаемся к базе по умолчанию "postgres"
	defaultDB := "postgres"

	// Удаление базы данных
	var stderr bytes.Buffer
	dropCmd := exec.Command("psql", defaultDB, "-c", fmt.Sprintf(`DROP DATABASE IF EXISTS "%s";`, dbName))
	dropCmd.Env = os.Environ()
	dropCmd.Stderr = &stderr

	if err := dropCmd.Run(); err != nil {
		log.Printf("Ошибка удаления базы данных: %v, %s", err, stderr.String())
		return err
	}
	log.Printf("База данных %s удалена", dbName)

	// Создание базы данных заново
	createCmd := exec.Command("psql", defaultDB, "-c", fmt.Sprintf(`CREATE DATABASE "%s";`, dbName))
	createCmd.Env = os.Environ()
	if err := createCmd.Run(); err != nil {
		log.Printf("Ошибка создания базы данных: %v", err)
		return err
	}
	log.Printf("База данных %s создана заново", dbName)
	return nil
}

// RestoreDatabase восстанавливает базу данных из бэкап файла
func RestoreDatabase(host, port, user, password, dbName string) error {
	// Удаляем и создаем базу данных заново
	if err := DropAndCreateDatabase(host, port, user, password, "testDB"); err != nil {
		return err
	}

	// Устанавливаем базу данных для восстановления
	os.Setenv("PGDATABASE", "testDB")

	// Восстанавливаем данные из бэкапа
	restoreCmd := exec.Command("psql", "-f", "./backups/backup.sql")
	restoreCmd.Env = os.Environ()

	if err := restoreCmd.Run(); err != nil {
		log.Printf("Ошибка восстановления базы данных: %v", err)
		return err
	}

	log.Printf("База данных %s успешно восстановлена из бэкапа: %s", "testDB", "./backups/backup.sql")
	return nil
}
