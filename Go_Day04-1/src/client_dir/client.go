package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Order struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

func main() {
	candyType := flag.String("k", "", "Type of candy (e.g., AA)")
	candyCount := flag.Int("c", 0, "Number of candies")
	money := flag.Int("m", 0, "Amount of money")
	flag.Parse()

	if *candyType == "" || *candyCount <= 0 || *money <= 0 {
		log.Fatalf("All flags (-k, -c, -m) are required and must be positive")
	}

	order := Order{
		Money:      *money,
		CandyType:  *candyType,
		CandyCount: *candyCount,
	}

	payload, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order: %v", err)
	}

	cert, err := tls.LoadX509KeyPair("clientcandy.tld/cert.pem", "clientcandy.tld/key.pem")
	if err != nil {
		log.Fatalf("Failed to load client certificate: %v", err)
	}

	caCert, err := os.ReadFile("minica.pem")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		ServerName:   "candy.tld",
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	resp, err := client.Post("https://candy.tld:39791/buy_candy", "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	formattedBody := strings.ReplaceAll(string(body), "\\n", "\n")

	fmt.Println("Response:", formattedBody)
}
