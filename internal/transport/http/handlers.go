package http

import (
	"encoding/json"
	nethttp "net/http"
	"strconv"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

// HTTPHandler agrupa los servicios que vamos a exponer por HTTP.
type HTTPHandler struct {
	userService *usecase.UserService
	bookService *usecase.BookService
}

// NewHTTPHandler crea un nuevo handler HTTP.
func NewHTTPHandler(userSvc *usecase.UserService, bookSvc *usecase.BookService) *HTTPHandler {
	return &HTTPHandler{
		userService: userSvc,
		bookService: bookSvc,
	}
}

// RegisterRoutes registra las rutas del servidor HTTP.
func (h *HTTPHandler) RegisterRoutes(mux *nethttp.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/books", h.handleBooks)
	mux.HandleFunc("/books/search", h.handleBooksSearch)
}

/*
HANDLER: /health
Método: GET
Descripción: Endpoint simple para verificar si el servidor está vivo.
*/
func (h *HTTPHandler) handleHealth(w nethttp.ResponseWriter, r *nethttp.Request) {
	writeJSON(w, nethttp.StatusOK, map[string]string{
		"status": "ok",
	})
}

/*
HANDLER: /users
Métodos:
  - GET  → lista usuarios
  - POST → crea usuario
*/
func (h *HTTPHandler) handleUsers(w nethttp.ResponseWriter, r *nethttp.Request) {
	switch r.Method {
	case nethttp.MethodGet:
		users, err := h.userService.ListUsers()
		if err != nil {
			writeError(w, nethttp.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, nethttp.StatusOK, users)

	case nethttp.MethodPost:
		var payload struct {
			Name  string      `json:"name"`
			Email string      `json:"email"`
			Role  domain.Role `json:"role"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, nethttp.StatusBadRequest, "JSON inválido")
			return
		}

		user, err := h.userService.RegisterUser(payload.Name, payload.Email, payload.Role)
		if err != nil {
			writeError(w, nethttp.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, nethttp.StatusCreated, user)

	default:
		writeError(w, nethttp.StatusMethodNotAllowed, "método no permitido")
	}
}

/*
HANDLER: /books
Métodos:
  - GET  → búsqueda simple por query params
  - POST → registra un nuevo libro
*/
func (h *HTTPHandler) handleBooks(w nethttp.ResponseWriter, r *nethttp.Request) {
	switch r.Method {
	case nethttp.MethodGet:
		query := r.URL.Query()
		filter := domain.BookFilter{
			TitleContains:  query.Get("title"),
			AuthorContains: query.Get("author"),
			CategoryTI:     query.Get("category"),
		}

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
		var payload struct {
			Title      string   `json:"title"`
			Author     string   `json:"author"`
			Year       int      `json:"year"`
			ISBN       string   `json:"isbn"`
			CategoryTI string   `json:"category_ti"`
			Tags       []string `json:"tags"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeError(w, nethttp.StatusBadRequest, "JSON inválido")
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
		writeError(w, nethttp.StatusMethodNotAllowed, "método no permitido")
	}
}

/*
HANDLER: /books/search
Método: GET
Descripción: ejemplo extra de búsqueda (puedes ampliarlo luego).
*/
func (h *HTTPHandler) handleBooksSearch(w nethttp.ResponseWriter, r *nethttp.Request) {
	if r.Method != nethttp.MethodGet {
		writeError(w, nethttp.StatusMethodNotAllowed, "método no permitido")
		return
	}

	query := r.URL.Query()
	filter := domain.BookFilter{
		TitleContains:  query.Get("title"),
		AuthorContains: query.Get("author"),
		CategoryTI:     query.Get("category"),
	}

	books, err := h.bookService.SearchBooks(filter)
	if err != nil {
		writeError(w, nethttp.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, nethttp.StatusOK, books)
}

/*
   FUNCIONES DE APOYO PARA RESPUESTAS JSON
*/

func writeJSON(w nethttp.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w nethttp.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}
