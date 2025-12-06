package main

import (
	"log"
	nethttp "net/http"

	"github.com/jfmg0509/sistema_libros_funcional_go/internal/infrastructure/db"
	httptransport "github.com/jfmg0509/sistema_libros_funcional_go/internal/transport/http"
	"github.com/jfmg0509/sistema_libros_funcional_go/internal/usecase"
)

/*
   ==========================================================
   main.go
   ==========================================================

   Este archivo es el PUNTO DE ENTRADA del sistema.

   Pasos:
   1. Crear repositorios en memoria.
   2. Crear servicios (capa de negocio/usecase).
   3. Crear el handler HTTP.
   4. Registrar rutas en un ServeMux.
   5. Iniciar el servidor HTTP en el puerto 8081.
*/

func main() {
	// 1. Crear repositorios (implementan interfaces del dominio).
	userRepo := db.NewInMemoryUserRepo()
	bookRepo := db.NewInMemoryBookRepo()
	accessRepo := db.NewInMemoryAccessLogRepo()

	// 2. Crear servicios de negocio, inyectando los repositorios.
	userService := usecase.NewUserService(userRepo)
	bookService := usecase.NewBookService(bookRepo, userRepo, accessRepo)

	// 3. Crear el handler HTTP, que usará los servicios.
	handler := httptransport.NewHTTPHandler(userService, bookService)

	// 4. Crear un enrutador (ServeMux) y registrar las rutas.
	mux := nethttp.NewServeMux()
	handler.RegisterRoutes(mux)

	// 5. Levantar el servidor HTTP.
	log.Println("Servidor HTTP iniciado en http://localhost:8081")
	if err := nethttp.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("error al iniciar servidor: %v", err)
	}
}
