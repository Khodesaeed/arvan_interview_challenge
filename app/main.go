package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/khodesaeed/arvan_interview_challenge/internal/api"
	"github.com/khodesaeed/arvan_interview_challenge/internal/db"
	_ "github.com/khodesaeed/arvan_interview_challenge/internal/metrics"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()
	dbPool, err := db.CreateDatabaseAndTable(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbPool.Close()

	token := os.Getenv("IPINFO_TOKEN")
	if token == "" {
		log.Println("IPINFO_TOKEN environment variable not set. API calls may fail.")
	}
	ipinfoClient := ipinfo.NewClient(nil, nil, token)

	apiHandlers := api.NewAPI(dbPool, ipinfoClient)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":5050", nil))
	}()

	http.HandleFunc("/get_country", apiHandlers.GetCountryHandler)
	http.HandleFunc("/live", apiHandlers.LiveHandler)
	http.HandleFunc("/ready", apiHandlers.ReadyHandler)
	log.Println("Starting API server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
