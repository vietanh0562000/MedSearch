package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CrawlExpand() {
	url := "https://api.nhathuoclongchau.com.vn/lccus/search-product-service/api/products/ecom/product/search/cate"

	payload := map[string]interface{}{
		"skipCount":      12,
		"maxResultCount": 8,
		"codes": []string{
			"productTypes", "priceRanges", "prescription", "objectUse", "skin", "flavor",
			"manufactor", "indications", "brand", "brandOrigin", "producer", "specification", "ingredient",
		},
		"sortType": 4,
		"category": []string{"thuoc/thuoc-di-ung"},
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
	fmt.Println(string(data))
}
