package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type SearchController interface {
	Search(w http.ResponseWriter, r *http.Request)
}

type searchController struct {
	client          *http.Client
	homeserverURL   string
}

func NewSearchController() SearchController {
	homeserverURL := os.Getenv("HOMESERVER_URL")
	if homeserverURL == "" {
		homeserverURL = "http://homeserver:4000"
	}

	return &searchController{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		homeserverURL: homeserverURL,
	}
}

// Search proxies search requests to the homeserver
func (c *searchController) Search(w http.ResponseWriter, r *http.Request) {
	// Get query parameter
	query := strings.TrimSpace(r.URL.Query().Get("text"))
	if query == "" {
		http.Error(w, "Missing 'text' parameter", http.StatusBadRequest)
		return
	}

	// Build homeserver URL with all query parameters
	homeserverURL := fmt.Sprintf("%s/v1/search", c.homeserverURL)
	
	// Forward all query parameters to homeserver
	if len(r.URL.RawQuery) > 0 {
		homeserverURL += "?" + r.URL.RawQuery
	}

	// Make request to homeserver
	resp, err := c.client.Get(homeserverURL)
	if err != nil {
		log.Error("Failed to query homeserver", "url", homeserverURL, "error", err)
		http.Error(w, "Search service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Forward status code
	w.WriteHeader(resp.StatusCode)

	// Forward response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Forward response body
	var response interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Error("Failed to decode homeserver response", "error", err)
		http.Error(w, "Invalid response from search service", http.StatusBadGateway)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Error("Failed to encode response", "error", err)
		return
	}
}