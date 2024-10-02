package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Конфигурация базы данных
type Config struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"` // Добавьте это, если используете ssl
}


func NewPostgresDB(cfg Config) (*pgxpool.Pool, error) {
	// Формируем строку подключения
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	// Настраиваем параметры пула
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Printf("Unable to parse config: %v\n", err)
		return nil, err
	}

	// Настраиваем контекст с тайм-аутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Инициализируем пул соединений
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return nil, err
	}

	// Пинг базы данных для проверки соединения
	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("Unable to ping database: %v\n", err)
		return nil, err
	}

	log.Println("Connected to database successfully")
	return pool, nil
}
