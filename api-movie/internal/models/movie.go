package models

// Movie representa os dados do filme
type Movie struct {
	Adult            bool    `json:"adult"`
	OriginalLanguage string  `json:"original_language"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
}

// APIResponse representa a estrutura da resposta da API
type APIResponse struct {
	Results []Movie `json:"results"`
}
