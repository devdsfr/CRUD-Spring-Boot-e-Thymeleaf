package models

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	// Adicione outros campos conforme necessário
}
