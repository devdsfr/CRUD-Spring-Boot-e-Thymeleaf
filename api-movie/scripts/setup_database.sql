CREATE DATABASE IF NOT EXISTS meu_projeto_tmdb;

USE meu_projeto_tmdb;

CREATE TABLE IF NOT EXISTS movies (
                                      id INT PRIMARY KEY AUTO_INCREMENT,
                                      title VARCHAR(255) NOT NULL,
    description TEXT,
    release_date DATE
    -- Adicione outros campos conforme necessário
    );