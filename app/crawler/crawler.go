package crawler

import (
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
				drug.Uses = value
			case "Lưu ý":
				drug.Notes = value
			}
		})
		fmt.Println(drug)
	})

	//TODO: Visit(url)
	url := "https://nhathuoclongchau.com.vn/thuoc/nolvadex-d-15679.html"
	collector.Visit(url)
	collector.Wait()
}
