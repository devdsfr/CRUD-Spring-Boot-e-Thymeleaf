package services

import (
	"api-movie/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

// MockHTTPClient é uma implementação do HTTPClient para testes
type MockHTTPClient struct {
	Response *http.Response
	Error    error
}

// Get simula a chamada HTTP
func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Error
}

func TestSearchMovies(t *testing.T) {
	// Cria um mock de resposta da API
	movie := models.Movie{
		Adult:            false,
		OriginalLanguage: "en",
		Overview:         "A great movie",
		Popularity:       10.0,
	}
	apiResponse := models.APIResponse{
		Results: []models.Movie{movie},
	}

	// Serializa a resposta em JSON
	responseBody, _ := json.Marshal(apiResponse)

	// Teste 1: Resposta bem-sucedida
	t.Run("success", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBuffer(responseBody)),
			},
			Error: nil,
		}

		client := NewTMDBClient("fake_api_key", mockClient)
		result, err := client.SearchMovies("great movie")

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(result.Results) != 1 {
			t.Fatalf("Expected 1 movie, got %d", len(result.Results))
		}

		if result.Results[0].Overview != movie.Overview {
			t.Fatalf("Expected overview %s, got %s", movie.Overview, result.Results[0].Overview)
		}
	})

	// Teste 2: Resposta com erro
	t.Run("error_response", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBuffer([]byte(`{"status_code": 7, "status_message": "Invalid API key"}`))),
			},
			Error: nil,
		}

		client := NewTMDBClient("fake_api_key", mockClient)
		_, err := client.SearchMovies("great movie")

		if err == nil {
			t.Fatal("Expected an error, got none")
		}
	})

	// Teste 3: Erro na chamada HTTP
	t.Run("http_error", func(t *testing.T) {
		mockClient := &MockHTTPClient{
			Response: nil,
			Error:    fmt.Errorf("network error"),
		}

		client := NewTMDBClient("fake_api_key", mockClient)
		_, err := client.SearchMovies("great movie")

		if err == nil {
			t.Fatal("Expected an error, got none")
		}
	})
}
