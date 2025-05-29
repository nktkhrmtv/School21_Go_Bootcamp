package render

import (
	"encoding/json"
	"fmt"
    "strconv"
    "net/http"
)

type JSONResponse struct {
    Name     string  `json:"name"`
    Total    int     `json:"total"`
    Places   []Place `json:"places"`
    PrevPage int     `json:"prev_page,omitempty"`
    NextPage int     `json:"next_page,omitempty"`
    LastPage int     `json:"last_page"`
}

type JSONError struct {
	Error    string  `json:"error,omitempty"`
}

func (es *ElasticStore) HandlePlacesAPI(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    pageStr := r.URL.Query().Get("page")
    if pageStr == "" {
        pageStr = "1"
    }
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(JSONError{
            Error: fmt.Sprintf("Invalid 'page' value: '%s'", pageStr),
        })
        return
    }

    limit := 10
    offset := (page - 1) * limit
    places, total, err := es.GetPlaces(limit, offset)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(JSONError{
            Error: fmt.Sprintf("Ошибка при получении данных: %s", err),
        })
        return
    }

    lastPage := (total + limit - 1) / limit
    if page > lastPage {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(JSONError{
            Error: fmt.Sprintf("Invalid 'page' value: '%d'", page),
        })
        return
    }

    response := JSONResponse{
        Name:     "Places",
        Total:    total,
        Places:   places,
        LastPage: lastPage,
    }
    if page > 1 {
        response.PrevPage = page - 1
    }
    if page < lastPage {
        response.NextPage = page + 1
    }

    json.NewEncoder(w).Encode(response)
}