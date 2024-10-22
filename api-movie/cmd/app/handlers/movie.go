package handlers

import (
	"encoding/json"
	"net/http"
)

// Movie representa um filme
type Movie struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Year  int    `json:"year"`
}

// Armazenamento em memória para exemplificação
var movies = []Movie{
	{ID: "1", Title: "The Matrix", Year: 1999},
	{ID: "2", Title: "Inception", Year: 2010},
}

// Handler para o endpoint /movies
func MoviesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMovies(w, r)
	case http.MethodPost:
		createMovie(w, r)
	case http.MethodPut:
		updateMovie(w, r)
	case http.MethodDelete:
		deleteMovie(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método não suportado"})
	}
}

// GET /movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// POST /movies
func createMovie(w http.ResponseWriter, r *http.Request) {
	var newMovie Movie
	err := json.NewDecoder(r.Body).Decode(&newMovie)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dados inválidos"})
		return
	}
	// Adicionar ID de forma simplificada
	newMovie.ID = string(len(movies) + 1)
	movies = append(movies, newMovie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newMovie)
}

// PUT /movies
func updateMovie(w http.ResponseWriter, r *http.Request) {
	var updatedMovie Movie
	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Dados inválidos"})
		return
	}

	for i, movie := range movies {
		if movie.ID == updatedMovie.ID {
			movies[i].Title = updatedMovie.Title
			movies[i].Year = updatedMovie.Year

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(movies[i])
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Filme não encontrado"})
}

// DELETE /movies?id={id}
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID é necessário"})
		return
	}

	for i, movie := range movies {
		if movie.ID == id {
			movies = append(movies[:i], movies[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Filme deletado com sucesso"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Filme não encontrado"})
}
