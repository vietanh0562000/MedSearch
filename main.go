package main

import (
	"MedSearch/app"
	"MedSearch/app/config"
	"MedSearch/app/crawler"
	"MedSearch/app/logger"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [api|crawl]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "api":
		cfg := config.GetNewAppConfig("8080", "mongodb://localhost:27017", "medsearch")
		log := logger.NewLogger("app.log")
		app.Start(cfg, log)
	case "crawl":
		cfg := config.GetNewCrawlerConfig("https://nhathuoclongchau.com.vn", "https://nhathuoclongchau.com.vn/thuoc", "mongodb://localhost:27017", "medsearch")
		log := logger.NewLogger("app.log")
		cr := crawler.NewCrawler(cfg, log)
		cr.Start()
	default:
		fmt.Println("Unknown command:", os.Args[1])
		fmt.Println("Usage: go run main.go [api|crawl]")
		os.Exit(1)
	}
}
