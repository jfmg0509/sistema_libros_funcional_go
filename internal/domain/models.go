package domain

import (
	"errors"
	"strings"
	"time"
)

/*
   ==========================================================
   TIPOS BÁSICOS DEL DOMINIO
   ==========================================================

   Aquí definimos tipos propios para darle más significado
   a los datos y no usar solo int64 o string en todo lado.
*/

// UserID representa el identificador único de un usuario.
type UserID int64

// BookID representa el identificador único de un libro.
type BookID int64

// AccessEventID representa el identificador único de un evento de acceso.
type AccessEventID int64

// Role representa el rol de un usuario dentro del sistema.
type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleReader Role = "READER"
)

// allowedRoles es un ARRAY con los roles permitidos.
// Se usa para validar en NewUser y ChangeRole.
var allowedRoles = [2]Role{RoleAdmin, RoleReader}

// AccessType representa el tipo de acceso a un libro.
type AccessType string

const (
	AccessTypeApertura AccessType = "APERTURA"
	AccessTypeLectura  AccessType = "LECTURA"
	AccessTypeDescarga AccessType = "DESCARGA"
)

/*
   ==========================================================
   ENTIDAD: USER (USUARIO)
   ==========================================================
*/

// User representa a un usuario del sistema.
// Notar que todos los campos son privados (inician en minúscula).
// Esto es ENCAPSULACIÓN: solo se accede a ellos mediante métodos.
type User struct {
	id        UserID
	name      string
	email     string
	role      Role
	active    bool
	createdAt time.Time
}

// NewUser es un CONSTRUCTOR de usuarios.
// Recibe los datos, valida y devuelve (*User, error).
func NewUser(name, email string, role Role) (*User, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	if strings.TrimSpace(email) == "" {
		return nil, errors.New("el email no puede estar vacío")
	}
	if !isValidRole(role) {
		return nil, errors.New("rol no válido")
	}

	return &User{
		id:        0, // se asigna luego en el repositorio
		name:      name,
		email:     email,
		role:      role,
		active:    true,
		createdAt: time.Now(),
	}, nil
}

// isValidRole revisa si el rol está dentro del ARRAY allowedRoles.
func isValidRole(r Role) bool {
	for _, allowed := range allowedRoles {
		if allowed == r {
			return true
		}
	}
	return false
}

// Métodos GETTER para leer los campos privados.

func (u *User) ID() UserID           { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() string        { return u.email }
func (u *User) Role() Role           { return u.role }
func (u *User) Active() bool         { return u.active }
func (u *User) CreatedAt() time.Time { return u.createdAt }

// SetID permite asignar el ID desde el repositorio.
func (u *User) SetID(id UserID) {
	u.id = id
}

// ChangeRole cambia el rol del usuario, validando que el nuevo rol sea permitido.
func (u *User) ChangeRole(newRole Role) error {
	if !isValidRole(newRole) {
		return errors.New("nuevo rol no válido")
	}
	u.role = newRole
	return nil
}

// Deactivate marca al usuario como inactivo (no borramos, solo inactivamos).
func (u *User) Deactivate() {
	u.active = false
}

/*
   ==========================================================
   ENTIDAD: BOOK (LIBRO)
   ==========================================================
*/

// Book representa un libro electrónico del sistema.
type Book struct {
	id         BookID
	title      string
	author     string
	year       int
	isbn       string
	categoryTI string
	tags       []string
	active     bool
	createdAt  time.Time
}

// NewBook es el CONSTRUCTOR de libros.
// Valida los datos de entrada y devuelve (*Book, error).
func NewBook(title, author string, year int, isbn, categoryTI string, tags []string) (*Book, error) {
	if strings.TrimSpace(title) == "" {
		return nil, errors.New("el título no puede estar vacío")
	}
	if strings.TrimSpace(author) == "" {
		return nil, errors.New("el autor no puede estar vacío")
	}
	if year <= 0 {
		return nil, errors.New("el año debe ser mayor que cero")
	}
	if strings.TrimSpace(isbn) == "" {
		return nil, errors.New("el ISBN no puede estar vacío")
	}
	if strings.TrimSpace(categoryTI) == "" {
		return nil, errors.New("la categoría TI no puede estar vacía")
	}

	return &Book{
		id:         0, // se asigna luego en el repositorio
		title:      title,
		author:     author,
		year:       year,
		isbn:       isbn,
		categoryTI: categoryTI,
		tags:       tags,
		active:     true,
		createdAt:  time.Now(),
	}, nil
}

// Getters del libro.

func (b *Book) ID() BookID           { return b.id }
func (b *Book) Title() string        { return b.title }
func (b *Book) Author() string       { return b.author }
func (b *Book) Year() int            { return b.year }
func (b *Book) ISBN() string         { return b.isbn }
func (b *Book) CategoryTI() string   { return b.categoryTI }
func (b *Book) Tags() []string       { return b.tags }
func (b *Book) Active() bool         { return b.active }
func (b *Book) CreatedAt() time.Time { return b.createdAt }

// SetID asigna el ID del libro desde el repositorio.
func (b *Book) SetID(id BookID) {
	b.id = id
}

// Archive marca el libro como no activo (archivado).
func (b *Book) Archive() {
	b.active = false
}

/*
   ==========================================================
   ENTIDAD: ACCESSEVENT (REGISTRO DE ACCESO)
   ==========================================================
*/

// AccessEvent representa un registro de acceso de un usuario a un libro.
type AccessEvent struct {
	id         AccessEventID
	bookID     BookID
	userID     UserID
	accessType AccessType
	timestamp  time.Time
}

// NewAccessEvent crea un nuevo evento de acceso validando datos.
func NewAccessEvent(bookID BookID, userID UserID, accessType AccessType) (*AccessEvent, error) {
	if bookID <= 0 {
		return nil, errors.New("bookID debe ser mayor que cero")
	}
	if userID <= 0 {
		return nil, errors.New("userID debe ser mayor que cero")
	}
	if accessType == "" {
		return nil, errors.New("el tipo de acceso no puede estar vacío")
	}

	return &AccessEvent{
		id:         0, // se asigna en el repositorio
		bookID:     bookID,
		userID:     userID,
		accessType: accessType,
		timestamp:  time.Now(),
	}, nil
}

// Getters del evento de acceso.

func (e *AccessEvent) ID() AccessEventID      { return e.id }
func (e *AccessEvent) BookID() BookID         { return e.bookID }
func (e *AccessEvent) UserID() UserID         { return e.userID }
func (e *AccessEvent) AccessType() AccessType { return e.accessType }
func (e *AccessEvent) Timestamp() time.Time   { return e.timestamp }

// SetID asigna el ID desde el repositorio.
func (e *AccessEvent) SetID(id AccessEventID) {
	e.id = id
}

/*
   ==========================================================
   FILTRO DE BÚSQUEDA DE LIBROS
   ==========================================================
*/

// BookFilter se usa para filtrar libros por distintos criterios.
type BookFilter struct {
	TitleContains  string
	AuthorContains string
	CategoryTI     string
	YearFrom       int
	YearTo         int
	Tags           []string
}

/*
   ==========================================================
   INTERFACES DE REPOSITORIO
   ==========================================================
*/

// UserRepository define las operaciones que se pueden hacer con usuarios.
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	FindByID(id UserID) (*User, error)
	FindByEmail(email string) (*User, error)
	ListAll() ([]*User, error)
}

// BookRepository define las operaciones de persistencia de libros.
type BookRepository interface {
	Create(book *Book) error
	Update(book *Book) error
	FindByID(id BookID) (*Book, error)
	SearchByFilters(filter BookFilter) ([]*Book, error)
	ListAll() ([]*Book, error)
}

// AccessLogRepository define cómo se guardan los eventos de acceso.
type AccessLogRepository interface {
	Store(event *AccessEvent) error
	ListByBook(bookID BookID) ([]*AccessEvent, error)
	ListByUser(userID UserID) ([]*AccessEvent, error)
}
