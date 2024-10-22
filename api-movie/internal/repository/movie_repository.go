package repository

import (
	"api-movie/internal/models"
	"database/sql"
)

type MovieRepository struct {
	DB *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{DB: db}
}

func (r *MovieRepository) Create(movie *models.Movie) error {
	query := `INSERT INTO movies (title, description, release_date) VALUES (?, ?, ?)`
	_, err := r.DB.Exec(query, movie.Title, movie.Description, movie.ReleaseDate)
	return err
}

// Adicione outros métodos como Get, Update, Delete conforme necessário
