package main

import (
	"LibMusic/internal/config"
	"LibMusic/internal/handlers"
	"LibMusic/internal/logger"
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/middleware"
	"LibMusic/internal/storage"
	"net/http"
	"os"

	_ "LibMusic/docs" // Импортируем сгенерированную документацию Swagger

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	// Пакет для Swagger UI
)

// @title Songs Library API
// @version 1.0
// @description API для управления библиотекой песен
// @host localhost:8080
// @BasePath /
func main() {
	// Загрузка конфигурации
	conf := config.LoadConfig()

	// Инициализация логгера
	log := logger.SetLogger(conf.Environment)
	log.Info("Logger starting, config setup")

	log.Debug("address external api: " + conf.ExternalApiUrl)

	// Инициализация базы данных
	db, err := storage.New(conf, log)
	if err != nil {
		log.Error("failed to init db", er.Err(err))
		os.Exit(1)
	}
	log.Info("Success connect to db: " + conf.DBHost + ":" + conf.DBPort)

	// Инициализация роутера
	r := mux.NewRouter()

	// Добавление middleware
	r.Use(middleware.LoggerMiddleware(log))

	// Инициализация мега-хендлера
	handler := handlers.NewHandler(db, log, conf)

	// Добавление обработчиков
	r.HandleFunc("/songs", handler.GetSongs).Methods("GET")
	r.HandleFunc("/songs/{id}/text", handler.GetSongText).Methods("GET")
	r.HandleFunc("/songs/{id}", handler.DeleteSong).Methods("DELETE")
	r.HandleFunc("/songs/{id}", handler.UpdateSong).Methods("PUT")
	r.HandleFunc("/songs", handler.AddSong).Methods("POST")

	// Добавление Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Запуск сервера
	log.Info("Server starting on port: " + conf.ServerPort)
	err = http.ListenAndServe(":"+conf.ServerPort, r)
	if err != nil {
		log.Error(err.Error())
	}
}
