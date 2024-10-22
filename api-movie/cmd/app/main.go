package main

import (
	handler "api-movie/cmd/app/handlers"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

const apiKey = "98cf6e1c64f9167cc7ee0ea8b6b4779b" // Substitua pela sua chave de API do TMDb
const tmdbURL = "https://api.themoviedb.org/3/movie/popular?api_key=" + apiKey

var db *sql.DB

func main() {
	// Carregar variáveis de ambiente
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo .env: %v", err)
	}

	// Obter configurações do banco de dados
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// String de conexão
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Conectar ao banco de dados
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer db.Close()

	// Testar a conexão
	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao testar conexão com o banco de dados: %v", err)
	}
	log.Println("Conectado ao banco de dados MySQL com sucesso!")

	// Configurar o banco de dados no handler
	handler.SetDB(db)

	// Configurar rotas
	http.HandleFunc("/movies", handler.MoviesHandler)
	http.HandleFunc("/movies/", handler.MovieHandler) // Para rotas com ID

	// Obter a porta
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Aplicação rodando na porta %s", port)

	// Iniciar o servidor
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
