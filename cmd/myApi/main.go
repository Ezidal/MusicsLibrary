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

	"github.com/gorilla/mux"
)

func main() {
	// Load config
	conf := config.LoadConfig()

	// Init logger
	log := logger.SetLogger(conf.Environment)
	log.Info("Logger starting, config setup")

	log.Debug("addres external api: " + conf.ExternalApiUrl)

	// Init db
	db, err := storage.New(conf)
	if err != nil {
		log.Error("failed to init db", er.Err(err))
		os.Exit(1)
	}
	log.Info("Sucsess connect to db: " + conf.DBHost + ":" + conf.DBPort)

	// Init server
	r := mux.NewRouter()
	// Add middleware
	r.Use(middleware.LoggerMiddleware(log))

	handler := handlers.NewHandler(db, log, conf)
	// Add handlers
	r.HandleFunc("/songs", handler.GetSongs).Methods("GET")
	// r.HandleFunc("/songs/{id}/text", handler.GetSongText).Methods("GET")
	// r.HandleFunc("/songs/{id}", handler.DeleteSong).Methods("DELETE")
	// r.HandleFunc("/songs/{id}", handler.UpdateSong).Methods("PUT")
	r.HandleFunc("/songs", handler.AddSong).Methods("POST")

	// Start server
	log.Info("Server starting on port: " + conf.ServerPort)
	err = http.ListenAndServe(":"+conf.ServerPort, r)
	if err != nil {
		log.Error(err.Error())
	}

}
