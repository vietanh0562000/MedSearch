package app

import (
	"MedSearch/app/crawler"
	"MedSearch/app/database"
)

func Start() {
	database.Connect("mongodb://localhost:27017", "medsearch")

	crawler := crawler.Crawler{}
	crawler.Start()
}
