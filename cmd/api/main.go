package main

import "fmt"

// función simple para mostrar un mensaje
func mensajeBienvenida() string {
	return "Sistema de Gestión de Libros Electrónicos – Estructura inicial"
}

func main() {
	// uso de variables, función y salida por pantalla (Unidad 1)
	mensaje := mensajeBienvenida()
	fmt.Println(mensaje)
}

