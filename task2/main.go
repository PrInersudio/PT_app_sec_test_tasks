package main

import (
	"FloatService/config"
	"FloatService/floatcalculation"
	"FloatService/handlers/handlefloatcalculation"
	mwLogger "FloatService/middleware/logger"
	"FloatService/response"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Запускать: %s <файл конфигурации>\n", os.Args[0])
		os.Exit(1)
	}
	cfg := config.MustLoad(os.Args[1])
	log, logfile := setupLogger(cfg.Env)
	if logfile != nil {
		defer logfile.Close()
	}
	log.Info("Запуск FloatService", slog.String("env", cfg.Env))
	log.Debug("Логгирование запущено на уровне DEBUG.")
	router := chi.NewRouter()
	// добавление к каждому запросу ID, чтобы потом отслеживать, что пошло не так
	router.Use(middleware.RequestID)
	// логгирование запросов
	router.Use(mwLogger.New(log))
	// восстановление в случае паники у обработчика
	router.Use(middleware.Recoverer)
	//добавление ограничения на количество запросов
	router.Use(httprate.Limit(
		cfg.Limit, cfg.Interval,
		httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
			log.Warn("Достигнут лимит запросов.")
			// Собственный http ответ с ошибкой 402 и json ответом из файла конфигурации
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(402)
			resp, _ := json.Marshal(response.Error(cfg.Msg))
			w.Write(append(resp, byte('\n')))
		})))
	// добавляем обработчик
	router.Get("/", handlefloatcalculation.New(log, &floatcalculation.FloatCalculator{}))
	log.Info("Запускаем сервер.", slog.String("address", cfg.Address))
	// обработка прерываний
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}
	// выносим запуск сервера в отдельную Go рутину
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Ошибка сервера.", slog.String("error", err.Error()))
		}
	}()
	log.Info("Сервер запущен")
	<-done
	log.Info("Остановка сервера.")
	// сервер остановится через timepout времени, если есть открытые подключения, иначе мгновенно
	ctx, cancel := context.WithTimeout(context.Background(), cfg.StopTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Ошибка остановки сервера", slog.String("error", err.Error()))
		return
	}
	log.Info("Сервер остановлен.")
}

// настройка логгирования
func setupLogger(env string) (log *slog.Logger, logfile *os.File) {
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
		logfile = nil
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		logfile = nil
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}
	return
}
