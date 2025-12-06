package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ------------------------------------------------------------
// TIPOS BÁSICOS DEL SISTEMA
// ------------------------------------------------------------

// Role representa el rol de un usuario (ADMIN o READER).
type Role string

const (
	RoleAdmin  Role = "ADMIN"
	RoleReader Role = "READER"
)

// AccessType representa el tipo de acceso a un libro.
type AccessType string

const (
	AccessApertura AccessType = "APERTURA"
	AccessLectura  AccessType = "LECTURA"
	AccessDescarga AccessType = "DESCARGA"
)

// User representa un usuario del sistema.
type User struct {
	ID    int
	Name  string
	Email string
	Role  Role
}

// Book representa un libro electrónico (sin tags, como pediste).
type Book struct {
	ID       int
	Title    string
	Author   string
	Year     int
	ISBN     string
	Category string
}

// AccessEvent representa un evento de acceso de un usuario a un libro.
type AccessEvent struct {
	ID     int
	BookID int
	UserID int
	Type   AccessType
}

// ------------------------------------------------------------
// "BASE DE DATOS" EN MEMORIA (slices)
// ------------------------------------------------------------

var users []User
var books []Book
var accessEvents []AccessEvent

// ------------------------------------------------------------
// FUNCIÓN PRINCIPAL: MENÚ INTERACTIVO
// ------------------------------------------------------------

func main() {
	// Scanner para leer desde la terminal (entrada estándar).
	scanner := bufio.NewScanner(os.Stdin)

	for {
		// Mostrar menú principal
		fmt.Println("==========================================")
		fmt.Println(" SISTEMA DE GESTIÓN DE LIBROS ELECTRÓNICOS")
		fmt.Println("           (MODO TERMINAL / CLI)")
		fmt.Println("==========================================")
		fmt.Println("1. Registrar usuario")
		fmt.Println("2. Listar usuarios")
		fmt.Println("3. Registrar libro")
		fmt.Println("4. Listar libros")
		fmt.Println("5. Buscar libros por título/autor")
		fmt.Println("6. Registrar acceso a un libro")
		fmt.Println("7. Ver estadísticas de accesos de un libro")
		fmt.Println("0. Salir")
		fmt.Print("Selecciona una opción: ")

		// Leer opción ingresada por el usuario
		if !scanner.Scan() {
			fmt.Println("No se pudo leer la opción. Saliendo...")
			return
		}
		opcion := strings.TrimSpace(scanner.Text())

		fmt.Println() // línea en blanco para separar

		switch opcion {
		case "1":
			registerUser(scanner)
		case "2":
			listUsers(scanner) // ahora recibe scanner para poder pausar
		case "3":
			registerBook(scanner)
		case "4":
			listBooks()
		case "5":
			searchBooks(scanner)
		case "6":
			registerAccess(scanner)
		case "7":
			showAccessStats(scanner)
		case "0":
			fmt.Println("Saliendo del sistema... ¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida. Intenta nuevamente.")
		}

		fmt.Println() // otra línea en blanco para separar ciclos
	}
}

// ------------------------------------------------------------
// USUARIOS
// ------------------------------------------------------------

// registerUser pide los datos por consola y crea un nuevo usuario.
func registerUser(scanner *bufio.Scanner) {
	fmt.Println("=== Registrar nuevo usuario ===")

	// ID se genera automáticamente como siguiente número.
	newID := len(users) + 1

	fmt.Print("Nombre: ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el nombre.")
		return
	}
	name := strings.TrimSpace(scanner.Text())

	fmt.Print("Email: ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el email.")
		return
	}
	email := strings.TrimSpace(scanner.Text())

	// 🔒 Validación básica: el email debe contener '@'.
	if !strings.Contains(email, "@") {
		fmt.Println("Email inválido. Debe contener '@'.")
		return
	}

	fmt.Print("Rol (ADMIN/READER): ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el rol.")
		return
	}
	roleInput := strings.ToUpper(strings.TrimSpace(scanner.Text()))

	var role Role
	switch roleInput {
	case "ADMIN":
		role = RoleAdmin
	case "READER":
		role = RoleReader
	default:
		fmt.Println("Rol inválido. Debe ser ADMIN o READER.")
		return
	}

	// 🔒 Validación: no permitir usuarios repetidos por nombre o email.
	for _, u := range users {
		if strings.EqualFold(u.Email, email) {
			fmt.Println("Ya existe un usuario con ese email.")
			return
		}
		if strings.EqualFold(u.Name, name) {
			fmt.Println("Ya existe un usuario con ese nombre.")
			return
		}
	}

	// Crear el usuario y agregarlo al slice en memoria.
	user := User{
		ID:    newID,
		Name:  name,
		Email: email,
		Role:  role,
	}

	users = append(users, user)
	fmt.Println("Usuario registrado correctamente con ID:", user.ID)
}

// listUsers imprime todos los usuarios y espera ENTER para volver al menú.
func listUsers(scanner *bufio.Scanner) {
	fmt.Println("=== Listado de usuarios ===")

	if len(users) == 0 {
		fmt.Println("No hay usuarios registrados.")
	} else {
		for _, u := range users {
			fmt.Printf("ID: %d | Nombre: %s | Email: %s | Rol: %s\n",
				u.ID, u.Name, u.Email, u.Role)
		}
	}

	// Pausa para que el usuario pueda leer antes de volver al menú.
	fmt.Print("\nPresiona ENTER para volver al menú...")
	scanner.Scan()
}

// ------------------------------------------------------------
// LIBROS
// ------------------------------------------------------------

// registerBook pide los datos y registra un libro nuevo.
func registerBook(scanner *bufio.Scanner) {
	fmt.Println("=== Registrar nuevo libro ===")

	newID := len(books) + 1

	fmt.Print("Título: ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el título.")
		return
	}
	title := strings.TrimSpace(scanner.Text())

	fmt.Print("Autor: ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el autor.")
		return
	}
	author := strings.TrimSpace(scanner.Text())

	fmt.Print("Año de publicación (número): ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el año.")
		return
	}
	yearInput := strings.TrimSpace(scanner.Text())
	year, err := strconv.Atoi(yearInput)
	if err != nil {
		fmt.Println("Año inválido. Debe ser un número.")
		return
	}

	fmt.Print("ISBN: ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el ISBN.")
		return
	}
	isbn := strings.TrimSpace(scanner.Text())

	fmt.Print("Categoría TI (por ejemplo: SEGURIDAD, REDES, BD, ETC.): ")
	if !scanner.Scan() {
		fmt.Println("Error al leer la categoría.")
		return
	}
	category := strings.TrimSpace(scanner.Text())

	// Crear el libro (sin tags, como acordamos).
	book := Book{
		ID:       newID,
		Title:    title,
		Author:   author,
		Year:     year,
		ISBN:     isbn,
		Category: category,
	}

	books = append(books, book)
	fmt.Println("Libro registrado correctamente con ID:", book.ID)
}

// listBooks imprime todos los libros registrados.
func listBooks() {
	fmt.Println("=== Listado de libros ===")

	if len(books) == 0 {
		fmt.Println("No hay libros registrados.")
		return
	}

	for _, b := range books {
		fmt.Printf("ID: %d | Título: %s | Autor: %s | Año: %d | ISBN: %s | Categoría: %s\n",
			b.ID, b.Title, b.Author, b.Year, b.ISBN, b.Category)
	}
}

// searchBooks permite buscar libros por título o autor (búsqueda parcial).
func searchBooks(scanner *bufio.Scanner) {
	fmt.Println("=== Buscar libros ===")

	fmt.Print("Texto a buscar (en título o autor): ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el texto de búsqueda.")
		return
	}
	query := strings.ToLower(strings.TrimSpace(scanner.Text()))

	if query == "" {
		fmt.Println("La búsqueda no puede estar vacía.")
		return
	}

	encontrados := 0
	for _, b := range books {
		titleLower := strings.ToLower(b.Title)
		authorLower := strings.ToLower(b.Author)

		if strings.Contains(titleLower, query) || strings.Contains(authorLower, query) {
			encontrados++
			fmt.Printf("ID: %d | Título: %s | Autor: %s | Año: %d\n",
				b.ID, b.Title, b.Author, b.Year)
		}
	}

	if encontrados == 0 {
		fmt.Println("No se encontraron libros que coincidan con la búsqueda.")
	}
}

// ------------------------------------------------------------
// ACCESOS A LIBROS
// ------------------------------------------------------------

// registerAccess registra un evento de acceso (APERTURA, LECTURA, DESCARGA).
func registerAccess(scanner *bufio.Scanner) {
	fmt.Println("=== Registrar acceso a libro ===")

	// Pedir ID de usuario
	fmt.Print("ID de usuario: ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el ID de usuario.")
		return
	}
	userIDInput := strings.TrimSpace(scanner.Text())
	userID, err := strconv.Atoi(userIDInput)
	if err != nil {
		fmt.Println("ID de usuario inválido.")
		return
	}

	user := findUserByID(userID)
	if user == nil {
		fmt.Println("No existe un usuario con ese ID.")
		return
	}

	// Pedir ID de libro
	fmt.Print("ID de libro: ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el ID de libro.")
		return
	}
	bookIDInput := strings.TrimSpace(scanner.Text())
	bookID, err := strconv.Atoi(bookIDInput)
	if err != nil {
		fmt.Println("ID de libro inválido.")
		return
	}

	book := findBookByID(bookID)
	if book == nil {
		fmt.Println("No existe un libro con ese ID.")
		return
	}

	// Tipo de acceso
	fmt.Print("Tipo de acceso (APERTURA/LECTURA/DESCARGA): ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el tipo de acceso.")
		return
	}
	accessInput := strings.ToUpper(strings.TrimSpace(scanner.Text()))

	var accessType AccessType
	switch accessInput {
	case "APERTURA":
		accessType = AccessApertura
	case "LECTURA":
		accessType = AccessLectura
	case "DESCARGA":
		accessType = AccessDescarga
	default:
		fmt.Println("Tipo de acceso inválido.")
		return
	}

	// ID del evento de acceso se genera automáticamente.
	newID := len(accessEvents) + 1

	event := AccessEvent{
		ID:     newID,
		BookID: book.ID,
		UserID: user.ID,
		Type:   accessType,
	}

	accessEvents = append(accessEvents, event)

	fmt.Printf("Acceso registrado correctamente (Usuario %s -> Libro %s, Tipo: %s)\n",
		user.Name, book.Title, accessType)
}

// showAccessStats muestra cuántos accesos tiene un libro por tipo.
func showAccessStats(scanner *bufio.Scanner) {
	fmt.Println("=== Estadísticas de accesos por libro ===")

	fmt.Print("ID de libro: ")
	if !scanner.Scan() {
		fmt.Println("Error al leer el ID de libro.")
		return
	}
	bookIDInput := strings.TrimSpace(scanner.Text())
	bookID, err := strconv.Atoi(bookIDInput)
	if err != nil {
		fmt.Println("ID de libro inválido.")
		return
	}

	book := findBookByID(bookID)
	if book == nil {
		fmt.Println("No existe un libro con ese ID.")
		return
	}

	// Construimos un mapa AccessType -> cantidad de accesos.
	stats := make(map[AccessType]int)

	for _, e := range accessEvents {
		if e.BookID == book.ID {
			stats[e.Type]++
		}
	}

	fmt.Printf("Estadísticas de accesos para el libro: %s (ID %d)\n", book.Title, book.ID)
	if len(stats) == 0 {
		fmt.Println("Este libro no tiene accesos registrados.")
		return
	}

	fmt.Printf("APERTURA: %d\n", stats[AccessApertura])
	fmt.Printf("LECTURA : %d\n", stats[AccessLectura])
	fmt.Printf("DESCARGA: %d\n", stats[AccessDescarga])
}

// ------------------------------------------------------------
// FUNCIONES AUXILIARES PARA BUSCAR POR ID
// ------------------------------------------------------------

// findUserByID devuelve un puntero al usuario con ese ID o nil si no existe.
func findUserByID(id int) *User {
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

// findBookByID devuelve un puntero al libro con ese ID o nil si no existe.
func findBookByID(id int) *Book {
	for i := range books {
		if books[i].ID == id {
			return &books[i]
		}
	}
	return nil
}
