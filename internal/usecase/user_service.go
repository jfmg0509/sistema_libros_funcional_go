package usecase

import (
	"errors"
	"strings"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

// UserService contiene la lógica de negocio para usuarios.
type UserService struct {
	repo domain.UserRepository
}

// NewUserService es una función CONSTRUCTOR de UserService.
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser registra un nuevo usuario en el sistema.
func (s *UserService) RegisterUser(name, email string, role domain.Role) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, errors.New("el email no puede estar vacío")
	}

	// Validamos si ya existe un usuario con ese email.
	if existing, _ := s.repo.FindByEmail(email); existing != nil {
		return nil, errors.New("ya existe un usuario con ese email")
	}

	// Creamos el usuario (sin ID, el ID lo asigna el repositorio).
	user, err := domain.NewUser(name, email, role)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangeUserRole cambia el rol de un usuario existente.
func (s *UserService) ChangeUserRole(userID domain.UserID, newRole domain.Role) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("usuario no encontrado")
	}

	if err := user.ChangeRole(newRole); err != nil {
		return err
	}

	return s.repo.Update(user)
}

// DeactivateUser desactiva un usuario.
func (s *UserService) DeactivateUser(userID domain.UserID) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("usuario no encontrado")
	}

	user.Deactivate()
	return s.repo.Update(user)
}

// ListUsers lista todos los usuarios.
func (s *UserService) ListUsers() ([]*domain.User, error) {
	return s.repo.ListAll()
}
