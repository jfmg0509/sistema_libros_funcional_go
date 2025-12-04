package usecase

import (
	"errors"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

// BookService contiene la lógica de negocio para libros y accesos.
type BookService struct {
	bookRepo   domain.BookRepository
	accessRepo domain.AccessLogRepository
}

// NewBookService crea un nuevo BookService.
func NewBookService(bookRepo domain.BookRepository, accessRepo domain.AccessLogRepository) *BookService {
	return &BookService{
		bookRepo:   bookRepo,
		accessRepo: accessRepo,
	}
}

// RegisterBook registra un nuevo libro en el sistema.
func (s *BookService) RegisterBook(
	title, author string,
	year int,
	isbn, categoryTI string,
	tags []string,
) (*domain.Book, error) {
	book, err := domain.NewBook(title, author, year, isbn, categoryTI, tags)
	if err != nil {
		return nil, err
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}
	return book, nil
}

// SearchBooks realiza búsquedas avanzadas de libros según un filtro.
func (s *BookService) SearchBooks(filter domain.BookFilter) ([]*domain.Book, error) {
	return s.bookRepo.SearchByFilters(filter)
}

// ArchiveBook archiva (desactiva) un libro.
func (s *BookService) ArchiveBook(id domain.BookID) error {
	book, err := s.bookRepo.FindByID(id)
	if err != nil {
		return err
	}
	if book == nil {
		return errors.New("libro no encontrado")
	}

	book.Archive()
	return s.bookRepo.Update(book)
}

// RecordAccess registra un evento de acceso a un libro.
func (s *BookService) RecordAccess(bookID domain.BookID, userID domain.UserID, access domain.AccessType) error {
	event, err := domain.NewAccessEvent(bookID, userID, access)
	if err != nil {
		return err
	}
	return s.accessRepo.Store(event)
}

// BuildAccessStatsByBook genera estadísticas de accesos por tipo de acceso.
// Aquí usamos un MAP (map[AccessType]int) cumpliendo el requerimiento.
func (s *BookService) BuildAccessStatsByBook(bookID domain.BookID) (map[domain.AccessType]int, error) {
	events, err := s.accessRepo.ListByBook(bookID)
	if err != nil {
		return nil, err
	}

	stats := make(map[domain.AccessType]int)
	for _, ev := range events {
		stats[ev.Access()]++
	}
	return stats, nil
}
