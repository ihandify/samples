package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Global configuration variables
var (
	APIURL string
	APIKey string
)

// InitConfig initializes the configuration from environment variables
func InitConfig() {
	APIURL = os.Getenv("API_URL")
	APIKey = os.Getenv("API_KEY")

	if APIURL == "" {
		log.Println("⚠️ Warning: API_URL is not set in .env")
	}
	if APIKey == "" {
		log.Println("⚠️ Warning: API_KEY is not set in .env")
	}
}

// Request payload matching Python structure
type ScopedKeyRequest struct {
    // Keep slice as-is (slices can be nil), but add omitempty
    Engines          []string `json:"engines,omitempty"`
    // Use a pointer to an int so it can be nil instead of defaulting to 0
    ExpiresInSeconds *int     `json:"expiresInSeconds,omitempty"`
}

// GenerateScopedPublicKey calls the upstream API
func GenerateScopedPublicKey(engines []string, expiresInSeconds *int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/plan/auth/generate-scoped-public-key", APIURL)

	payloadBytes, err := json.Marshal(ScopedKeyRequest{
        Engines:          engines,
        ExpiresInSeconds: expiresInSeconds,
    })
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", APIKey)

	// 30-second timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("No response received from server. Please check your URL or ensure the backend server is running.")
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Generate Scoped Public Key API Request Failed! Status Code: %d\n", resp.StatusCode)
		
		if resp.StatusCode == http.StatusUnauthorized {
			log.Println("💡 Hint: API Key is missing or invalid.")
		} else if resp.StatusCode == http.StatusForbidden {
			log.Println("💡 Hint: Your IP address is not included in the IP allowlist.")
		}
		return nil, fmt.Errorf("upstream error status %d", resp.StatusCode)
	}

	log.Println("✅ Generate Scoped Public Key API Request Successful!")

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
