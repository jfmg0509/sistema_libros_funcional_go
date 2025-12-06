package db

import (
	"errors"
	"strings"
	"sync"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

/*
   ==========================================================
   InMemoryUserRepo
   ==========================================================

   Implementación EN MEMORIA de UserRepository usando MAPS.

   - users:      map[UserID]*User
   - emailIndex: map[string]UserID (para buscar rápido por email)

   Ideal para prácticas y prototipos sin base de datos real.
*/

// InMemoryUserRepo implementa domain.UserRepository usando mapas en memoria.
type InMemoryUserRepo struct {
	mu         sync.RWMutex  // mutex para acceso concurrente
	seq        domain.UserID // secuencia para generar IDs
	users      map[domain.UserID]*domain.User
	emailIndex map[string]domain.UserID
}

// NewInMemoryUserRepo crea un repositorio vacío listo para usar.
func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users:      make(map[domain.UserID]*domain.User),
		emailIndex: make(map[string]domain.UserID),
	}
}

// nextID incrementa la secuencia y devuelve un nuevo ID de usuario.
func (r *InMemoryUserRepo) nextID() domain.UserID {
	r.seq++
	return r.seq
}

// Create guarda un nuevo usuario en el mapa.
// Valida que no exista otro usuario con el mismo email.
func (r *InMemoryUserRepo) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.emailIndex[user.Email()]; exists {
		return errors.New("ya existe un usuario con ese email")
	}

	id := r.nextID()
	user.SetID(id)

	r.users[id] = user
	r.emailIndex[user.Email()] = id
	return nil
}

// Update actualiza un usuario ya existente.
func (r *InMemoryUserRepo) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID() == 0 {
		return errors.New("el usuario no tiene ID asignado")
	}

	if _, exists := r.users[user.ID()]; !exists {
		return errors.New("no existe un usuario con ese ID")
	}

	// Actualizar el índice de email si cambió el correo.
	for email, id := range r.emailIndex {
		if id == user.ID() && email != user.Email() {
			delete(r.emailIndex, email)
			r.emailIndex[user.Email()] = user.ID()
			break
		}
	}

	r.users[user.ID()] = user
	return nil
}

// FindByID busca un usuario por su ID.
func (r *InMemoryUserRepo) FindByID(id domain.UserID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

// FindByEmail busca un usuario por su email usando el índice.
func (r *InMemoryUserRepo) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.emailIndex[strings.ToLower(email)]
	if !ok {
		return nil, nil
	}
	user, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

// ListAll devuelve todos los usuarios en una slice.
func (r *InMemoryUserRepo) ListAll() ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

/*
   ==========================================================
   InMemoryBookRepo
   ==========================================================
*/

// InMemoryBookRepo implementa BookRepository usando mapas en memoria.
type InMemoryBookRepo struct {
	mu    sync.RWMutex
	seq   domain.BookID
	books map[domain.BookID]*domain.Book
}

// NewInMemoryBookRepo crea un repositorio de libros en memoria.
func NewInMemoryBookRepo() *InMemoryBookRepo {
	return &InMemoryBookRepo{
		books: make(map[domain.BookID]*domain.Book),
	}
}

// nextID genera un nuevo ID para libros.
func (r *InMemoryBookRepo) nextID() domain.BookID {
	r.seq++
	return r.seq
}

// Create guarda un nuevo libro en el mapa.
func (r *InMemoryBookRepo) Create(book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID()
	book.SetID(id)
	r.books[id] = book
	return nil
}

// Update actualiza un libro existente.
func (r *InMemoryBookRepo) Update(book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if book.ID() == 0 {
		return errors.New("el libro no tiene ID asignado")
	}
	if _, exists := r.books[book.ID()]; !exists {
		return errors.New("no existe un libro con ese ID")
	}

	r.books[book.ID()] = book
	return nil
}

// FindByID busca un libro por su ID.
func (r *InMemoryBookRepo) FindByID(id domain.BookID) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.books[id]
	if !ok {
		return nil, nil
	}
	return book, nil
}

// SearchByFilters aplica filtros básicos sobre todos los libros.
func (r *InMemoryBookRepo) SearchByFilters(filter domain.BookFilter) ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Book, 0)
	for _, b := range r.books {
		// Solo libros activos.
		if !b.Active() {
			continue
		}

		// Filtro por título.
		if filter.TitleContains != "" &&
			!strings.Contains(strings.ToLower(b.Title()), strings.ToLower(filter.TitleContains)) {
			continue
		}

		// Filtro por autor.
		if filter.AuthorContains != "" &&
			!strings.Contains(strings.ToLower(b.Author()), strings.ToLower(filter.AuthorContains)) {
			continue
		}

		// Filtro por categoría TI.
		if filter.CategoryTI != "" &&
			!strings.EqualFold(b.CategoryTI(), filter.CategoryTI) {
			continue
		}

		// Filtro por año desde.
		if filter.YearFrom > 0 && b.Year() < filter.YearFrom {
			continue
		}

		// Filtro por año hasta.
		if filter.YearTo > 0 && b.Year() > filter.YearTo {
			continue
		}

		// (Opcional) Filtro por tags: aquí podrías validar si contiene ciertos tags.

		result = append(result, b)
	}

	return result, nil
}

// ListAll devuelve todos los libros.
func (r *InMemoryBookRepo) ListAll() ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Book, 0, len(r.books))
	for _, b := range r.books {
		result = append(result, b)
	}
	return result, nil
}

/*
   ==========================================================
   InMemoryAccessLogRepo
   ==========================================================
*/

// InMemoryAccessLogRepo implementa AccessLogRepository en memoria.
type InMemoryAccessLogRepo struct {
	mu     sync.RWMutex
	seq    domain.AccessEventID
	events map[domain.AccessEventID]*domain.AccessEvent
}

// NewInMemoryAccessLogRepo crea un repositorio de accesos vacío.
func NewInMemoryAccessLogRepo() *InMemoryAccessLogRepo {
	return &InMemoryAccessLogRepo{
		events: make(map[domain.AccessEventID]*domain.AccessEvent),
	}
}

// nextID genera un nuevo ID para eventos de acceso.
func (r *InMemoryAccessLogRepo) nextID() domain.AccessEventID {
	r.seq++
	return r.seq
}

// Store guarda un nuevo evento de acceso.
func (r *InMemoryAccessLogRepo) Store(event *domain.AccessEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID()
	event.SetID(id)
	r.events[id] = event
	return nil
}

// ListByBook devuelve todos los eventos para un libro.
func (r *InMemoryAccessLogRepo) ListByBook(bookID domain.BookID) ([]*domain.AccessEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.AccessEvent, 0)
	for _, ev := range r.events {
		if ev.BookID() == bookID {
			result = append(result, ev)
		}
	}
	return result, nil
}

// ListByUser devuelve todos los eventos para un usuario (no lo usamos aún, pero está listo).
func (r *InMemoryAccessLogRepo) ListByUser(userID domain.UserID) ([]*domain.AccessEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.AccessEvent, 0)
	for _, ev := range r.events {
		if ev.UserID() == userID {
			result = append(result, ev)
		}
	}
	return result, nil
}
