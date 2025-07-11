package app

import (
	"MedSearch/app/config"
	"MedSearch/app/crawler"
	"MedSearch/app/database"
	"MedSearch/app/database/repository"
)

func Start() {
	database.Connect("mongodb://localhost:27017", "medsearch")
	repository.CreateDrugTextIndex()

	config := config.Config{
		BaseURL:  "https://nhathuoclongchau.com.vn",
		StartURL: "https://nhathuoclongchau.com.vn/thuoc",
	}
	crawler := crawler.NewCrawler(&config)
	crawler.Start()
}
