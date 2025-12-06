package usecase

import (
	"fmt"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
)

/*
   ==========================================================
   UserService
   ==========================================================

   Esta estructura representa la "capa de negocio" para usuarios.
   No sabe cómo se guardan los datos (eso lo hace el repositorio).
   Solo sabe QUÉ reglas aplicar al registrar o listar usuarios.
*/

// UserService contiene un repositorio que cumple la interfaz UserRepository.
type UserService struct {
	repo domain.UserRepository
}

// NewUserService es el CONSTRUCTOR del servicio de usuarios.
// Recibe un objeto que implemente domain.UserRepository (por ejemplo, el repositorio en memoria).
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

/*
RegisterUser registra un nuevo usuario en el sistema.

Pasos:
1. Verifica si ya existe un usuario con el mismo email.
2. Si existe, devuelve un error controlado.
3. Si no existe, usa el CONSTRUCTOR de dominio (NewUser) para crear el usuario.
4. Pide al repositorio que lo guarde.
5. Devuelve el usuario creado.
*/
func (s *UserService) RegisterUser(name, email string, role domain.Role) (*domain.User, error) {
	// 1. Verificar si ya existe un usuario con ese email.
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// Mensaje que viste cuando probaste con curl.
		return nil, fmt.Errorf("ya existe un usuario con ese email")
	}

	// 2. Crear el usuario usando el constructor del dominio.
	user, err := domain.NewUser(name, email, role)
	if err != nil {
		return nil, err
	}

	// 3. Guardar el usuario en el repositorio.
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 4. Devolver el usuario creado.
	return user, nil
}

/*
ListUsers devuelve todos los usuarios registrados.

La lógica es sencilla:
- Llamar al repositorio.
- Devolver la lista.
*/
func (s *UserService) ListUsers() ([]*domain.User, error) {
	return s.repo.ListAll()
}

/*
FindUserByID busca un usuario por su ID.

Esto no lo usamos todavía en los handlers HTTP,
pero es útil si más adelante quieres agregar endpoints
como GET /users/{id}.
*/
func (s *UserService) FindUserByID(id domain.UserID) (*domain.User, error) {
	return s.repo.FindByID(id)
}
