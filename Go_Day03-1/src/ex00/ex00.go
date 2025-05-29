package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", 
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Ошибка при создании клиента: %s", err)
	}

	createIndexWithMapping(es)

	restaurants, err := readCSV("../../materials/data.csv")
	if err != nil {
		log.Fatalf("Ошибка при чтении CSV: %s", err)
	}

	loadData(es, restaurants)
}

func createIndexWithMapping(es *elasticsearch.Client) {
	mapping := `{
		"mappings": {
			"properties": {
				"name": {
					"type": "text"
				},
				"address": {
					"type": "text"
				},
				"phone": {
					"type": "text"
				},
				"location": {
					"type": "geo_point"
				}
			}
		}
	}`

	req := esapi.IndicesCreateRequest{
		Index: "places",
		Body:  strings.NewReader(mapping),
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Ошибка при создании индекса: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Ошибка при создании индекса: %s", res.String())
	}

	fmt.Println("Индекс 'places' создан с маппингом")
}

func readCSV(filename string) ([]map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t' 
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении CSV: %s", err)
	}

	var restaurants []map[string]interface{}
	for i, record := range records {
		if i == 0 {
			continue
		}

		lon, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при парсинге долготы: %s", err)
		}

		lat, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка при парсинге широты: %s", err)
		}

		restaurant := map[string]interface{}{
			"id": 		record[0],
			"name":     record[1],
			"address":  record[2],
			"phone":    record[3],
			"location": map[string]float64{"lat": lat, "lon": lon},
		}

		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

func loadData(es *elasticsearch.Client, restaurants []map[string]interface{}) {
	var buf bytes.Buffer
	for i, restaurant := range restaurants {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "places",
				"_id":    i + 1,
			},
		}
		metaJSON, _ := json.Marshal(meta)
		buf.Write(metaJSON)
		buf.WriteString("\n")

		docJSON, _ := json.Marshal(restaurant)
		buf.Write(docJSON)
		buf.WriteString("\n")
	}

	res, err := es.Bulk(
		strings.NewReader(buf.String()),
		es.Bulk.WithIndex("places"),
		es.Bulk.WithContext(context.Background()),
	)
	if err != nil {
		log.Fatalf("Ошибка при отправке пакета операций: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Ошибка при выполнении Bulk API: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		log.Fatalf("Ошибка при декодировании ответа: %s", err)
	}

	fmt.Println("Данные загружены в индекс 'places'")
}
