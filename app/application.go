package app

import (
	"MedSearch/app/config"
	"MedSearch/app/crawler"
	"MedSearch/app/database"
)

func Start() {
	database.Connect("mongodb://localhost:27017", "medsearch")

	config := config.Config{
		BaseURL:  "https://nhathuoclongchau.com.vn",
		StartURL: "https://nhathuoclongchau.com.vn/thuoc",
	}
	crawler := crawler.NewCrawler(&config)
	crawler.Start()
}
