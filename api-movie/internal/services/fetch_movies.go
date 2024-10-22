package services

import (
	"api-movie/internal/api"
	"api-movie/internal/config"
	"api-movie/internal/repository"
	"database/sql"
	"fmt"
	"log"
)

type FetchMoviesService struct {
	TMDbClient *api.TMDbClient
	MovieRepo  *repository.MovieRepository
}

func NewFetchMoviesService(cfg *config.Config) (*FetchMoviesService, error) {
	// Inicialize a conexão com o banco de dados
	db, err := initDB(cfg)
	if err != nil {
		return nil, err
	}

	movieRepo := repository.NewMovieRepository(db)
	tmdbClient := api.NewTMDbClient(cfg.TMDbAPIKey)

	return &FetchMoviesService{
		TMDbClient: tmdbClient,
		MovieRepo:  movieRepo,
	}, nil
}

func initDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (s *FetchMoviesService) Execute() error {
	movies, err := s.TMDbClient.FetchPopularMovies()
	if err != nil {
		return err
	}

	for _, movie := range movies {
		if err := s.MovieRepo.Create(&movie); err != nil {
			log.Printf("Erro ao salvar filme %s: %v", movie.Title, err)
		}
	}

	return nil
}
