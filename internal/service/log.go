package service

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
)

type LogService struct {
	repository repository.LogRepo
}

func NewLogService(repository repository.LogRepo) *LogService {
	return &LogService{repository: repository}
}

func (l LogService) GetLogByID(id uint) (*domain.Log, error) {
	return l.repository.GetLogByID(id)
}

func (l LogService) GetLogs() ([]*domain.Log, error) {
	return l.repository.GetLogs()
}

func (l LogService) CreateLog(logInput *CreateLogInput) error {
	log := &domain.Log{
		UserID:  logInput.UserID,
		Action:  logInput.Action,
		Date:    logInput.Date,
		Details: logInput.Details,
	}

	return l.repository.CreateLog(log)
}

// Вспомогательный метод для создания лога с простыми параметрами (для обратной совместимости)
func (l LogService) CreateSimpleLog(userID uint, action string) error {
	logInput := &CreateLogInput{
		UserID:  userID,
		Action:  action,
		Date:    time.Now(),
		Details: action,
	}
	return l.CreateLog(logInput)
}

func (l LogService) UpdateLog(logInput *UpdateLogInput) error {
	log := &domain.Log{
		ID:      logInput.ID,
		UserID:  logInput.UserID,
		Action:  logInput.Action,
		Date:    logInput.Date,
		Details: logInput.Details,
	}

	return l.repository.UpdateLog(log)
}

func (l LogService) DeleteLog(id uint) error {
	return l.repository.DeleteLog(id)
}

func (l LogService) GetLogsByUserID(id uint) ([]*domain.Log, error) {
	return l.repository.GetLogsByUserID(id)
}

func (l LogService) CreateLogWithCookie(cookie *http.Cookie, action string) error {
	if cookie.Value == "" {
		return nil
	}

	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return err
	}

	logInput := &CreateLogInput{
		UserID:  uint(userID), //nolint:gosec
		Action:  action,
		Date:    time.Now(),
		Details: action,
	}

	return l.CreateLog(logInput)
}
