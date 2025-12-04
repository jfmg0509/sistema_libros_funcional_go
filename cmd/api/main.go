package main

import (
	"log"
	"net/http"

	dbInfra "github.com/jfmg0509/sistema_libros_funcional_go/internal/infrastructure/db"
	httpTransport "github.com/jfmg0509/sistema_libros_funcional_go/internal/transport/http"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

func main() {
	// 1) Creamos repositorios en memoria
	userRepo := dbInfra.NewInMemoryUserRepo()
	bookRepo := dbInfra.NewInMemoryBookRepo()
	accessRepo := dbInfra.NewInMemoryAccessLogRepo()

	// 2) Creamos servicios de negocio
	userService := usecase.NewUserService(userRepo)
	bookService := usecase.NewBookService(bookRepo, accessRepo)

	// 3) Creamos capa de transporte HTTP
	handler := httpTransport.NewHTTPHandler(userService, bookService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Println("Servidor HTTP iniciado en http://localhost:8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("error al iniciar servidor: %v", err)
	}
}
