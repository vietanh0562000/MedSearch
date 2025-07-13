package crawler

import (
	"MedSearch/app/config"
	"MedSearch/app/database/repository"
	"MedSearch/app/logger"
	"fmt"
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

const maxResultCount = 100

func (c *Crawler) Start() {
	c.logger.Log("🚀 Starting crawler...")

	//TODO: Create collector by colly
	collector := colly.NewCollector(
		colly.MaxDepth(10), // Increased depth
	)

	// Rotate user agents to avoid detection
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/121.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	collector.OnRequest(func(r *colly.Request) {
		// Randomly select a user agent
		userAgent := userAgents[time.Now().UnixNano()%int64(len(userAgents))]
		r.Headers.Set("User-Agent", userAgent)
	})

	collector.SetRequestTimeout(30 * time.Second)

	// Optional: Add proxy support (uncomment if you have proxies)
	// collector.SetProxy("http://your-proxy-here:port")

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,               // Single request at a time
		Delay:       5 * time.Second, // Much longer delay
		RandomDelay: 3 * time.Second, // More random delay
	})

	// Add error handling
	collector.OnError(func(r *colly.Response, err error) {
		c.logger.Log("❌ Error visiting %s: %v", r.Request.URL, err)

		// Simple retry for 403/429 errors (rate limiting)
		if r.StatusCode == 403 || r.StatusCode == 429 {
			c.logger.Log("🔄 Rate limited, waiting 30 seconds before retry...")
			time.Sleep(30 * time.Second)
			r.Request.Retry()
		}
	})

	//TODO: Set rate limit for collector: OnRequest, OnResponse. OnHTML, OnError
	// Get Info in OnHTML function
	collector.OnRequest(func(r *colly.Request) {
		c.logger.Log("🌐 Visiting: %s", r.URL.String())

		// Add more realistic browser headers
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9,vi;q=0.8")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		r.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
		r.Headers.Set("Sec-Ch-Ua-Platform", `"macOS"`)
		r.Headers.Set("DNT", "1")
	})

	// Add response callback to see what we're getting
	collector.OnResponse(func(r *colly.Response) {
		c.logger.Log("📥 Response received for: %s (Status: %d)", r.Request.URL, r.StatusCode)
		c.logger.Log("📥 Response length: %d bytes", len(r.Body))
		if len(r.Body) < 1000 {
			c.logger.Log("📥 Response body: %s", string(r.Body))
		}
	})

	//TODO: Set callback for collector
	collector.OnHTML("*", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		c.logger.Log("📄 Processing: %s", url)

		// // Print the full HTML content for debugging
		// html, _ := e.DOM.Html()
		// fmt.Println("================ HTML START ================")
		// fmt.Printf("URL: %s\n", url)
		// fmt.Println(html)
		// fmt.Println("================ HTML END ================")

		if strings.Contains(url, ".html") {
			c.logger.Log("💊 Found drug page: %s", url)
			c.extractDrugHtml(e)
		} else {
			e.DOM.Find("a[href]").Each(func(i int, s *goquery.Selection) {
				fmt.Println("Found a lin")
				href, exists := s.Attr("href")
				if exists {
					c.visit(href, e)

					skipCount := 0
					slugs := GetDrugLink(href[1:], skipCount)
					for len(slugs) == maxResultCount {
						// Visit each slug in the current batch
						for _, slug := range slugs {
							c.visit(slug, e)
						}
						skipCount += maxResultCount
						slugs = GetDrugLink(href, skipCount)
					}
					// Visit remaining slugs in the last batch
					for _, slug := range slugs {
						c.visit(slug, e)
					}
				}
			})
		}
	})

	c.logger.Log("🎯 Starting crawl from: %s", c.config.StartURL)
	collector.Visit(c.config.StartURL)
	collector.Wait()
	log.Println("🏁 Crawler finished!")
}

func (c *Crawler) extractDrugHtml(e *colly.HTMLElement) {
	drug := ParseDrug(e)
	start := time.Now()
	err := repository.InsertDrug(&drug)
	elapsed := time.Since(start)
	if err != nil {
		c.logger.Log("❌ Failed to insert drug: %v (took %v)", err, elapsed)
	} else {
		c.logger.Log("✅ Successfully inserted drug: %s (took %v)", drug.Name, elapsed)
	}
}

func (c *Crawler) visit(href string, e *colly.HTMLElement) {
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
		e.Request.Visit(href)
	}
}
