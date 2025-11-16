# 📚 Sistema de Gestión de Libros Electrónicos (Go – Programación Funcional)

Este repositorio contiene el desarrollo del **Sistema de Gestión de Libros Electrónicos** implementado en **Go (Golang)**, siguiendo un enfoque de **programación funcional**

---

## 🎯 Objetivo del sistema
Desarrollar una plataforma interna que permita gestionar libros técnicos digitales utilizados en consultoría TI, facilitando la búsqueda, clasificación y acceso a la información utilizando un diseño modular y funcional.

---

## 🧱 Tecnologías y paradigma
- Lenguaje: **Go 1.20+**
- Paradigma: **Programación funcional**
  - funciones puras  
  - closures  
  - composición  
  - evitar estados mutables  
- Dependencias externas: **Ninguna**

---

## 📁 Estructura del repositorio (propuesta inicial)

```text
cmd/
  api/
    main.go

internal/
  domain/
  usecase/
  infrastructure/
    db/
  config/
  transport/
    http/

docs/
go run ./cmd/api
Autor
Juan Francisco Morán

Licencia
pendiente de selección
