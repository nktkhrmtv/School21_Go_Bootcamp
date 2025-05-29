package render

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    Name  string `json:"name"`
    Admin bool   `json:"admin"`
    jwt.RegisteredClaims
}

func (es *ElasticStore) HandleGetToken(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Name:  "Nikita",
        Admin: true,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(es.JwtKey)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(map[string]string{
            "error": "Ошибка при создании токена",
        })
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "token": tokenString,
    })
}

func (es *ElasticStore) JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Требуется авторизация",
            })
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        if tokenString == authHeader {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Неверный формат токена",
            })
            return
        }

        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return es.JwtKey, nil
        })
        if err != nil || !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(map[string]string{
                "error": "Неверный токен",
            })
            return
        }

        next(w, r)
    }
}