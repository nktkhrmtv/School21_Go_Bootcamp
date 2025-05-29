package render

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"context"
	"strings"
	"strconv"
    "net/http"
)

type RecommendationResponse struct {
    Name   string  `json:"name"`
    Places []Place `json:"places"`
}

func (es *ElasticStore) HandleRecommendations(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    latStr := r.URL.Query().Get("lat")
    lonStr := r.URL.Query().Get("lon")

    lat, err := strconv.ParseFloat(latStr, 64)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid 'lat' value",
        })
        return
    }

    lon, err := strconv.ParseFloat(lonStr, 64)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Invalid 'lon' value",
        })
        return
    }

    places, err := es.GetNearestPlaces(lat, lon)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": fmt.Sprintf("Ошибка при получении данных: %s", err),
        })
        return
    }

    response := RecommendationResponse{
        Name:   "Recommendation",
        Places: places,
    }

    json.NewEncoder(w).Encode(response)
}

func (es *ElasticStore) GetNearestPlaces(lat, lon float64) ([]Place, error) {
    query := map[string]interface{}{
        "size": 3,
        "sort": []map[string]interface{}{
            {
                "_geo_distance": map[string]interface{}{
                    "location": map[string]float64{
                        "lat": lat,
                        "lon": lon,
                    },
                    "order":         "asc",
                    "unit":         "km",
                    "mode":         "min",
                    "distance_type": "arc",
                    "ignore_unmapped": true,
                },
            },
        },
    }

    queryJSON, err := json.Marshal(query)
    if err != nil {
        return nil, fmt.Errorf("ошибка при сериализации запроса: %s", err)
    }

    req := esapi.SearchRequest{
        Index: []string{es.Index},
        Body:  strings.NewReader(string(queryJSON)),
    }

    res, err := req.Do(context.Background(), es.Client)
    if err != nil {
        return nil, fmt.Errorf("ошибка при выполнении запроса: %s", err)
    }
    defer res.Body.Close()

    if res.IsError() {
        return nil, fmt.Errorf("ошибка ответа от Elasticsearch: %s", res.String())
    }

    var result map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("ошибка при декодировании ответа: %s", err)
    }

    hits := result["hits"].(map[string]interface{})
    var places []Place
    for _, hit := range hits["hits"].([]interface{}) {
        var place Place
        source := hit.(map[string]interface{})["_source"]
        sourceJSON, _ := json.Marshal(source)
        if err := json.Unmarshal(sourceJSON, &place); err != nil {
            return nil, fmt.Errorf("ошибка при декодировании документа: %s", err)
        }
        places = append(places, place)
    }

    return places, nil
}