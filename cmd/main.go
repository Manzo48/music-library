package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"music-library/pkg/handler"
	"music-library/pkg/repository"
	"music-library/pkg/service"
	_ "music-library/docs" // Импортируйте сгенерированные документации Swagger
)

// @title Music Library API
// @version 1.0
// @description This is a simple API for managing a music library.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT 

// @externalDocs.description Genius API
// @externalDocs.url https://docs.genius.com/

// @host localhost:9000

// @BasePath /
func main() {
	// Чтение конфигурации из файла, переданного в аргументах командной строки.
	if len(os.Args) < 2 {
		logrus.Fatalf("Usage: %v config_filename\n", os.Args[0])
	}
	fmt.Println("Using config file:", os.Args[1])

	// Инициализация конфигураций с помощью Viper
	if err := initConfig(os.Args[1]); err != nil {
		logrus.Fatalf("Error initializing configs: %s", err.Error())
	}

	// Инициализация базы данных с использованием параметров из конфигурации
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Username: viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.name"),
		SSLMode:  viper.GetString("database.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize DB: %s", err.Error())
	}

	// Инициализация репозитория для работы с песнями
	songRepo := repository.NewPostgresSongRepository(db)

	// Инициализация логгера для записи логов
	logger := logrus.New()

	// Инициализация сервисов, отвечающих за бизнес-логику
	services := service.NewSongService(songRepo, logger)

	// Создание нового маршрутизатора Gin
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery()) // Подключение логирования и обработки ошибок
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Инициализация обработчиков маршрутов для песен
	handlers := handler.NewSongHandler(services, logger)
	
	// Инициализация маршрутов, определенных в обработчиках
	handlers.InitRoutes(router)

	// Запуск HTTP-сервера
	srv := &http.Server{
		Addr:    ":9000",
		Handler: router,
	}

	go func() {
		// Запуск сервера и обработка ошибок
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to listen: %s", err.Error())
		}
	}()
	log.Println("Server started at", ":9000")

	// Обработка завершения работы сервера (graceful shutdown)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Gracefully shutting down...")

	// Установка тайм-аута для завершения работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Попытка завершить работу сервера
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Error occurred during server shutdown: %s", err.Error())
	}

	log.Println("Server exited properly")
}

// initConfig инициализирует конфигурацию с помощью Viper
func initConfig(filename string) error {
	viper.SetConfigFile(filename)
	return viper.ReadInConfig()
}
