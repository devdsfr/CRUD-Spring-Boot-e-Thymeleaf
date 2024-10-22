package main

import (
	handler "api-movie/cmd/app/handlers"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const apiKey = "98cf6e1c64f9167cc7ee0ea8b6b4779b" // Substitua pela sua chave de API do TMDb
const tmdbURL = "https://api.themoviedb.org/3/movie/popular?api_key=" + apiKey

func main() {
	// Carregar variáveis de ambiente
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Definir handlers
	http.HandleFunc("/movies", handler.MoviesHandler)

	log.Printf("Aplicação rodando na porta %s", port)

	// Iniciar o servidor e bloquear
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
