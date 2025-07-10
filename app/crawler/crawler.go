package crawler

import (
	"MedSearch/app/database/repository"
	"MedSearch/app/models"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
		var drug models.Drug
		drug.Name = e.ChildText("h1[data-test='product_name']")
		e.DOM.Find("tr.content-container").Each(func(i int, s *goquery.Selection) {
			label := strings.TrimSpace(s.Find("td:first-child").Text())
			value := strings.TrimSpace(s.Find("td:nth-child(2)").Text())

			switch label {
			case "Danh mục":
				drug.Category = value
			case "Số đăng ký":
				drug.RegistedNumber = value
			case "Dạng bào chế":
				drug.DosageForm = value
			case "Quy cách":
				drug.Packaging = value
			case "Thành phần":
				drug.Ingredients = strings.Split(value, ", ")
			case "Chỉ định":
				drug.Indication = value
			case "Chống chỉ định":
				drug.Contraindication = value
			case "Nhà sản xuất":
				drug.Manufacturer = value
			case "Nước sản xuất":
				drug.MAH = value
			case "Xuất xứ thương hiệu":
				drug.MAH = value
			case "Mô tả ngắn":
				drug.Description = value
			case "Lưu ý":
				drug.Notes = value
			}
		})
		drug.Price = e.DOM.Find(`span[data-test='price']`).First().Text()
		fmt.Println("Price:", drug.Price)

		e.DOM.Find("div.usage").Each(func(i int, s *goquery.Selection) {
			s.Find("h3").Each(func(_ int, h3 *goquery.Selection) {
				if strings.Contains(h3.Text(), "Chỉ định") {
					h3.NextUntil("h3").Each(func(_ int, p *goquery.Selection) {
						drug.Uses += p.Text()
					})
				}
			})
		})

		if drug.Uses == "" {
			e.DOM.Find(`div.usage`).Each(func(i int, s *goquery.Selection) {
				s.Find("h3").Each(func(_ int, h3 *goquery.Selection) {
					if h3.Text() == "Chỉ định" {
						// Tìm thẻ <ul> ngay sau h3
						ul := h3.NextFiltered("p").NextFiltered("ul")
						ul.Find("li").Each(func(_ int, li *goquery.Selection) {
							drug.Uses += li.Text() + "\n"
						})
					}
				})
			})
		}

		e.DOM.Find("div.dosage").Each(func(i int, s *goquery.Selection) {
			s.Find("h3").Each(func(_ int, h3 *goquery.Selection) {
				if strings.Contains(h3.Text(), "Cách dùng") {
					h3.NextUntil("h3").Each(func(_ int, p *goquery.Selection) {
						drug.Administration += p.Text()
					})
				}
				if strings.Contains(h3.Text(), "Liều dùng") {
					h3.NextUntil("h3").Each(func(_ int, p *goquery.Selection) {
						drug.Dosage += p.Text()
					})
				}
			})
		})

		e.DOM.Find("div.adverseEffect").Each(func(i int, s *goquery.Selection) {
			s.Find("p").Each(func(_ int, p *goquery.Selection) {
				drug.SideEffects += p.Text()
			})
		})

		e.DOM.Find("div.preservation").Each(func(i int, s *goquery.Selection) {
			s.Find("p").Each(func(_ int, p *goquery.Selection) {
				drug.Storage += p.Text()
			})
		})

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
