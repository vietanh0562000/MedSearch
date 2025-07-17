package database

import (
	"MedSearch/app/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"go.mongodb.org/mongo-driver/bson"
)

func MigrateToElastic() {
	// --- Bắt đầu vòng lặp ---
	log.Println("📖 Bắt đầu đọc và ghi từng document...")
	collection := DB.Collection("drugs")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatalf("❌ Lỗi khi tìm kiếm document trong Mongo: %s", err)
	}
	defer cursor.Close(context.TODO())

	var totalCount int
	for cursor.Next(context.TODO()) {
		var drug models.Drug
		if err := cursor.Decode(&drug); err != nil {
			log.Printf("⚠️ Lỗi khi decode document: %s", err)
			continue
		}

		// Ánh xạ dữ liệu từ MongoDoc sang ElasticDoc
		drugLite := models.DrugLite{
			Name: drug.Name,
		}

		fmt.Println("ID Doc: ", drug.ID)
		// Ghi ngay lập tức document này vào Elasticsearch
		err := indexSingleDoc(ElasticClient, drug.ID.Hex(), drugLite)
		if err != nil {
			log.Printf("🔥 Lỗi khi ghi document ID %s: %s", drug.ID.Hex(), err)
			continue
		}

		totalCount++
		// Log tiến trình mỗi 100 document để tránh làm chậm console
		if totalCount%100 == 0 {
			log.Printf("📝 Đã xử lý %d document...", totalCount)
		}
	}

	if err := cursor.Err(); err != nil {
		log.Fatalf("❌ Lỗi con trỏ MongoDB: %s", err)
	}

	log.Printf("🎉 Hoàn tất! Đã di chuyển tổng cộng %d document.", totalCount)
	//log.Printf("Tổng thời gian: %s", time.Since(startTime))
}

// indexSingleDoc thực hiện ghi một document duy nhất vào Elasticsearch
func indexSingleDoc(esClient *elasticsearch.Client, docID string, doc models.DrugLite) error {
	// Chuyển đổi struct ElasticDoc thành JSON
	payload, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	// Tạo yêu cầu ghi
	req := esapi.IndexRequest{
		Index:      "drugs",
		DocumentID: docID,
		Body:       bytes.NewReader(payload),
		Refresh:    "false", // "false" để tăng tốc độ ghi
	}

	// Thực hiện yêu cầu
	res, err := req.Do(context.Background(), esClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Kiểm tra lỗi từ response của Elasticsearch
	if res.IsError() {
		log.Printf("Response Body: %s", res.String())
		return err
	}

	return nil
}
