package api

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/khodesaeed/arvan_interview_challenge/internal/db"
	"github.com/khodesaeed/arvan_interview_challenge/internal/metrics"

	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/jackc/pgx/v5/pgxpool"
)

// API holds the database connection and the IPInfo client
type API struct {
	DB     *pgxpool.Pool
	IPInfo *ipinfo.Client
}

// NewAPI creates a new API instance with initialized clients
func NewAPI(dbPool *pgxpool.Pool, ipinfoClient *ipinfo.Client) *API {
	return &API{
		DB:     dbPool,
		IPInfo: ipinfoClient,
	}
}

// GetCountryHandler handles the /get_country API endpoint
func (a *API) GetCountryHandler(w http.ResponseWriter, r *http.Request) {
	// Start timer for latency metric
	startTime := time.Now()
	metrics.InflightRequests.WithLabelValues("/get_country", r.Method).Inc()
	
	// Create a variable to hold the HTTP status code
	statusCode := http.StatusOK
	
	// Defer the metrics updates to run at the end of the function
	defer func() {
		metrics.InflightRequests.WithLabelValues("/get_country", r.Method).Dec()
		metrics.RequestLatency.WithLabelValues("/get_country", r.Method, strconv.Itoa(statusCode)).Observe(time.Since(startTime).Seconds())
		metrics.RequestTotal.WithLabelValues("/get_country", strconv.Itoa(statusCode)).Inc()
	}()

	// Get IP from query parameter
	ipStr := r.URL.Query().Get("ip")
	if ipStr == "" {
		statusCode = http.StatusBadRequest
		http.Error(w, `{"error": "IP parameter is missing"}`, http.StatusBadRequest)
		return
	}

	// Validate IP
	if !db.ValidateIP(ipStr) {
		statusCode = http.StatusBadRequest
		http.Error(w, `{"error": "Invalid IP address"}`, http.StatusBadRequest)
		return
	}

	// Check for country in the database
	country, err := db.GetCountryFromDB(r.Context(), a.DB, ipStr)
	if err == nil {
		// Found in cache, return it
		json.NewEncoder(w).Encode(map[string]string{"country": country})
		return
	}

	// Not in cache, fetch from ipinfo.io
	info, err := a.IPInfo.GetIPInfo(net.ParseIP(ipStr))
	if err != nil {
		log.Printf("Error fetching IP info: %v", err)
		statusCode = http.StatusInternalServerError
		http.Error(w, `{"error": "something happened in the external API"}`, http.StatusInternalServerError)
		return
	}

	country = info.Country

	// Save to database for caching
	go func() {
		err := db.SaveCountryToDB(context.Background(), a.DB, ipStr, country)
		if err != nil {
			log.Printf("Error saving IP to database: %v", err)
		}
	}()

	json.NewEncoder(w).Encode(map[string]string{"country": country})
}

// LiveHandler checks if the database connection is alive.
func (a *API) LiveHandler(w http.ResponseWriter, r *http.Request) {
	if err := a.DB.Ping(r.Context()); err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Live!"))
}

// ReadyHandler checks if the application is ready to serve requests.
func (a *API) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	if err := a.DB.Ping(r.Context()); err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Ready!"))
}