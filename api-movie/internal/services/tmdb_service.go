package services

import (
	"api-movie/internal/models"
	"fmt"
)

// MovieService representa o serviço de filmes
type MovieService struct {
	Client *TMDBClient
}

// NewMovieService cria uma nova instância de MovieService
func NewMovieService(client *TMDBClient) *MovieService {
	return &MovieService{Client: client}
}

// GetMovieDetails busca e retorna os detalhes de um filme.
func (m *MovieService) GetMovieDetails(query string) (*models.Movie, error) {
	movies, err := m.Client.SearchMovies(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar filmes: %w", err) // Wrap the error
	}

	if len(movies.Results) == 0 {
		return nil, nil // Ou retorne um erro se nenhum filme for encontrado, dependendo da sua lógica
	}

	return &movies.Results[0], nil // Retorna o primeiro resultado da busca
}
