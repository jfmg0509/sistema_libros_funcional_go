package usecase

import (
	"fmt"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

/*
   ==========================================================
   BookService
   ==========================================================

   Este servicio maneja la lógica de negocio de los libros y
   de los accesos a los libros (eventos de lectura, apertura, etc.).

   Depende de TRES interfaces:
   - BookRepository: para guardar y buscar libros.
   - UserRepository: para verificar que el usuario exista.
   - AccessLogRepository: para guardar eventos de acceso.
*/

// BookService representa los casos de uso relacionados con libros.
type BookService struct {
	bookRepo      domain.BookRepository
	userRepo      domain.UserRepository
	accessLogRepo domain.AccessLogRepository
}

// NewBookService es el CONSTRUCTOR de BookService.
func NewBookService(
	bookRepo domain.BookRepository,
	userRepo domain.UserRepository,
	accessLogRepo domain.AccessLogRepository,
) *BookService {
	return &BookService{
		bookRepo:      bookRepo,
		userRepo:      userRepo,
		accessLogRepo: accessLogRepo,
	}
}

/*
RegisterBook registra un nuevo libro en el sistema.

Pasos:
1. Usa el constructor de dominio (NewBook) para validar los datos.
2. Pide al repositorio que cree el libro.
3. Devuelve el libro creado.
*/
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

/*
SearchBooks permite buscar libros por filtros.

Recibe un BookFilter que puede tener:
- parte del título
- parte del autor
- categoría
- años "from" y "to"
- tags

La implementación exacta del filtro se hace en el repositorio.
*/
func (s *BookService) SearchBooks(filter domain.BookFilter) ([]*domain.Book, error) {
	return s.bookRepo.SearchByFilters(filter)
}

/*
RecordAccess registra un acceso de un usuario a un libro.

Pasos:
1. Verificar que el libro exista.
2. Verificar que el usuario exista.
3. Crear un AccessEvent (dominio).
4. Guardar el evento en el AccessLogRepository.
*/
func (s *BookService) RecordAccess(
	bookID domain.BookID,
	userID domain.UserID,
	accessType domain.AccessType,
) error {

	// 1. Verificar libro.
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return err
	}
	if book == nil {
		return fmt.Errorf("libro no encontrado")
	}

	// 2. Verificar usuario.
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("usuario no encontrado")
	}

	// 3. Crear el evento de acceso.
	event, err := domain.NewAccessEvent(bookID, userID, accessType)
	if err != nil {
		return err
	}

	// 4. Guardar el evento.
	if err := s.accessLogRepo.Store(event); err != nil {
		return err
	}

	return nil
}

/*
BuildAccessStatsByBook genera estadísticas de accesos por libro.

Devuelve un MAP:
- clave: AccessType (LECTURA, APERTURA, etc.)
- valor: cantidad de veces que ocurrió

Ejemplo de retorno:

	{
	  "LECTURA":  3,
	  "APERTURA": 5
	}
*/
func (s *BookService) BuildAccessStatsByBook(bookID domain.BookID) (map[domain.AccessType]int, error) {
	// 1. Traer todos los eventos de ese libro.
	events, err := s.accessLogRepo.ListByBook(bookID)
	if err != nil {
		return nil, err
	}

	// 2. Crear el MAP donde vamos a contar.
	stats := make(map[domain.AccessType]int)

	// 3. Recorrer eventos e incrementar el contador por tipo.
	for _, ev := range events {
		stats[ev.AccessType()] = stats[ev.AccessType()] + 1
	}

	return stats, nil
}
