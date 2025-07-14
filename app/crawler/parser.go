package crawler

import (
	"MedSearch/app/models"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func ParseDrug(e *colly.HTMLElement) models.Drug {
	var drug models.Drug

	drug.Link = e.Request.URL.String()

	parseGeneralInfo(e, &drug)
	parseSideEffects(e, &drug)
	parseStorage(e, &drug)
	parseDrugUses(e, &drug)
	parseDosage(e, &drug)

	return drug
}

func parseGeneralInfo(e *colly.HTMLElement, drug *models.Drug) {
	drug.Name = e.ChildText("h1[data-test='product_name']")
	drug.Price = e.DOM.Find(`span[data-test='price']`).First().Text()
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
}

func parseDrugUses(e *colly.HTMLElement, drug *models.Drug) {
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
}

func parseDosage(e *colly.HTMLElement, drug *models.Drug) {
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
}

func parseStorage(e *colly.HTMLElement, drug *models.Drug) {
	e.DOM.Find("div.preservation").Each(func(i int, s *goquery.Selection) {
		s.Find("p").Each(func(_ int, p *goquery.Selection) {
			drug.Storage += p.Text()
		})
	})
}

func parseSideEffects(e *colly.HTMLElement, drug *models.Drug) {
	e.DOM.Find("div.adverseEffect").Each(func(i int, s *goquery.Selection) {
		s.Find("p").Each(func(_ int, p *goquery.Selection) {
			drug.SideEffects += p.Text()
		})
	})
}
