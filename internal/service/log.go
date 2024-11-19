package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
	"net/http"
	"strconv"
	"time"
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

	log := &domain.Log{
		UserID:  uint(userID), //nolint:gosec
		Action:  action,
		Date:    time.Now(),
		Details: action,
	}

	return l.repository.CreateLog(log)
}
