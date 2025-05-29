package main

import (
    "log"
    "net/http"
    "day03/render"
)

func main() {
    client, err := render.GetClient()
    if err != nil {
        log.Fatalf("Ошибка при создании клиента Elasticsearch: %s", err)
    }

    store := render.ElasticStore{
        Client: client,
        Index:  "places",
        JwtKey: []byte("meteoriw_key"),
    }

    http.HandleFunc("/", store.HandlePlaces)
    http.HandleFunc("/api/places", store.HandlePlacesAPI)
    http.HandleFunc("/api/get_token", store.HandleGetToken)
    http.HandleFunc("/api/recommend", store.JWTMiddleware(store.HandleRecommendations))

    log.Println("Сервер запущен на http://127.0.0.1:8888")
    log.Fatal(http.ListenAndServe(":8888", nil))
}




