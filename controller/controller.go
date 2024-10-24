package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	uniqueRequests     = make(map[int]struct{})
	uniqueRequestCount int
	mutex              sync.Mutex
)

// handleAccept handles the API requests to /api/verve/accept
func HandleVerveAccept(c *gin.Context) {
	// Parse the mandatory "id" query parameter
	idParam := c.Query("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"message":    "id is required",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"message":    "Invalid id",
		})
		return
	}

	// Check for the optional "endpoint" query parameter
	endpoint := c.Query("endpoint")

	log.Info().Int("id", id).Str("endpoint", endpoint).Msg("Received a http endpoint")

	// Handle uniqueness
	mutex.Lock()
	if _, exists := uniqueRequests[id]; !exists {
		uniqueRequests[id] = struct{}{}
		uniqueRequestCount++
		log.Info().Int("unique request count: ", uniqueRequestCount).Msg("Unique request count increases")
	}
	mutex.Unlock()

	// If endpoint is provided, send the current unique count as a query parameter
	if endpoint != "" {
		go sendRequestToEndpoint(endpoint)
		// go makePostRequest(endpoint) // uncomment in case of testing
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"message":    "OK, Done!",
	})
}

// sendRequestToEndpoint sends the unique request count to the provided endpoint
func sendRequestToEndpoint(endpoint string) {
	mutex.Lock()
	count := uniqueRequestCount
	mutex.Unlock()

	// Send HTTP GET request to the endpoint
	resp, err := http.Get(fmt.Sprintf("%s?count=%d", endpoint, count))
	if err != nil {
		// Handle error (optional logging)
		log.Err(err).Str("endpoint", endpoint).Msg("Failed to send request to endpoint")
		return
	}
	defer resp.Body.Close()

	log.Info().Str("endpoint", endpoint).Msgf("Request processed successfully: Status code: %s", resp.Status)
}

// makePostRequest fires a POST request to the provided endpoint with the unique request count
func makePostRequest(endpoint string) {
	// Prepare the data structure for the POST request
	data := map[string]interface{}{
		"unique_request_count": uniqueRequestCount,
		"message":              "Unique requests logged",
	}

	// Convert the data to JSON
	body, err := json.Marshal(data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal JSON")
		return
	}

	// Create a new POST request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create POST request")
		return
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")

	// Send the POST request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Str("url", endpoint).Msg("Failed to make POST request")
		return
	}
	defer resp.Body.Close()

	// Log the HTTP status code of the response
	log.Info().Str("url", endpoint).Int("status_code", resp.StatusCode).Msg("POST request status")
}

func LogUniqueRequestsEveryMinute() {
	log.Info().Msg("Inside log unique request every minute")

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mutex.Lock()
		count := uniqueRequestCount
		uniqueRequestCount = 0
		uniqueRequests = make(map[int]struct{}) // reset for the new minute
		mutex.Unlock()

		// Log the unique request count
		log.Info().Int("unique_request_count", count).Msg("Logged unique requests")
	}
}
