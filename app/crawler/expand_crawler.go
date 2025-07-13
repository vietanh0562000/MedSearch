package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetDrugLink(category string, skipCount int) []string {
	fmt.Println(category)
	url := "https://api.nhathuoclongchau.com.vn/lccus/search-product-service/api/products/ecom/product/search/cate"

	payload := map[string]interface{}{
		"skipCount":      skipCount,
		"maxResultCount": 100,
		"codes": []string{
			"productTypes", "priceRanges", "prescription", "objectUse", "skin", "flavor",
			"manufactor", "indications", "brand", "brandOrigin", "producer", "specification", "ingredient",
		},
		"sortType": 4,
		"category": []string{category},
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)

	var slugs []string

	// Parse the response as JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		fmt.Println("Raw response:", string(data))
		return slugs
	}

	// Extract slugs from product array

	if products, exists := jsonData["products"]; exists {
		if productArray, ok := products.([]interface{}); ok {
			for _, product := range productArray {
				if productMap, ok := product.(map[string]interface{}); ok {
					if slug, exists := productMap["slug"]; exists {
						if slugStr, ok := slug.(string); ok {
							slugs = append(slugs, slugStr)
						}
					}
				}
			}
		}
	}

	return slugs
}
