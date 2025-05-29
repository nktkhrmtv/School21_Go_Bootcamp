package render

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type Place struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

type renderData struct {
	Places   []Place
	Total    int
	Page     int
	LastPage int
}

type Store interface {
	GetPlaces(limit int, offset int) ([]Place, int, error)
}

type ElasticStore struct {
	Client *elasticsearch.Client
	Index  string
	JwtKey []byte
}

func GetClient() (*elasticsearch.Client, error) {
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return client, err
}

func (es *ElasticStore) HandlePlaces(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, fmt.Sprintf("Invalid 'page' value: '%s'", pageStr), http.StatusBadRequest)
		return
	}

	limit := 10
	offset := (page - 1) * limit

	places, total, err := es.GetPlaces(limit, offset)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при получении данных: %s", err), http.StatusInternalServerError)
		return
	}

	lastPage := (total + limit - 1) / limit
	if page > lastPage {
		http.Error(w, fmt.Sprintf("Invalid 'page' value: '%d'", page), http.StatusBadRequest)
		return
	}

	data := renderData{
		Places:   places,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
	}

	RenderFuncHTML(data, &w)
}

func (es *ElasticStore) GetPlaces(limit int, offset int) ([]Place, int, error) {
	query := map[string]interface{}{
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		"track_total_hits": true,
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при сериализации запроса: %s", err)
	}

	req := esapi.SearchRequest{
		Index: []string{es.Index},
		Body:  strings.NewReader(string(queryJSON)),
	}

	res, err := req.Do(context.Background(), es.Client)
	if err != nil {
		return nil, 0, fmt.Errorf("ошибка при выполнении запроса: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("ошибка ответа от Elasticsearch: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("ошибка при декодировании ответа: %s", err)
	}

	hits := result["hits"].(map[string]interface{})
	total := int(hits["total"].(map[string]interface{})["value"].(float64))

	var places []Place
	for _, hit := range hits["hits"].([]interface{}) {
		var place Place
		source := hit.(map[string]interface{})["_source"]
		sourceJSON, _ := json.Marshal(source)
		if err := json.Unmarshal(sourceJSON, &place); err != nil {
			return nil, 0, fmt.Errorf("ошибка при декодировании документа: %s", err)
		}
		places = append(places, place)
	}

	return places, total, nil
}

func RenderFuncHTML(data renderData, w *http.ResponseWriter) {
	tmpl := template.New("places.html").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	})

	tmpl, err := tmpl.Parse(` 
    <!doctype html>
    <html>
    <head>
        <meta charset="utf-8">
        <title>Places</title>
    </head>
    <body>
    <h5>Total: {{.Total}}</h5>
    <ul>
        {{range .Places}}
        <li>
            <div>{{.Name}}</div>
            <div>{{.Address}}</div>
            <div>{{.Phone}}</div>
        </li>
        {{end}}
    </ul>
    {{if gt .Page 1}}<a href="/?page={{sub .Page 1}}">Previous</a>{{end}}
    {{if lt .Page .LastPage}}<a href="/?page={{add .Page 1}}">Next</a>{{end}}
    <a href="/?page={{.LastPage}}">Last</a>
    </body>
    </html>
    `)
	if err != nil {
		http.Error(*w, fmt.Sprintf("Ошибка при парсинге шаблона: %s", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(*w, data); err != nil {
		http.Error(*w, fmt.Sprintf("Ошибка при рендеринге шаблона: %s", err), http.StatusInternalServerError)
	}
}
