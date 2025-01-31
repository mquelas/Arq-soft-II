package main

import (
	"fmt"
	"log"
	"os"
	"search-api/initializers"
	"search-api/routes"
)

func init() {
	initializers.LoadEnv() // Cargar variables de entorno
}

func main() {
	// Si no hay variable de entorno, usar 8081 por defecto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Configurar las rutas
	r := routes.SetupRouter()

	// Iniciar el servidor
	log.Printf("Starting server on port %s...", port)
	r.Run(fmt.Sprintf(":%s", port)) // Inicia el servidor en el puerto
}
