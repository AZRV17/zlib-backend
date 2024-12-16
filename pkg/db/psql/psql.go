package psql

import (
	"bytes"
	"fmt"
	"github.com/AZRV17/zlib-backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
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

	// Создание ENUM типа для ролей (если его еще нет)
	DB.Exec("CREATE TYPE role AS ENUM ('user', 'admin', 'librarian');")

	// Создание функции sign_in
	DB.Exec(
		`
        CREATE OR REPLACE FUNCTION sign_in(login TEXT, pass TEXT) 
        RETURNS TABLE(id INT, role TEXT) AS $$
        BEGIN
            RETURN QUERY
            SELECT id, role
            FROM users
            WHERE login = login AND password = pass;

            IF NOT FOUND THEN
                RAISE EXCEPTION 'Invalid login or password';
            END IF;
        END;
        $$ LANGUAGE plpgsql;
    `,
	)

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

	DB.Exec("CREATE INDEX idx_book_title ON books USING gin (to_tsvector('simple', title));")
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

// DropAndCreateDatabase удаляет существующую базу данных и создает ее заново
func DropAndCreateDatabase(host, port, user, password, dbName string) error {
	// Устанавливаем переменные среды для подключения
	err := os.Setenv("PGHOST", host)
	if err != nil {
		return err
	}
	err = os.Setenv("PGPORT", port)
	if err != nil {
		return err
	}
	err = os.Setenv("PGUSER", user)
	if err != nil {
		return err
	}
	err = os.Setenv("PGPASSWORD", password)
	if err != nil {
		return err
	}

	// Подключаемся к базе по умолчанию "postgres"
	defaultDB := "postgres"

	// Удаление базы данных
	var stderr bytes.Buffer
	dropCmd := exec.Command(
		//nolint:gosec
		"psql",
		defaultDB,
		"-c",
		fmt.Sprintf(`DROP DATABASE IF EXISTS "%s";`, dbName),
	)
	dropCmd.Env = os.Environ()
	dropCmd.Stderr = &stderr

	if err := dropCmd.Run(); err != nil {
		log.Printf("Ошибка удаления базы данных: %v, %s", err, stderr.String())
		return err
	}
	log.Printf("База данных %s удалена", dbName)

	// Создание базы данных заново
	createCmd := exec.Command("psql", defaultDB, "-c", fmt.Sprintf(`CREATE DATABASE "%s";`, dbName)) //nolint:gosec
	createCmd.Env = os.Environ()
	if err := createCmd.Run(); err != nil {
		log.Printf("Ошибка создания базы данных: %v", err)
		return err
	}
	log.Printf("База данных %s создана заново", dbName)
	return nil
}

// BackupDatabase создает бэкап базы данных и возвращает его как []byte
func BackupDatabase(dsn string) ([]byte, error) {
	// Создаем буфер для хранения вывода pg_dump
	var outBuffer bytes.Buffer

	// Запускаем pg_dump и направляем вывод в буфер
	cmd := exec.Command("pg_dump", dsn)
	cmd.Stdout = &outBuffer
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		log.Printf("Ошибка создания бэкапа: %v", err)
		return nil, err
	}

	return outBuffer.Bytes(), nil
}

// RestoreDatabase восстанавливает базу данных из предоставленного ридера
func RestoreDatabase(host, port, user, password, dbName string, backupData io.Reader) error {
	// Удаляем и создаем базу данных заново
	if err := DropAndCreateDatabase(host, port, user, password, dbName); err != nil {
		return err
	}

	// Устанавливаем базу данных для восстановления
	err := os.Setenv("PGDATABASE", dbName)
	if err != nil {
		return err
	}

	// Создаем команду psql для восстановления
	cmd := exec.Command("psql")
	cmd.Env = os.Environ()

	// Присоединяем входные данные к stdin команды
	cmd.Stdin = backupData

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Ошибка восстановления базы данных: %v, %s", err, stderr.String())
		return fmt.Errorf("error restoring database: %w, %s", err, stderr.String())
	}

	log.Printf("База данных %s успешно восстановлена", dbName)
	return nil
}
