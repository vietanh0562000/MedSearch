package app

import (
	"MedSearch/app/config"
	"MedSearch/app/crawler"
	"MedSearch/app/database"
	"MedSearch/app/database/repository"
	"MedSearch/app/logger"
	"sync"
)

type App struct {
	Config  *config.Config
	Logger  *logger.MLogger
	Crawler *crawler.Crawler
}

var Instance *App
var once sync.Once

func Start() {
	once.Do(func() {
		var newApp App
		database.Connect("mongodb://localhost:27017", "medsearch")
		repository.CreateDrugTextIndex()

		config := config.Config{
			BaseURL:  "https://nhathuoclongchau.com.vn",
			StartURL: "https://nhathuoclongchau.com.vn/thuoc",
		}

		logger := logger.NewLogger("app.log")

		crawler := crawler.NewCrawler(&config, logger)
		crawler.Start()

		newApp.Config = &config
		newApp.Crawler = crawler
		newApp.Logger = logger

		Instance = &newApp
	})

}
