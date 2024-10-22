package api

import (
	"api-movie/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

type TMDbClient struct {
	APIKey string
}

func NewTMDbClient(apiKey string) *TMDbClient {
	return &TMDbClient{APIKey: apiKey}
}

func (c *TMDbClient) FetchPopularMovies() ([]models.Movie, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/popular?api_key=%s", c.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Description string `json:"overview"`
			ReleaseDate string `json:"release_date"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var movies []models.Movie
	for _, m := range result.Results {
		movies = append(movies, models.Movie{
			ID:          m.ID,
			Title:       m.Title,
			Description: m.Description,
			ReleaseDate: m.ReleaseDate,
		})
	}

	return movies, nil
}
