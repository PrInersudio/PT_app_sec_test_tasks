package main

import (
	"float_service/config"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Запускать: %s <файл конфигурации>", os.Args[0])
	}
	cfg := config.MustLoad(os.Args[1])
	log := setupLogger(cfg.Env)
	log.Info("Запуск float_service", slog.String("env", cfg.Env))
	log.Debug("Логгирование запущено на уровне DEBUG.")
	router := chi.NewRouter()
	// добавление к каждому запросу ID, чтобы потом отслеживать, что пошло не так
	router.Use(middleware.RequestID)
	// логгирование запросов
	router.Use(middleware.Logger)
	// восстановление в случае паники у обработчика
	router.Use(middleware.Recoverer)
	// "красивые" url у обработчиков
	router.Use(middleware.URLFormat)
}

// настройка логгирования
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		logfile, err := os.OpenFile(
			"app.log",
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666,
		)
		if err != nil {
			fmt.Println("Ошибка создания файла логгирования:", err)
			os.Exit(1)
		}
		log = slog.New(
			slog.NewTextHandler(
				logfile,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}
	return log
}
