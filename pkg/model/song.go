package model

// Song представляет собой основную модель песни в библиотеке
type Song struct {
	ID      int       `json:"id"`       // Уникальный идентификатор песни
	Group   string    `json:"group"`    // Исполнитель или группа
	Title   string    `json:"song"`     // Название песни
	Details SongDetail `json:"details"` // Дополнительные детали о песне
}

// SongDetail представляет собой детализированную модель песни, полученной из Genius API
type SongDetail struct {
    Link        string `json:"link"`
    Artist      string `json:"artist"`
    Album       string `json:"album"`
    ReleaseDate string `json:"release_date"`
    Text        string `json:"text"`
    Genre       string `json:"genre"`      // Жанр
    Duration    string `json:"duration"`   // Продолжительность
    Key         string `json:"key"`         // Ключ
    Tempo       string `json:"tempo"`       // Темп
}

// SongFilter представляет собой модель фильтрации для поиска песен
type SongFilter struct {
	Group string `json:"group"` // Фильтр по исполнителю или группе
    Artist string `json:"artist"`
    Album string `json:"album"`
	Song  string `json:"song"`  
    ReleaseDate string `json:"release"`
}

// Pagination представляет собой модель для пагинации
type Pagination struct {
	Page        int `json:"page"`        // Номер страницы
	PageSize    int `json:"pageSize"`    // Количество элементов на странице
	TotalCount  int `json:"totalCount"`  // Общее количество элементов
	TotalPages  int `json:"totalPages"`  // Общее количество страниц
}

// Response представляет собой модель ответа для API
type Response struct {
	Songs      []Song     `json:"songs"`      // Список песен
	Pagination Pagination `json:"pagination"` // Пагинация
}
