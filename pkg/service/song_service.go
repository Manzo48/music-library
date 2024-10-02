package service

import (
	"context"
	"music-library/pkg/error_message"
	"music-library/pkg/model"
	"music-library/pkg/repository"

	"github.com/sirupsen/logrus"
)

// SongService интерфейс для бизнес-логики управления песнями.
type SongService interface {
	GetSongs(ctx context.Context, filter model.SongFilter, page int, pageSize int) (model.Response, error) // Новая функция
	GetSongText(ctx context.Context, songID int, verse int, limit int) ([]string, error)
	DeleteSong(ctx context.Context, songID int) error
	UpdateSong(ctx context.Context, song model.Song) error
	AddSong(ctx context.Context, song model.Song) (int, error)
}

// songService реализует интерфейс SongService.
type songService struct {
	repo   repository.SongRepository
	logger *logrus.Logger
}

// NewSongService создает новый экземпляр songService.
func NewSongService(repo repository.SongRepository, logger *logrus.Logger) SongService {
	return &songService{
		repo:   repo,
		logger: logger,
	}
}

// GetSongs возвращает песни с учетом фильтрации и пагинации.
func (s *songService) GetSongs(ctx context.Context, filter model.SongFilter, page int, pageSize int) (model.Response, error) {
	s.logger.Debug("Fetching songs with filter: ", filter, " page: ", page, " pageSize: ", pageSize)

	// Получаем песни с фильтрацией и пагинацией из репозитория
	response, err := s.repo.GetSongs(ctx, filter, page, pageSize)
	if err != nil {
		s.logger.Error("Error fetching filtered songs: ", err)
		return model.Response{}, error_message.ErrInternal
	}

	s.logger.Debug("Fetched songs successfully: ", response)
	return response, nil
}

// GetSongText возвращает текст песни по ID и номеру куплета.
func (s *songService) GetSongText(ctx context.Context, songID int, verse int, limit int) ([]string, error) {
	s.logger.Debug("Fetching text for song ID: ", songID, " verse: ", verse, " limit: ", limit)

	verses, err := s.repo.GetSongText(ctx, songID, verse, limit)
	if err != nil {
		s.logger.Error("Error fetching song text: ", err)
		return nil, error_message.ErrInternal
	}

	s.logger.Debug("Fetched song text successfully: ", verses)
	return verses, nil
}

// DeleteSong удаляет песню по ID.
func (s *songService) DeleteSong(ctx context.Context, songID int) error {
	s.logger.Debug("Attempting to delete song ID: ", songID)

	if err := s.repo.DeleteSong(ctx, songID); err != nil {
		s.logger.Error("Failed to delete song: ", err)
		return error_message.ErrInternal
	}

	s.logger.Info("Song deleted successfully: ", songID)
	return nil
}

// UpdateSong обновляет данные песни.
func (s *songService) UpdateSong(ctx context.Context, song model.Song) error {
	s.logger.Debug("Attempting to update song: ", song)

	if err := s.repo.UpdateSong(ctx, song); err != nil {
		s.logger.Error("Failed to update song: ", err)
		return error_message.ErrInternal
	}

	s.logger.Info("Song updated successfully: ", song.ID)
	return nil
}

// AddSong добавляет новую песню и получает дополнительную информацию через Genius API.
func (s *songService) AddSong(ctx context.Context, song model.Song) (int, error) {
	s.logger.Debug("Attempting to add a new song: ", song)
	
	newSongID, err := s.repo.AddSong(ctx, song)
	if err != nil {
		s.logger.Error("Failed to add song to repository: ", err)
		return 0, error_message.ErrInternal
	}

	s.logger.Info("Song added successfully with ID: ", newSongID)
	return newSongID, nil
}
