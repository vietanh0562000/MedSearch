package app

import "MedSearch/app/crawler"

func Start() {
	crawler := crawler.Crawler{}
	crawler.Start()
}
