package crawler

import (
	"MedSearch/app/config"
	"MedSearch/app/database/repository"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	config config.Config
}

func NewCrawler(config *config.Config) *Crawler {
	newCrawler := Crawler{
		config: *config,
	}

	return &newCrawler
}

func (c *Crawler) Start() {
	log.Println("🚀 Starting crawler...")

	//TODO: Create collector by colly
	collector := colly.NewCollector(
		colly.MaxDepth(3), // Increased depth
	)

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		Delay:       1 * time.Second, // Reduced delay
		RandomDelay: 1 * time.Second,
	})

	// Add error handling
	collector.OnError(func(r *colly.Response, err error) {
		log.Printf("❌ Error visiting %s: %v", r.Request.URL, err)
	})

	//TODO: Set rate limit for collector: OnRequest, OnResponse. OnHTML, OnError
	// Get Info in OnHTML function
	collector.OnRequest(func(r *colly.Request) {
		log.Printf("🌐 Visiting: %s", r.URL.String())
	})

	//TODO: Set callback for collector
	collector.OnHTML("body", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		log.Printf("📄 Processing: %s", url)

		if strings.Contains(url, "/thuoc/") && strings.Contains(url, ".html") {
			log.Printf("💊 Found drug page: %s", url)
			drug := ParseDrug(e)
			err := repository.InsertDrug(&drug)
			if err != nil {
				log.Printf("❌ Failed to insert drug: %v", err)
			} else {
				log.Printf("✅ Successfully inserted drug: %s", drug.Name)
			}
		} else {
			log.Printf("🔗 Looking for links on: %s", url)
			linkCount := 0
			e.DOM.Find("a[href]").Each(func(i int, s *goquery.Selection) {
				href, exists := s.Attr("href")
				if exists {
					// Convert relative URLs to absolute URLs
					if strings.HasPrefix(href, "/") {
						href = c.config.BaseURL + href
					} else if strings.HasPrefix(href, "./") {
						href = c.config.BaseURL + href[1:]
					} else if !strings.HasPrefix(href, "http") {
						href = c.config.BaseURL + "/" + href
					}

					// Improved URL filtering - check if it's a drug page
					if strings.Contains(href, "/thuoc/") && !strings.Contains(href, "#") {
						log.Printf("🔗 Found drug link: %s", href)
						linkCount++
						e.Request.Visit(href)
					}
				}
			})
			log.Printf("📊 Found %d drug links on %s", linkCount, url)
		}
	})

	log.Printf("🎯 Starting crawl from: %s", c.config.BaseURL)
	collector.Visit(c.config.BaseURL)
	collector.Wait()
	log.Println("🏁 Crawler finished!")
}
