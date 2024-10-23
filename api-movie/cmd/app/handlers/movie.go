package handlers

import (
	service "api-movie/internal/services"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const apiKey = "98cf6e1c64f9167cc7ee0ea8b6b4779b" // Substitua pela sua chave de API do TMDb

// Movie representa um filme
type Movie struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Year             int     `json:"year"`
	Adult            bool    `json:"adult"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
}

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

// Handler para /movies (GET e POST)
func MoviesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMovies(w, r)
	case http.MethodPost:
		createMovie(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método não suportado"})
	}
}

// Handler para /movies/{id} (GET, PUT e DELETE)
func MovieHandler(w http.ResponseWriter, r *http.Request) {
	// Extrair ID da URL
	id, err := extractID(r.URL.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID inválido"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		getMovieByID(w, r, id)
	case http.MethodPut:
		updateMovie(w, r, id)
	case http.MethodDelete:
		deleteMovie(w, r, id)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método não suportado"})
	}
}

// Extrai o ID da URL /movies/{id}
func extractID(path string) (int, error) {
	var id int
	n, err := fmt.Sscanf(path, "/movies/%d", &id)
	if err != nil || n != 1 {
		return 0, fmt.Errorf("ID inválido")
	}
	return id, nil
}

// GET /movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, year FROM movies")
	if err != nil {
		http.Error(w, "Erro ao buscar filmes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Year)
		if err != nil {
			http.Error(w, "Erro ao escanear dados", http.StatusInternalServerError)
			return
		}
		movies = append(movies, movie)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// GET /movies/{id}
func getMovieByID(w http.ResponseWriter, r *http.Request, id int) {
	var movie Movie
	err := db.QueryRow("SELECT id, title, year FROM movies WHERE id = ?", id).Scan(&movie.ID, &movie.Title, &movie.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Filme não encontrado"})
			return
		}
		http.Error(w, "Erro ao buscar o filme", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}

// POST /movies
func createMovie(w http.ResponseWriter, r *http.Request) {
	httpClient := &http.Client{}

	// Cria uma nova instância do cliente TMDB
	tmdbClient := service.NewTMDBClient(apiKey, httpClient)
	movieService := service.NewMovieService(tmdbClient)

	// Imprime os detalhes dos filmes
	//movieService.GetMovieDetails("Inception")

	var newMovie Movie
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dados inválidos"})
		return
	}

	// Buscar detalhes do filme na API do TMDB
	movieDetails, err := movieService.GetMovieDetails(newMovie.Title)
	if err != nil {
		http.Error(w, "Erro ao buscar detalhes do filme na API", http.StatusInternalServerError)
		return
	}

	// Preencher os novos campos com os dados da API
	newMovie.Adult = movieDetails.Adult
	newMovie.OriginalLanguage = movieDetails.OriginalLanguage
	newMovie.Overview = movieDetails.Overview
	newMovie.Popularity = movieDetails.Popularity

	// Inserir no banco de dados com os novos campos
	result, err := db.Exec(
		"INSERT INTO movies ("+
			"title, year, adult, original_language, overview, popularity) VALUES (?, ?, ?, ?, ?, ?)",
		newMovie.Title,
		newMovie.Year,
		newMovie.Adult,
		newMovie.OriginalLanguage,
		newMovie.Overview,
		newMovie.Popularity)
	if err != nil {
		http.Error(w, "Erro ao criar o filme", http.StatusInternalServerError)
		return
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Erro ao obter ID do filme", http.StatusInternalServerError)
		return
	}

	newMovie.ID = strconv.FormatInt(lastInsertID, 10)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMovie)
}

// PUT /movies/{id}
func updateMovie(w http.ResponseWriter, r *http.Request, id int) {
	var updatedMovie Movie
	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dados inválidos"})
		return
	}

	// Atualizar no banco de dados
	result, err := db.Exec("UPDATE movies SET title = ?, year = ? WHERE id = ?", updatedMovie.Title, updatedMovie.Year, id)
	if err != nil {
		http.Error(w, "Erro ao atualizar o filme", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Erro ao verificar atualização", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Filme não encontrado"})
		return
	}

	// Buscar o filme atualizado
	getMovieByID(w, r, id)
}

// DELETE /movies/{id}
func deleteMovie(w http.ResponseWriter, r *http.Request, id int) {
	result, err := db.Exec("DELETE FROM movies WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Erro ao deletar o filme", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Erro ao verificar deleção", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Filme não encontrado"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Filme deletado com sucesso"})
}
