package api_client_test

import (
	"fmt"
	"music-library/pkg/api_client"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Пример мок-сервера, который имитирует ответ от API
func TestGetSongDetails(t *testing.T) {
	// Мок-сервер
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"response": {
				"hits": [{
					"result": {
						"id": 12345,
						"title": "Fix You",
						"primary_artist": {
							"name": "Coldplay"
						},
						"url": "https://genius.com/The-weeknd-and-playboi-carti-timeless-lyrics",
						"album": "X&Y",
						"release_date": "2005-03-22",
						"genre": "Alternative Rock",
						"duration": "4:55",
						"key": "B♭",
						"tempo": "138"
					}
				}]
			}
		}`))
	}))
	defer mockServer.Close()

	// Инициализация API клиента
	client := api_client.NewAPIClient(mockServer.URL, "fake_token")

	// Вызов метода
	songDetails, err := client.GetSongDetails("Coldplay", "Fix You")

	// Вывод результата в консоль
	fmt.Println("Song Details:", songDetails)

	// Проверка на ошибку
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Проверка данных
	if songDetails.Artist != "Coldplay" {
		t.Errorf("Expected artist Coldplay, got %s", songDetails.Artist)
	}
	if songDetails.Link != "https://genius.com/The-weeknd-and-playboi-carti-timeless-lyrics" {
		t.Errorf("Expected Link %s, got %s", "https://genius.com/The-weeknd-and-playboi-carti-timeless-lyrics", songDetails.Link)
	}
	if songDetails.Album != "X&Y" {
		t.Errorf("Expected Album %s, got %s", "X&Y", songDetails.Album)
	}
	if songDetails.ReleaseDate != "2005-03-22" {
		t.Errorf("Expected ReleaseDate %s, got %s", "2005-03-22", songDetails.ReleaseDate)
	}
	if songDetails.Genre != "Alternative Rock" {
		t.Errorf("Expected Genre %s, got %s", "Alternative Rock", songDetails.Genre)
	}
	if songDetails.Duration != "4:55" {
		t.Errorf("Expected Duration %s, got %s", "4:55", songDetails.Duration)
	}
	if songDetails.Key != "B♭" {
		t.Errorf("Expected Key %s, got %s", "B♭", songDetails.Key)
	}
	if songDetails.Tempo != "138" {
		t.Errorf("Expected Tempo %s, got %s", "138", songDetails.Tempo)
	}
}
