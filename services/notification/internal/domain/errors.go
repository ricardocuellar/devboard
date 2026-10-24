package domain

import "errors"

var (
	//ErrNotFound reutilizable
	ErrNotFound = errors.New("recurso no encontrado")
	//ErrAlreadyExists reutilizable
	ErrAlreadyExists = errors.New("el recurso ya existe")
	//ErrUnauthorized reutilizable
	ErrUnauthorized = errors.New("no autorizado")
	//ErrForbidden reutilizable
	ErrForbidden = errors.New("acceso prohibido")
	//ErrInvalidInput reutilizable
	ErrInvalidInput = errors.New("datos de entrada inválidos")
	//ErrInternal reutilizable
	ErrInternal = errors.New("error interno")
	//ErrDependencyUnavailable reutilizable
	ErrDependencyUnavailable = errors.New("Dependencia no disponible")
	//ErrDependencyTimeout reutilizable
	ErrDependencyTimeout = errors.New("Tiempo excedido")
)
