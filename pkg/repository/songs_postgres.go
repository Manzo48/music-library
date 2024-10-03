package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"unicode"

	"music-library/pkg/api_client"
	"music-library/pkg/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
)


type SongRepository interface {
	GetSongText(ctx context.Context, songID int, verse int, limit int) ([]string, error)
	DeleteSong(ctx context.Context, songID int) error
	UpdateSong(ctx context.Context, song model.Song) error
	AddSong(ctx context.Context, song model.Song) (int, error)
    GetSongs(ctx context.Context, filter model.SongFilter, page int, pageSize int) (model.Response, error)
}


type postgresSongRepository struct {
	db      *pgxpool.Pool
	api     *api_client.APIClient // Клиент для Genius API
}

func NewPostgresSongRepository(db *pgxpool.Pool) SongRepository {
	baseURL := viper.GetString("api.base_url")
    accessToken := viper.GetString("api.access_token")

	apiClient := api_client.NewAPIClient(baseURL, accessToken)

	return &postgresSongRepository{
		db:  db,
		api: apiClient,
	}
}



// AddSong добавляет новую песню в базу и получает информацию о песне через API.
func (r *postgresSongRepository) AddSong(ctx context.Context, song model.Song) (int, error) {
    // Проверяем, существует ли песня в базе данных
    var existingSongID int
    checkSongQuery := "SELECT id FROM songs WHERE group_name = $1 AND song_name = $2"
    err := r.db.QueryRow(ctx, checkSongQuery, song.Group, song.Title).Scan(&existingSongID)

    // Получаем информацию о песне через API Genius
    songDetail, err := r.api.GetSongDetails(song.Group, song.Title)
    if err != nil {
        return 0, fmt.Errorf("failed to fetch song details from Genius API for group: %s, song: %s, error: %w", song.Group, song.Title, err)
    }

    // Проверяем наличие необходимых данных
    if songDetail.Link == "" {
        return 0, fmt.Errorf("invalid song details received from Genius API for group: %s, song: %s", song.Group, song.Title)
    }

    // Если ReleaseDate отсутствует, передаем nil для колонки даты
    var releaseDate interface{}
    if songDetail.ReleaseDate == "" {
        releaseDate = nil
    } else {
        releaseDate = songDetail.ReleaseDate
    }

    // Форматируем текст песни, добавляя переносы строк
    songDetail.Text = strings.ReplaceAll(songDetail.Text, "]", "]\n")

    // Вставляем данные в таблицу songs и получаем ID новой записи
    var newSongID int
    insertSongQuery := "INSERT INTO songs (group_name, song_name) VALUES ($1, $2) RETURNING id"
    err = r.db.QueryRow(ctx, insertSongQuery, song.Group, song.Title).Scan(&newSongID)
    if err != nil {
        return 0, fmt.Errorf("error adding song to database, group: %s, song: %s, error: %w", song.Group, song.Title, err)
    }

    // Вставляем данные в таблицу song_details
    insertDetailQuery := `
    INSERT INTO song_details 
    (song_id, release_date, text, link, artist, album, genre, duration, key, tempo) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`



    _, err = r.db.Exec(ctx, insertDetailQuery, 
        newSongID, 
        releaseDate, 
        songDetail.Text, 
        songDetail.Link, 
        songDetail.Artist, 
        songDetail.Album, 
        songDetail.Genre, 
        songDetail.Duration, 
        songDetail.Key, 
        songDetail.Tempo,
    )
    if err != nil {
        return 0, fmt.Errorf("error adding song details to database, song ID: %d, error: %w", newSongID, err)
    }

    return newSongID, nil
}

func (r *postgresSongRepository) GetSongText(ctx context.Context, songID int, verse int, limit int) ([]string, error) {
    var text string
    err := r.db.QueryRow(ctx, "SELECT text FROM song_details WHERE song_id = $1", songID).Scan(&text)
    if err != nil {
        return nil, fmt.Errorf("ошибка получения текста песни: %w", err)
    }
    
    // Вызываем функцию для форматирования текста
    formattedLines := formatSongText(text)
    
    // Пагинация
    start := (verse - 1) * limit
    if start > len(formattedLines) {
        return nil, nil // Если строк меньше, чем запрашивается
    }
    end := start + limit
    if end > len(formattedLines) {
        end = len(formattedLines)
    }
    
    return formattedLines[start:end], nil
}








// GetSongs получает песни с фильтрацией  и пагинацией
func (r *postgresSongRepository) GetSongs(ctx context.Context, filter model.SongFilter, page int, pageSize int) (model.Response, error) {
    var response model.Response
    var songs []model.Song
    
    queryBuilder := strings.Builder{}
    queryBuilder.WriteString(`
    SELECT s.id, s.group_name, s.song_name, d.artist, COALESCE(d.album, '') as album, d.release_date, d.text, d.link
    FROM songs s
    LEFT JOIN song_details d ON s.id = d.song_id
    `)
    
    countQueryBuilder := strings.Builder{}
    countQueryBuilder.WriteString(`
        SELECT COUNT(*)
        FROM songs s
        LEFT JOIN song_details d ON s.id = d.song_id
        `)

        var args []interface{}
        var conditions []string
        argPos := 1
        
        if filter.Group != "" {
        conditions = append(conditions, fmt.Sprintf("s.group_name ILIKE $%d", argPos))
        args = append(args, "%"+filter.Group+"%")
        argPos++
    }
    
    if filter.Artist != "" {
        conditions = append(conditions, fmt.Sprintf("d.artist ILIKE $%d", argPos))
        args = append(args, "%"+filter.Artist+"%")
        argPos++
    }
    
    if filter.Album != "" {
        conditions = append(conditions, fmt.Sprintf("COALESCE(d.album, '') ILIKE $%d", argPos))
        args = append(args, "%"+filter.Album+"%")
        argPos++
    }
    
    if filter.Song != "" {
        conditions = append(conditions, fmt.Sprintf("s.song_name ILIKE $%d", argPos))
        args = append(args, "%"+filter.Song+"%")
        argPos++
    }
    
    if filter.ReleaseDate != "" {
        conditions = append(conditions, fmt.Sprintf("d.release_date = $%d", argPos))
        args = append(args, filter.ReleaseDate)
        argPos++
    }
    
    if len(conditions) > 0 {
        whereClause := " WHERE " + strings.Join(conditions, " AND ")
        queryBuilder.WriteString(whereClause)
        countQueryBuilder.WriteString(whereClause)
    }
    
    var totalCount int
    countQuery := countQueryBuilder.String()
    err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
    if err != nil {
        return response, fmt.Errorf("error counting songs: %w", err)
    }
    
    queryBuilder.WriteString(" ORDER BY s.id ASC")
    queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argPos, argPos+1))
    args = append(args, pageSize, (page-1)*pageSize)
    
    finalQuery := queryBuilder.String()
    
    rows, err := r.db.Query(ctx, finalQuery, args...)
    if err != nil {
        return response, fmt.Errorf("error querying songs with filters: %w", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var song model.Song
        var songDetails model.SongDetail
        
        var artist, album, text, link *string
        var releaseDateTime sql.NullTime
        
        err := rows.Scan(&song.ID, &song.Group, &song.Title, &artist, &album, &releaseDateTime, &text, &link)
        if err != nil {
            return response, fmt.Errorf("error scanning song: %w", err)
        }
        
        if artist != nil {
            songDetails.Artist = *artist
        }
        if album != nil {
            songDetails.Album = *album
        }
        if releaseDateTime.Valid {
            songDetails.ReleaseDate = releaseDateTime.Time.Format("2006-01-02")
        }
        if text != nil {
            songDetails.Text = *text
        }
        if link != nil {
            songDetails.Link = *link
        }
        
        song.Details = songDetails
        songs = append(songs, song)
    }
    
    if err := rows.Err(); err != nil {
        return response, fmt.Errorf("error iterating songs: %w", err)
    }
    
    totalPages := (totalCount + pageSize - 1) / pageSize
    response = model.Response{
        Songs: songs,
        Pagination: model.Pagination{
            Page:       page,
            PageSize:   pageSize,
            TotalCount: totalCount,
            TotalPages: totalPages,
        },
    }
    
    return response, nil
}

func (r *postgresSongRepository) UpdateSong(ctx context.Context, song model.Song) error {
   
    _, err := r.db.Exec(ctx, "UPDATE songs SET group_name = $1, song_name = $2 WHERE id = $3", 
    song.Group, song.Title, song.ID)
    if err != nil {
        return err
    }
    

    _, err = r.db.Exec(ctx, `
    UPDATE song_details 
    SET album = $1, text = $2 
    WHERE song_id = $3`, 
    song.Details.Album, song.Details.Text, song.ID)
    
    return err
}


func (r *postgresSongRepository) DeleteSong(ctx context.Context, songID int) error {
    _, err := r.db.Exec(ctx, "DELETE FROM songs WHERE id = $1", songID)
    return err
}



// Функция для форматирования текста песни в стихи
func formatSongText(text string) []string {
    // Разделяем текст по заголовкам, таким как [Verse], [Chorus] и т.д.
    text = strings.ReplaceAll(text, "[Verse", "\n[Verse")
    text = strings.ReplaceAll(text, "[Chorus", "\n[Chorus")
    text = strings.ReplaceAll(text, "[Bridge", "\n[Bridge")

  
    text = strings.ReplaceAll(text, ". ", ".\n")
    text = strings.ReplaceAll(text, "! ", "!\n")
    text = strings.ReplaceAll(text, "? ", "?\n")


    lines := strings.Split(text, "\n")

    // Теперь проверяем каждую строку, чтобы добавить переносы в случае, если строки склеены
    var formattedLines []string
    for _, line := range lines {
        // Убираем лишние пробелы
        trimmedLine := strings.TrimSpace(line)

        // Добавляем пробел перед заглавной буквой, если слово склеено с предыдущим
        fixedLine := insertSpaceBeforeUppercase(trimmedLine)

        if fixedLine != "" {
            formattedLines = append(formattedLines, fixedLine)
        }
    }

    return formattedLines
}

// Функция, которая добавляет пробел перед заглавной буквой, если перед ней нет пробела
func insertSpaceBeforeUppercase(line string) string {
    var result []rune
    for i, char := range line {
        if i > 0 && unicode.IsUpper(char) && line[i-1] != ' ' {
            result = append(result, ' ') 
        }
        result = append(result, char)
    }
    return string(result)
}
