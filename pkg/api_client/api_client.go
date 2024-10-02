package api_client

import (
    "encoding/json"
    "fmt"
    "music-library/pkg/model"
    "net/http"
    "net/url"

    "github.com/PuerkitoBio/goquery" // Импортируем goquery для парсинга страницы
)

type APIClient struct {
    BaseURL     string
    AccessToken string
}

func NewAPIClient(baseURL, accessToken string) *APIClient {
    return &APIClient{
        BaseURL:     baseURL,
        AccessToken: accessToken,
    }
}

// GetSongDetails возвращает информацию о песне, включая текст и метаданные
// Package api_client provides a client for interacting with the Genius API.
//
// This package uses the Genius API to fetch song details, including the song text.
// The Genius API is a third-party service that provides information about songs, artists, and albums.
//
// To use the Genius API, you'll need to obtain an access token. You can sign up for a Genius API account
// and generate an access token on the Genius API website.

func (c *APIClient) GetSongDetails(group, song string) (*model.SongDetail, error) {
    // Создаем URL для поиска песни через API Genius
    endpoint, err := url.Parse(c.BaseURL + "/search")
    if err != nil {
        return nil, fmt.Errorf("error parsing URL: %w", err)
    }

    query := endpoint.Query()
    query.Set("q", fmt.Sprintf("%s %s", group, song)) // Добавляем название группы и песни
    endpoint.RawQuery = query.Encode()

    // Создаем HTTP запрос
    req, err := http.NewRequest("GET", endpoint.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // Добавляем авторизационный токен
    req.Header.Add("Authorization", "Bearer "+c.AccessToken)

    // Отправляем запрос
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    // Проверяем код ответа
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
    }

    // Декодируем JSON ответ
    var searchResult struct {
        Response struct {
            Hits []struct {
                Result struct {
                    ID int `json:"id"`
                } `json:"result"`
            } `json:"hits"`
        } `json:"response"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
        return nil, fmt.Errorf("error decoding JSON response: %w", err)
    }

    // Проверяем, есть ли результаты поиска
    if len(searchResult.Response.Hits) == 0 {
        return nil, fmt.Errorf("no song found for group: %s, song: %s", group, song)
    }

    // Получаем ID первой найденной песни
    songID := searchResult.Response.Hits[0].Result.ID

    // Запрашиваем детали песни по ID
    songDetail, err := c.getSongDetailByID(songID)
    if err != nil {
        return nil, fmt.Errorf("error fetching song details: %w", err)
    }

    // Возвращаем объект SongDetail
    return songDetail, nil
}

// getSongDetailByID запрашивает детали песни по её ID
func (c *APIClient) getSongDetailByID(songID int) (*model.SongDetail, error) {
    // Создаем URL для получения деталей песни
    endpoint, err := url.Parse(fmt.Sprintf("%s/songs/%d", c.BaseURL, songID))
    if err != nil {
        return nil, fmt.Errorf("error parsing URL: %w", err)
    }

    // Создаем HTTP запрос
    req, err := http.NewRequest("GET", endpoint.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request: %w", err)
    }

    // Добавляем авторизационный токен
    req.Header.Add("Authorization", "Bearer "+c.AccessToken)

    // Отправляем запрос
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %w", err)
    }
    defer resp.Body.Close()

    // Проверяем код ответа
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
    }

    // Декодируем JSON ответ
    var songResponse struct {
        Response struct {
            Song struct {
                ID          int    `json:"id"`
                Title       string `json:"title"`
                Artist      struct {
                    Name string `json:"name"`
                } `json:"primary_artist"`
                URL          string `json:"url"`
                Album        struct {
                    Name string `json:"name"`
                } `json:"album"` // Изменено на структуру
                ReleaseDate  string `json:"release_date"`
                Genre        string `json:"genre"`
                Duration     string `json:"duration"`
                Key          string `json:"key"`
                Tempo        string `json:"tempo"`
                // Добавьте другие поля, если необходимо
            } `json:"song"`
        } `json:"response"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&songResponse); err != nil {
        return nil, fmt.Errorf("error decoding JSON response: %w", err)
    }

    // Получаем текст песни через парсинг страницы
    lyrics, err := c.GetSongLyrics(songResponse.Response.Song.URL)
    if err != nil {
        return nil, fmt.Errorf("error fetching song lyrics: %w", err)
    }

    // Создаем объект SongDetail и возвращаем его
    songDetail := &model.SongDetail{
        Link:            songResponse.Response.Song.URL,
        Artist:          songResponse.Response.Song.Artist.Name,
        Album:           songResponse.Response.Song.Album.Name, // Обновлено для получения имени альбома
        ReleaseDate:     songResponse.Response.Song.ReleaseDate,
        Text:            lyrics,
        Genre:           songResponse.Response.Song.Genre,
        Duration:        songResponse.Response.Song.Duration,
        Key:             songResponse.Response.Song.Key,
        Tempo:           songResponse.Response.Song.Tempo,
    }

    return songDetail, nil
}

// GetSongLyrics парсит текст песни с страницы Genius
func (c *APIClient) GetSongLyrics(songURL string) (string, error) {
    // Отправляем запрос на страницу с текстом песни
    resp, err := http.Get(songURL)
    if err != nil {
        return "", fmt.Errorf("failed to fetch song page: %w", err)
    }
    defer resp.Body.Close()

    // Проверка статуса ответа
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to fetch song page: status %d", resp.StatusCode)
    }

    // Парсим HTML страницы
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to parse song page: %w", err)
    }

    // Ищем текст песни
    lyrics := ""
    doc.Find(".lyrics").Each(func(i int, s *goquery.Selection) {
        lyrics += s.Text() + "\n"
    })

    // Если текст не найден, пробуем другой селектор
    if lyrics == "" {
        doc.Find("div[class^='Lyrics__Container']").Each(func(i int, s *goquery.Selection) {
            lyrics += s.Text() + "\n"
        })
    }

    if lyrics == "" {
        return "", fmt.Errorf("no lyrics found on the page")
    }

    return lyrics, nil
}
