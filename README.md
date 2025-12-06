# Sistema de Gestión de Libros Electrónicos en Go

Proyecto desarrollado para la carrera de **Ingeniería en Sistemas**, cuyo objetivo es implementar un **Sistema de Gestión de Libros Electrónicos** utilizando el lenguaje **Go (Golang)**, aplicando arquitectura por capas, encapsulación, interfaces, manejo de errores y estructuras de datos (arrays, slices y maps).

---

## 🎯 Objetivo del Sistema

Este sistema permite:

- Registrar **usuarios** con roles (por ejemplo: `ADMIN`, `READER`).
- Registrar **libros electrónicos** con:
  - Título
  - Autor
  - Año
  - ISBN
  - Categoría TI
  - Tags
- Registrar **accesos de usuarios a libros**:
  - APERTURA
  - LECTURA
  - DESCARGA
- Generar **estadísticas de accesos por libro** usando mapas (`map[AccessType]int`).

Todo se expone mediante una **API REST**.

---

## 🧱 Arquitectura por Capas

El proyecto está organizado en los siguientes paquetes:

### 1. `internal/domain`
Define el **modelo de dominio**:

- Entidades:
  - `User`
  - `Book`
  - `AccessEvent`
- Tipos:
  - `UserID`, `BookID`, `AccessEventID`
  - `Role` (`ADMIN`, `READER`)
  - `AccessType` (`APERTURA`, `LECTURA`, `DESCARGA`)
- Filtro de libros:
  - `BookFilter`
- Interfaces:
  - `UserRepository`
  - `BookRepository`
  - `AccessLogRepository`

También implementa **encapsulación** mediante campos privados y métodos públicos (`ID()`, `Name()`, `Email()`, etc.)

---

### 2. `internal/usecase`

Contiene la **lógica de negocio** (casos de uso):

- `UserService`
  - `RegisterUser(name, email, role)`
  - `ListUsers()`
- `BookService`
  - `RegisterBook(...)`
  - `SearchBooks(filter)`
  - `RecordAccess(bookID, userID, accessType)`
  - `BuildAccessStatsByBook(bookID)`

Aquí se aplican reglas como:
- Validar que no exista un usuario con el mismo email.
- Verificar que el usuario y el libro existan antes de registrar un acceso.
- Construir estadísticas usando `map[AccessType]int`.

---

### 3. `internal/infrastructure/db`

Implementación de repositorios **en memoria** usando `map`:

- `InMemoryUserRepo`:
  - `users: map[UserID]*User`
  - `emailIndex: map[string]UserID`
- `InMemoryBookRepo`:
  - `books: map[BookID]*Book`
- `InMemoryAccessLogRepo`:
  - `events: map[AccessEventID]*AccessEvent`

Esta capa simula una base de datos y es ideal para prácticas y prototipos.

---

### 4. `internal/transport/http`

Expone los servicios como una **API REST** usando `net/http`:

Rutas principales:

- `GET    /health`
- `GET    /users`
- `POST   /users`
- `GET    /books`
- `POST   /books`
- `POST   /access`
- `GET    /access/stats?book_id={id}`

Cada handler:
- Lee parámetros o JSON de entrada.
- Llama a la capa de negocio (`usecase`).
- Devuelve respuestas JSON (datos o errores).

---

### 5. `cmd/api/main.go`

Punto de entrada de la aplicación:

1. Crea los repositorios en memoria.
2. Crea los servicios (`UserService`, `BookService`).
3. Crea el `HTTPHandler` y registra las rutas.
4. Levanta el servidor HTTP:

```go
log.Println("Servidor HTTP iniciado en http://localhost:8081")
http.ListenAndServe(":8081", mux)
