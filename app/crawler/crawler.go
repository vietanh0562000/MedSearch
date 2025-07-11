package crawler

import (
	"MedSearch/app/config"
	"MedSearch/app/database/repository"
	"MedSearch/app/logger"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Crawler struct {
	config *config.Config
	logger *logger.MLogger
}

func NewCrawler(config *config.Config, logger *logger.MLogger) *Crawler {
	newCrawler := Crawler{
		config: config,
		logger: logger,
	}

	return &newCrawler
}

func (c *Crawler) Start() {
	c.logger.Log("🚀 Starting crawler...")

	//TODO: Create collector by colly
	collector := colly.NewCollector(
		colly.MaxDepth(10), // Increased depth
	)

	collector.SetRequestTimeout(30 * time.Second)

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 4,
		Delay:       1 * time.Second, // Reduced delay
		RandomDelay: 1 * time.Second,
	})

	// Add error handling
	collector.OnError(func(r *colly.Response, err error) {
		c.logger.Log("❌ Error visiting %s: %v", r.Request.URL, err)
	})

	//TODO: Set rate limit for collector: OnRequest, OnResponse. OnHTML, OnError
	// Get Info in OnHTML function
	collector.OnRequest(func(r *colly.Request) {
		c.logger.Log("🌐 Visiting: %s", r.URL.String())
	})

	//TODO: Set callback for collector
	collector.OnHTML("body", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		c.logger.Log("📄 Processing: %s", url)

		if strings.Contains(url, "/thuoc/") && strings.Contains(url, ".html") {
			c.logger.Log("💊 Found drug page: %s", url)
			drug := ParseDrug(e)
			start := time.Now()
			err := repository.InsertDrug(&drug)
			elapsed := time.Since(start)
			if err != nil {
				c.logger.Log("❌ Failed to insert drug: %v (took %v)", err, elapsed)
			} else {
				c.logger.Log("✅ Successfully inserted drug: %s (took %v)", drug.Name, elapsed)
			}
		} else {
			c.logger.Log("🔗 Looking for links on: %s", url)
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
					if strings.Contains(href, "/thuoc/") {
						c.logger.Log("🔗 Found drug link: %s", href)
						linkCount++
						e.Request.Visit(href)
					}
				}
			})
			c.logger.Log("📊 Found %d drug links on %s", linkCount, url)
		}
	})

	c.logger.Log("🎯 Starting crawl from: %s", c.config.StartURL)
	collector.Visit(c.config.StartURL)
	collector.Wait()
	log.Println("🏁 Crawler finished!")
}
