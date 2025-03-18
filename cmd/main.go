package main

import (
	"LibMusic/internal/config"
	"LibMusic/internal/logger"
	er "LibMusic/internal/logger/err"
	"LibMusic/internal/storage"
	"os"
)

func main() {
	conf := config.LoadConfig()
	log := logger.SetLogger(conf.Environment)
	log.Info("Logger starting, config setup")

	log.Debug("Адрес внешнего api: " + conf.ExternalApiUrl)

	_, err := storage.New(conf)
	if err != nil {
		log.Error("failed to init db", er.Err(err))
		os.Exit(1)
	}
	log.Info("Sucsess connect to db: " + conf.DBHost + ":" + conf.DBPort)

}
