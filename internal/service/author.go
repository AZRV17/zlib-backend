package service

import (
	"github.com/AZRV17/zlib-backend/internal/domain"
	"github.com/AZRV17/zlib-backend/internal/repository"
)

type AuthorService struct {
	repository repository.AuthorRepo
}

func NewAuthorService(repository repository.AuthorRepo) *AuthorService {
	return &AuthorService{repository: repository}
}

func (a AuthorService) GetAuthorByID(id uint) (*domain.Author, error) {
	return a.repository.GetAuthorByID(id)
}

func (a AuthorService) GetAuthors() ([]*domain.Author, error) {
	return a.repository.GetAuthors()
}

func (a AuthorService) CreateAuthor(authorInput *CreateAuthorInput) error {
	author := domain.Author{
		Name:      authorInput.Name,
		Lastname:  authorInput.Lastname,
		Biography: authorInput.Biography,
		Birthdate: authorInput.Birthdate,
	}

	if err := a.repository.CreateAuthor(&author); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) UpdateAuthor(authorInput *UpdateAuthorInput) error {
	author := domain.Author{
		ID:        authorInput.ID,
		Name:      authorInput.Name,
		Lastname:  authorInput.Lastname,
		Biography: authorInput.Biography,
		Birthdate: authorInput.Birthdate,
	}

	if err := a.repository.UpdateAuthor(&author); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) DeleteAuthor(id uint) error {
	if err := a.repository.DeleteAuthor(id); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) GetAuthorBooks(id uint) ([]*domain.Book, error) {
	return a.repository.GetAuthorBooks(id)
}

func (a AuthorService) CreateAuthorBook(authorBook *domain.AuthorBook) error {

	if err := a.repository.CreateAuthorBook(authorBook); err != nil {
		return err
	}

	return nil
}

func (a AuthorService) DeleteAuthorBook(id uint) error {
	if err := a.repository.DeleteAuthorBook(id); err != nil {
		return err
	}

	return nil
}
