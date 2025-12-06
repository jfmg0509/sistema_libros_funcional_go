package http

import (
	"encoding/json"
	nethttp "net/http"
	"strconv"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

/*
   ==========================================================
   HTTPHandler
   ==========================================================

   Esta estructura se encarga de conectar la capa HTTP (peticiones)
   con la capa de negocio (usecase).

   - No guarda datos.
   - Solo recibe requests, llama a servicios, y responde JSON.
*/

type HTTPHandler struct {
	userService *usecase.UserService
	bookService *usecase.BookService
}

// NewHTTPHandler es el CONSTRUCTOR del handler HTTP.
func NewHTTPHandler(userSvc *usecase.UserService, bookSvc *usecase.BookService) *HTTPHandler {
	return &HTTPHandler{
		userService: userSvc,
		bookService: bookSvc,
	}
}

/*
RegisterRoutes recibe un *http.ServeMux y registra todas
las rutas que soporta nuestra API.

Aquí definimos:
- /health
- /users
- /books
- /access
- /access/stats
*/
func (h *HTTPHandler) RegisterRoutes(mux *nethttp.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/books", h.handleBooks)
	mux.HandleFunc("/access", h.handleAccess)
	mux.HandleFunc("/access/stats", h.handleAccessStats)
}

/*
==========================================================
ENDPOINT /health
==========================================================

Método: GET

Sirve para verificar que el servidor está vivo.
Responde: {"status": "ok"}
*/
func (h *HTTPHandler) handleHealth(w nethttp.ResponseWriter, r *nethttp.Request) {
	writeJSON(w, nethttp.StatusOK, map[string]string{
		"status": "ok",
	})
}

/*
==========================================================
ENDPOINT /users
==========================================================

Métodos soportados:
- GET  /users  → Lista todos los usuarios
- POST /users  → Crea un nuevo usuario

Formato JSON para crear usuario:

	{
	  "name": "Marleen",
	  "email": "marleen@example.com",
	  "role": "ADMIN"
	}
*/
func (h *HTTPHandler) handleUsers(w nethttp.ResponseWriter, r *nethttp.Request) {
	switch r.Method {
	case nethttp.MethodGet:
		// Obtener lista de usuarios desde la capa de negocio.
		users, err := h.userService.ListUsers()
		if err != nil {
			writeError(w, nethttp.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, nethttp.StatusOK, users)

	case nethttp.MethodPost:
		// Estructura auxiliar para leer el JSON de entrada.
		var payload struct {
			Name  string      `json:"name"`
			Email string      `json:"email"`
			Role  domain.Role `json:"role"`
		}

		// Decodificar el JSON del body en la estructura payload.
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, nethttp.StatusBadRequest, "JSON inválido en creación de usuario")
			return
		}

		// Llamar al caso de uso para registrar el usuario.
		user, err := h.userService.RegisterUser(payload.Name, payload.Email, payload.Role)
		if err != nil {
			writeError(w, nethttp.StatusBadRequest, err.Error())
			return
		}

		// Responder con el usuario creado (aunque por encapsulación se ve como {} en JSON).
		writeJSON(w, nethttp.StatusCreated, user)

	default:
		writeError(w, nethttp.StatusMethodNotAllowed, "método no permitido en /users")
	}
}

/*
==========================================================
ENDPOINT /books
==========================================================

Métodos soportados:
- GET  /books         → Lista o busca libros por filtros.
- POST /books         → Crea un nuevo libro.

Ejemplo JSON para crear libro:

	{
	  "title": "Seguridad Informática",
	  "author": "Baca Urbina",
	  "year": 2016,
	  "isbn": "123-456",
	  "category_ti": "Seguridad",
	  "tags": ["seguridad","ciberseguridad"]
	}
*/
func (h *HTTPHandler) handleBooks(w nethttp.ResponseWriter, r *nethttp.Request) {
	switch r.Method {
	case nethttp.MethodGet:
		// Leer parámetros de consulta (query params).
		query := r.URL.Query()
		filter := domain.BookFilter{
			TitleContains:  query.Get("title"),
			AuthorContains: query.Get("author"),
			CategoryTI:     query.Get("category"),
		}

		// Convertir year_from y year_to si vienen en la URL.
		if yearFromStr := query.Get("year_from"); yearFromStr != "" {
			if yearFrom, err := strconv.Atoi(yearFromStr); err == nil {
				filter.YearFrom = yearFrom
			}
		}
		if yearToStr := query.Get("year_to"); yearToStr != "" {
			if yearTo, err := strconv.Atoi(yearToStr); err == nil {
				filter.YearTo = yearTo
			}
		}

		books, err := h.bookService.SearchBooks(filter)
		if err != nil {
			writeError(w, nethttp.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, nethttp.StatusOK, books)

	case nethttp.MethodPost:
		// Estructura auxiliar para el JSON de entrada.
		var payload struct {
			Title      string   `json:"title"`
			Author     string   `json:"author"`
			Year       int      `json:"year"`
			ISBN       string   `json:"isbn"`
			CategoryTI string   `json:"category_ti"`
			Tags       []string `json:"tags"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, nethttp.StatusBadRequest, "JSON inválido en creación de libro")
			return
		}

		book, err := h.bookService.RegisterBook(
			payload.Title,
			payload.Author,
			payload.Year,
			payload.ISBN,
			payload.CategoryTI,
			payload.Tags,
		)
		if err != nil {
			writeError(w, nethttp.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, nethttp.StatusCreated, book)

	default:
		writeError(w, nethttp.StatusMethodNotAllowed, "método no permitido en /books")
	}
}

/*
==========================================================
ENDPOINT /access
==========================================================

Método soportado:
- POST /access   → registra un acceso de un usuario a un libro.

Ejemplo JSON:

	{
	  "book_id": 1,
	  "user_id": 1,
	  "access_type": "LECTURA"
	}
*/
func (h *HTTPHandler) handleAccess(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Method != nethttp.MethodPost {
		writeError(w, nethttp.StatusMethodNotAllowed, "método no permitido en /access")
		return
	}

	var payload struct {
		BookID     int64             `json:"book_id"`
		UserID     int64             `json:"user_id"`
		AccessType domain.AccessType `json:"access_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, nethttp.StatusBadRequest, "JSON inválido en registro de acceso")
		return
	}

	if payload.BookID <= 0 {
		writeError(w, nethttp.StatusBadRequest, "book_id debe ser mayor que cero")
		return
	}
	if payload.UserID <= 0 {
		writeError(w, nethttp.StatusBadRequest, "user_id debe ser mayor que cero")
		return
	}

	// Llamar a la lógica de negocio para registrar el acceso.
	err := h.bookService.RecordAccess(
		domain.BookID(payload.BookID),
		domain.UserID(payload.UserID),
		payload.AccessType,
	)
	if err != nil {
		writeError(w, nethttp.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, nethttp.StatusCreated, map[string]string{
		"message": "acceso registrado correctamente",
	})
}

/*
==========================================================
ENDPOINT /access/stats
==========================================================

Método soportado:
- GET /access/stats?book_id=1  → devuelve estadísticas de accesos

Ejemplo de respuesta:

	{
	  "LECTURA":  3,
	  "APERTURA": 5,
	  "DESCARGA": 1
	}
*/
func (h *HTTPHandler) handleAccessStats(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Method != nethttp.MethodGet {
		writeError(w, nethttp.StatusMethodNotAllowed, "método no permitido en /access/stats")
		return
	}

	query := r.URL.Query()
	bookIDStr := query.Get("book_id")
	if bookIDStr == "" {
		writeError(w, nethttp.StatusBadRequest, "parámetro book_id es obligatorio")
		return
	}

	bookIDInt, err := strconv.ParseInt(bookIDStr, 10, 64)
	if err != nil || bookIDInt <= 0 {
		writeError(w, nethttp.StatusBadRequest, "parámetro book_id debe ser un número válido mayor que cero")
		return
	}

	stats, err := h.bookService.BuildAccessStatsByBook(domain.BookID(bookIDInt))
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, nethttp.StatusOK, stats)
}

/*
   ==========================================================
   Funciones auxiliares para respuestas JSON
   ==========================================================
*/

// writeJSON escribe una respuesta JSON con el código de estado indicado.
func writeJSON(w nethttp.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// writeError simplifica el envío de errores en formato JSON.
func writeError(w nethttp.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}
