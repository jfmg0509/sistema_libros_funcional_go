package db

import (
	"errors"
	"strings"
	"sync"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

/*
   ======================
   REPOSITORIO EN MEMORIA PARA USUARIOS
   ======================
*/

type InMemoryUserRepo struct {
	mu         sync.RWMutex
	seq        domain.UserID
	users      map[domain.UserID]*domain.User
	emailIndex map[string]domain.UserID
}

// NewInMemoryUserRepo crea un repositorio de usuarios en memoria.
func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users:      make(map[domain.UserID]*domain.User),
		emailIndex: make(map[string]domain.UserID),
	}
}

func (r *InMemoryUserRepo) nextID() domain.UserID {
	r.seq++
	return r.seq
}

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

func (r *InMemoryUserRepo) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID()]; !ok {
		return errors.New("usuario no encontrado para actualizar")
	}
	r.users[user.ID()] = user
	r.emailIndex[user.Email()] = user.ID()
	return nil
}

func (r *InMemoryUserRepo) FindByID(id domain.UserID) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (r *InMemoryUserRepo) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.emailIndex[email]
	if !ok {
		return nil, nil
	}
	return r.users[id], nil
}

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
   ======================
   REPOSITORIO EN MEMORIA PARA LIBROS
   ======================
*/

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

func (r *InMemoryBookRepo) nextID() domain.BookID {
	r.seq++
	return r.seq
}

func (r *InMemoryBookRepo) Create(book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID()
	book.SetID(id)

	r.books[id] = book
	return nil
}

func (r *InMemoryBookRepo) Update(book *domain.Book) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.books[book.ID()]; !ok {
		return errors.New("libro no encontrado para actualizar")
	}
	r.books[book.ID()] = book
	return nil
}

func (r *InMemoryBookRepo) FindByID(id domain.BookID) (*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	book, ok := r.books[id]
	if !ok {
		return nil, nil
	}
	return book, nil
}

// SearchByFilters aplica filtros simples EN MEMORIA.
func (r *InMemoryBookRepo) SearchByFilters(filter domain.BookFilter) ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Book, 0)
	for _, b := range r.books {
		if !b.Active() {
			continue
		}

		// Filtros por título, autor, categoría, años...
		if filter.TitleContains != "" &&
			!strings.Contains(strings.ToLower(b.Title()), strings.ToLower(filter.TitleContains)) {
			continue
		}

		if filter.AuthorContains != "" &&
			!strings.Contains(strings.ToLower(b.Author()), strings.ToLower(filter.AuthorContains)) {
			continue
		}

		if filter.CategoryTI != "" &&
			!strings.EqualFold(b.CategoryTI(), filter.CategoryTI) {
			continue
		}

		if filter.YearFrom > 0 && b.Year() < filter.YearFrom {
			continue
		}
		if filter.YearTo > 0 && b.Year() > filter.YearTo {
			continue
		}

		// Filtro por tags (si se enviaron).
		if len(filter.Tags) > 0 && !bookHasAllTags(b, filter.Tags) {
			continue
		}

		result = append(result, b)
	}

	return result, nil
}

func (r *InMemoryBookRepo) ListAll() ([]*domain.Book, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*domain.Book, 0, len(r.books))
	for _, b := range r.books {
		result = append(result, b)
	}
	return result, nil
}

// bookHasAllTags ayuda a verificar si un libro contiene TODOS los tags del filtro.
func bookHasAllTags(b *domain.Book, tags []string) bool {
	tagSet := make(map[string]bool)
	for _, t := range b.Tags() {
		tagSet[strings.ToLower(t)] = true
	}

	for _, filterTag := range tags {
		if !tagSet[strings.ToLower(filterTag)] {
			return false
		}
	}
	return true
}

/*
   ==========================
   REPOSITORIO EN MEMORIA PARA ACCESS LOGS
   ==========================
*/

type InMemoryAccessLogRepo struct {
	mu     sync.RWMutex
	seq    domain.AccessEventID
	events map[domain.AccessEventID]*domain.AccessEvent
}

// NewInMemoryAccessLogRepo crea un repositorio en memoria para logs de acceso.
func NewInMemoryAccessLogRepo() *InMemoryAccessLogRepo {
	return &InMemoryAccessLogRepo{
		events: make(map[domain.AccessEventID]*domain.AccessEvent),
	}
}

func (r *InMemoryAccessLogRepo) nextID() domain.AccessEventID {
	r.seq++
	return r.seq
}

func (r *InMemoryAccessLogRepo) Store(event *domain.AccessEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := r.nextID()
	event.SetID(id)

	r.events[id] = event
	return nil
}

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
