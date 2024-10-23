package services

import (
	"api-movie/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

// const apiKey = "98cf6e1c64f9167cc7ee0ea8b6b4779b" // Substitua pela sua chave de API do TMDb
const tmdbURL = "https://api.themoviedb.org/3"

// HTTPClient interface para permitir a injeção de dependências
type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

// TMDBClient representa o cliente para a API do TMDB
type TMDBClient struct {
	APIKey  string
	BaseURL string
	Client  HTTPClient
}

// NewTMDBClient cria uma nova instância de TMDBClient
func NewTMDBClient(apiKey string, client HTTPClient) *TMDBClient {
	return &TMDBClient{
		APIKey:  apiKey,
		BaseURL: tmdbURL,
		Client:  client,
	}
}

// SearchMovies busca filmes com base na consulta fornecida
func (t *TMDBClient) SearchMovies(query string) (*models.APIResponse, error) {
	url := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s", t.BaseURL, t.APIKey, query)
	resp, err := t.Client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Erro na resposta da API: %s", resp.Status)
	}

	var apiResponse models.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
