// pkg/handler/song_handler.go
package handler

import (
	"music-library/pkg/error_message"
	"music-library/pkg/model"
	"music-library/pkg/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


type SongHandler struct {
	service service.SongService
	logger   *logrus.Logger
}


func NewSongHandler(service service.SongService, logger *logrus.Logger) *SongHandler {
	return &SongHandler{
		service: service,
		logger:  logger,
	}
}


func (h *SongHandler) InitRoutes(router *gin.Engine) {
	songs := router.Group("/songs")
	{	
		songs.GET("/", h.GetSongs)
		songs.GET("/:id", h.GetSongByID)
		songs.POST("/", h.AddSong)
		songs.PUT("/:id", h.UpdateSong)
		songs.DELETE("/:id", h.DeleteSong)
	}
}
// GetSongs получает список песен с фильтрацией и пагинацией
// @Summary Get list of songs
// @Description Get songs with optional filters (group, artist, album, song) and pagination
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string false "Filter by group"
// @Param artist query string false "Filter by artist"
// @Param album query string false "Filter by album"
// @Param song query string false "Filter by song title"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} []model.Song "List of songs"
// @Failure 400 {object} error_message.ErrorResponse "Invalid input"
// @Failure 500 {object} error_message.ErrorResponse "Internal server error"
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
    // Получаем параметры фильтров из запроса
    group := c.Query("group")
    artist := c.Query("artist")
    album := c.Query("album")
    song := c.Query("song")


    pageStr := c.DefaultQuery("page", "1")
    pageSizeStr := c.DefaultQuery("pageSize", "10")


    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        h.logger.Warn("Invalid page number: ", pageStr)
        c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid page number"})
        return
    }

  
    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil || pageSize <= 0 {
        h.logger.Warn("Invalid page size: ", pageSizeStr)
        c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid page size"})
        return
    }

   
    if pageSize > 100 {
        pageSize = 100
    }

  
    filter := model.SongFilter{
        Group:      group,
        Artist:     artist,
        Album:      album,
        Song:       song,
    }

   
    response, err := h.service.GetSongs(c.Request.Context(), filter, page, pageSize)
    if err != nil {
        h.logger.Error("GetSongs error: ", err)
        c.JSON(http.StatusInternalServerError, error_message.ErrorResponse{Message: "Internal server error"})
        return
    }

    
    c.JSON(http.StatusOK, response)
}


// GetSongByID получает песню по ID

// GetSongByID получает текст песни по ID с указанием куплета и лимита строк
// @Summary Get song text by ID
// @Description Get song text by song ID, verse, and limit
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param verse query int false "Verse number" default(1)
// @Param limit query int false "Number of lines per verse" default(4)
// @Success 200 {object} []string "Song text"
// @Failure 400 {object} error_message.ErrorResponse "Invalid input"
// @Failure 404 {object} error_message.ErrorResponse "Song not found"
// @Failure 500 {object} error_message.ErrorResponse "Internal server error"
// @Router /songs/{id} [get]
func (h *SongHandler) GetSongByID(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        h.logger.Warn("Invalid song ID: ", idParam)
        c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid ID"})
        return
    }
    
   
    verseParam := c.DefaultQuery("verse", "1") 
    verse, err := strconv.Atoi(verseParam)
    if err != nil {
        h.logger.Warn("Invalid verse parameter: ", verseParam)
        c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid verse parameter"})
        return
    }

    limitParam := c.DefaultQuery("limit", "4") 
    limit, err := strconv.Atoi(limitParam)
    if err != nil {
        h.logger.Warn("Invalid limit parameter: ", limitParam)
        c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid limit parameter"})
        return
    }

    text, err := h.service.GetSongText(c.Request.Context(), id, verse, limit)
    if err != nil {
        if err.Error() == "song not found" {
            c.JSON(http.StatusNotFound, error_message.ErrorResponse{Message: "Song not found"})
            return
        }
        h.logger.Error("GetSongByID error: ", err)
        c.JSON(http.StatusInternalServerError, error_message.ErrorResponse{Message: "Internal server error"})
        return
    }

    c.JSON(http.StatusOK, text)
}


// AddSong добавляет новую песню
// @Summary Add a new song
// @Description Add a new song to the library
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body model.Song true "Song data"
// @Success 201 {object} map[string]interface{} "Song added successfully with ID"
// @Failure 400 {object} error_message.ErrorResponse "Invalid request body"
// @Failure 502 {object} error_message.ErrorResponse "Failed to fetch data from external API"
// @Failure 500 {object} error_message.ErrorResponse "Internal server error"
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var input model.Song
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("AddSong binding error: ", err)
		c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid request body"})
		return
	}

	newSongID, err := h.service.AddSong(c.Request.Context(), input)
	if err != nil {
		if err.Error() == error_message.ErrExternalAPI.Error() {
			c.JSON(http.StatusBadGateway, error_message.ErrorResponse{Message: "Failed to fetch data from external API"})
			return
		}
		h.logger.Error("AddSong error: ", err)
		c.JSON(http.StatusInternalServerError, error_message.ErrorResponse{Message: "Failed to add song"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Song added successfully", "id": newSongID})
}


// UpdateSong обновляет данные песни
// @Summary Update an existing song
// @Description Update song details by ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "Song ID"
// @Param song body model.Song true "Updated song data"
// @Success 200 {object} map[string]interface{} "Song updated successfully"
// @Failure 400 {object} error_message.ErrorResponse "Invalid input"
// @Failure 500 {object} error_message.ErrorResponse "Internal server error"
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Warn("Invalid song ID: ", idParam)
		c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid ID"})
		return
	}

	var song model.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.logger.Warn("UpdateSong binding error: ", err)
		c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid request body"})
		return
	}

	song.ID = id
	if err := h.service.UpdateSong(c.Request.Context(), song); err != nil {
		h.logger.Error("UpdateSong error: ", err)
		c.JSON(http.StatusInternalServerError, error_message.ErrorResponse{Message: "Failed to update song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// DeleteSong удаляет песню по ID
// @Summary Delete a song by ID
// @Description Remove a song from the library by its ID
// @Tags songs
// @Produce  json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]interface{} "Song deleted successfully"
// @Failure 400 {object} error_message.ErrorResponse "Invalid song ID"
// @Failure 500 {object} error_message.ErrorResponse "Internal server error"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Warn("Invalid song ID: ", idParam)
		c.JSON(http.StatusBadRequest, error_message.ErrorResponse{Message: "Invalid ID"})
		return
	}

	if err := h.service.DeleteSong(c.Request.Context(), id); err != nil {
		h.logger.Error("DeleteSong error: ", err)
		c.JSON(http.StatusInternalServerError, error_message.ErrorResponse{Message: "Failed to delete song"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
