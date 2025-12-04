package http

import (
	"encoding/json"
	nethttp "net/http"
	"strconv"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/domain"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

/*
HTTPHandler agrupa los servicios de la capa de negocio (usecase)
que vamos a exponer mediante HTTP (API REST).
*/
type HTTPHandler struct {
	userService *usecase.UserService
	bookService *usecase.BookService
}

// NewHTTPHandler crea un nuevo HTTPHandler listo para registrar rutas.
func NewHTTPHandler(userSvc *usecase.UserService, bookSvc *usecase.BookService) *HTTPHandler {
	return &HTTPHandler{
		userService: userSvc,
		bookService: bookSvc,
	}
}

/*
RegisterRoutes recibe un *http.ServeMux y registra todas
las rutas/endpoints de nuestro servicio.
*/
func (h *HTTPHandler) RegisterRoutes(mux *nethttp.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
	mux.HandleFunc("/users", h.handleUsers)
	mux.HandleFunc("/books", h.handleBooks)

	// NUEVAS rutas:
	mux.HandleFunc("/access", h.handleAccess)            // registrar acceso
	mux.HandleFunc("/access/stats", h.handleAccessStats) // estadísticas
}

/*
==========================
ENDPOINT /health
==========================
*/
func (h *HTTPHandler) handleHealth(w nethttp.ResponseWriter, r *nethttp.Request) {
	writeJSON(w, nethttp.StatusOK, map[string]string{
		"status": "ok",
	})
}

/*
==========================
ENDPOINT /users
==========================
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
			writeError(w, nethttp.StatusBadRequest, "JSON inválido en creación de usuario")
			return
		}

		user, err := h.userService.RegisterUser(payload.Name, payload.Email, payload.Role)
		if err != nil {
			writeError(w, nethttp.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, nethttp.StatusCreated, user)

	default:
		writeError(w, nethttp.StatusMethodNotAllowed, "método no permitido en /users")
	}
}

/*
==========================
ENDPOINT /books
==========================
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
==========================
ENDPOINT /access  (POST)
==========================
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
==========================
ENDPOINT /access/stats  (GET)
==========================
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
   FUNCIONES AUXILIARES JSON
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
