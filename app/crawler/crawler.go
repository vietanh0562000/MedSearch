package crawler

import (
	"MedSearch/app/database/repository"
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Crawler struct {
}

func (c *Crawler) Start() {

	//TODO: Create collector by colly
	collector := colly.NewCollector()
	//TODO: Set rate limit for collector: OnRequest, OnResponse. OnHTML, OnError
	// Get Info in OnHTML function
	//TODO: Set callback for collector
	// 1. Thương hiệu
	collector.OnHTML("body", func(e *colly.HTMLElement) {
		drug := ParseDrug(e)
		err := repository.InsertDrug(&drug)
		if err != nil {
			fmt.Println("Failed to insert drug:", err)
		} else {
			fmt.Println("Inserted drug:", drug.Name)
		}
	})

	//TODO: Visit(url)
	url := "https://nhathuoclongchau.com.vn/thuoc/exopadin-60-mg-truong-tho-3-x10.html"
	collector.Visit(url)
	collector.Wait()
}
