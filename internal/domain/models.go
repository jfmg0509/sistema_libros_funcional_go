package domain

import (
	"errors"
	"strings"
	"time"
)

/*
   ================
   TIPOS Y CONSTANTES
   ================
*/

// Role representa el rol de un usuario dentro del sistema.
type Role string

const (
	RoleAdmin       Role = "ADMIN"
	RoleConsultorTI Role = "CONSULTOR_TI"
	RoleLectura     Role = "SOLO_LECTURA"
)

// Aquí usamos un ARRAY (tamaño fijo) para cumplir el requerimiento del ejercicio.
var allowedRoles = [3]Role{
	RoleAdmin,
	RoleConsultorTI,
	RoleLectura,
}

// Tipos para IDs, solo para darle más claridad al código.
type UserID int64
type BookID int64
type AccessEventID int64

/*
   ======================
   ESTRUCTURA: USER (USUARIO)
   ======================
*/

// User representa un usuario del sistema de libros.
type User struct {
	id        UserID
	name      string
	email     string
	role      Role
	active    bool
	createdAt time.Time
}

// NewUser es el CONSTRUCTOR de User.
// Aquí validamos los datos y devolvemos un *User o un error.
func NewUser(name, email string, role Role) (*User, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(strings.ToLower(email))

	if name == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("el email no es válido")
	}
	if !isValidRole(role) {
		return nil, errors.New("el rol indicado no es válido")
	}

	return &User{
		// el ID se asignará después en el repositorio
		name:      name,
		email:     email,
		role:      role,
		active:    true,
		createdAt: time.Now(),
	}, nil
}

// isValidRole revisa en el ARRAY allowedRoles si el rol existe.
func isValidRole(r Role) bool {
	for _, allowed := range allowedRoles {
		if allowed == r {
			return true
		}
	}
	return false
}

// ====== Métodos (encapsulación) ======

// SetID permite al repositorio asignar el ID.
func (u *User) SetID(id UserID) {
	u.id = id
}

// ID es un "getter" para leer el ID.
func (u *User) ID() UserID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Role() Role {
	return u.role
}

func (u *User) Active() bool {
	return u.active
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// ChangeRole permite cambiar el rol de forma controlada.
func (u *User) ChangeRole(newRole Role) error {
	if !isValidRole(newRole) {
		return errors.New("no se puede asignar un rol inválido")
	}
	u.role = newRole
	return nil
}

// Deactivate desactiva el usuario (no lo borra físicamente).
func (u *User) Deactivate() {
	u.active = false
}

/*
   ======================
   ESTRUCTURA: BOOK (LIBRO)
   ======================
*/

// Book representa un libro electrónico dentro del sistema.
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

// NewBook es el constructor para Book.
func NewBook(title, author string, year int, isbn, categoryTI string, tags []string) (*Book, error) {
	title = strings.TrimSpace(title)
	author = strings.TrimSpace(author)
	isbn = strings.TrimSpace(isbn)
	categoryTI = strings.TrimSpace(categoryTI)

	if title == "" {
		return nil, errors.New("el título no puede estar vacío")
	}
	if author == "" {
		return nil, errors.New("el autor no puede estar vacío")
	}
	if year <= 0 {
		return nil, errors.New("el año debe ser mayor que cero")
	}

	// Normalizamos tags a minúsculas y sin espacios.
	normalizedTags := make([]string, 0, len(tags))
	for _, t := range tags {
		t = strings.ToLower(strings.TrimSpace(t))
		if t != "" {
			normalizedTags = append(normalizedTags, t)
		}
	}

	return &Book{
		title:      title,
		author:     author,
		year:       year,
		isbn:       isbn,
		categoryTI: categoryTI,
		tags:       normalizedTags,
		active:     true,
		createdAt:  time.Now(),
	}, nil
}

// ====== Métodos de Book ======

func (b *Book) SetID(id BookID) {
	b.id = id
}

func (b *Book) ID() BookID {
	return b.id
}

func (b *Book) Title() string {
	return b.title
}

func (b *Book) Author() string {
	return b.author
}

func (b *Book) Year() int {
	return b.year
}

func (b *Book) ISBN() string {
	return b.isbn
}

func (b *Book) CategoryTI() string {
	return b.categoryTI
}

// Tags devuelve una COPIA del slice para proteger la estructura interna.
func (b *Book) Tags() []string {
	copia := make([]string, len(b.tags))
	copy(copia, b.tags)
	return copia
}

func (b *Book) Active() bool {
	return b.active
}

func (b *Book) CreatedAt() time.Time {
	return b.createdAt
}

// Archive "archiva" el libro (lo marca como inactivo).
func (b *Book) Archive() {
	b.active = false
}

/*
   ==========================
   ESTRUCTURA: AccessEvent (Registro de acceso)
   ==========================
*/

type AccessType string

const (
	AccessOpen     AccessType = "APERTURA"
	AccessRead     AccessType = "LECTURA"
	AccessDownload AccessType = "DESCARGA"
)

// AccessEvent representa un registro de acceso a un libro.
type AccessEvent struct {
	id        AccessEventID
	bookID    BookID
	userID    UserID
	access    AccessType
	timestamp time.Time
}

// NewAccessEvent crea un nuevo evento de acceso.
func NewAccessEvent(bookID BookID, userID UserID, access AccessType) (*AccessEvent, error) {
	if bookID <= 0 {
		return nil, errors.New("bookID inválido")
	}
	if userID <= 0 {
		return nil, errors.New("userID inválido")
	}
	return &AccessEvent{
		bookID:    bookID,
		userID:    userID,
		access:    access,
		timestamp: time.Now(),
	}, nil
}

// Métodos de AccessEvent.

func (e *AccessEvent) SetID(id AccessEventID) {
	e.id = id
}

func (e *AccessEvent) ID() AccessEventID {
	return e.id
}

func (e *AccessEvent) BookID() BookID {
	return e.bookID
}

func (e *AccessEvent) UserID() UserID {
	return e.userID
}

func (e *AccessEvent) Access() AccessType {
	return e.access
}

func (e *AccessEvent) Timestamp() time.Time {
	return e.timestamp
}

/*
   ==================
   FILTRO DE BÚSQUEDA
   ==================
*/

// BookFilter se usa para hacer búsquedas avanzadas de libros.
// Aquí usamos SLICES (Tags) para el requerimiento del ejercicio.
type BookFilter struct {
	TitleContains  string
	AuthorContains string
	CategoryTI     string
	Tags           []string
	YearFrom       int
	YearTo         int
}

/*
   =================
   INTERFACES (contratos)
   =================
*/

// UserRepository define qué operaciones debe soportar cualquier repositorio de usuarios.
type UserRepository interface {
	Create(user *User) error
	Update(user *User) error
	FindByID(id UserID) (*User, error)
	FindByEmail(email string) (*User, error)
	ListAll() ([]*User, error)
}

// BookRepository define operaciones para los libros.
type BookRepository interface {
	Create(book *Book) error
	Update(book *Book) error
	FindByID(id BookID) (*Book, error)
	SearchByFilters(filter BookFilter) ([]*Book, error)
	ListAll() ([]*Book, error)
}

// AccessLogRepository define operaciones para el registro de accesos.
type AccessLogRepository interface {
	Store(event *AccessEvent) error
	ListByBook(bookID BookID) ([]*AccessEvent, error)
	ListByUser(userID UserID) ([]*AccessEvent, error)
}
