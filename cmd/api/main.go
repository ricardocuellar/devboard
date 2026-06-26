// Package main
package main

import (
	"log"

	"github.com/ricardocuellar/devboard/internal/server"
)

func main() {
	// Aqui arranca el servidor
	srv := server.New(":8080")

	log.Println("Servidor iniciando en :8080")

	if err := srv.Start(); err != nil {
		log.Fatalf("error al iniciar el servidor %v", err)
	}

}
