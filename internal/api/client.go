package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/MikelMelnichuk/mycal/internal/models"
)

type APIError struct {
	Error   string
	Details string
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GetDayEvents(targetDate string, all bool, after string) ([]models.Event, error) {

	// Get the endpoint
	endpoint := c.BaseURL + "api/v1/events/day"
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Invalid base url: %w", err)
	}

	// Construct the query
	q := u.Query()
	q.Set("target_date_str", targetDate)
	q.Set("all", strconv.FormatBool(all))
	if after != "" {
		q.Set("after", after)
	}
	u.RawQuery = q.Encode()

	// Create a HTTP request
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Handle errors
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		var apiErr APIError
		if decodeErr := json.NewDecoder(resp.Body).Decode(&apiErr); decodeErr != nil {
			return nil, fmt.Errorf(
				"API returned status %d, failed to parse error response. Text is: %q",
				resp.StatusCode,
				resp.Body,
			)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Error, apiErr.Details)
	}

	// Parse successful JSON array
	var events []models.Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return events, nil
}

func (c *Client) GetWeekEvents(next bool, all bool) ([]models.Event, error) {

	// Define the endpoint to be used
	endpoint := c.BaseURL + "api/v1/events/week"
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Invalid base url: %w", err)
	}

	// Add the expected variables
	q := u.Query()
	q.Set("next", strconv.FormatBool(next))
	q.Set("all", strconv.FormatBool(all))
	u.RawQuery = q.Encode()

	// Send the request to the backend server
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Process the received response
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}

	// Handle error that come from the server as a response
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		var apiErr APIError
		if decodeErr := json.NewDecoder(resp.Body).Decode(&apiErr); decodeErr != nil {
			return nil, fmt.Errorf(
				"API returned status %d, failed to parse error response. Text is: %q",
				resp.StatusCode,
				resp.Body,
			)
		}
		return nil, fmt.Errorf("API error: %s - %s", apiErr.Error, apiErr.Details)
	}

	var weekEvents []models.Event
	if err := json.NewDecoder(resp.Body).Decode(&weekEvents); err != nil {
		return nil, fmt.Errorf("failed to decode events: %w", err)
	}

	return weekEvents, nil
}

func (c *Client) GetEventByID(eventID string) (*models.Event, error) {
	return nil, nil
}

func (c *Client) HealthCheckBackend() error {
	// Define the endpoint to be used
	endpoint := c.BaseURL
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("Invalid base url: %w", err)
	}

	// Send the request to the backend server
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Process the received response
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	// Handle error that come from the server as a response
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, resp.Body)
	}

	fmt.Println("Communication with backend was established")

	return nil
}

func (c *Client) HealthCheckDB() error {
	// Define the endpoint to be used
	endpoint := c.BaseURL + "api/v1/health"
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("Invalid base url: %w", err)
	}

	// Send the request to the backend server
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Process the received response
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	// Handle error that come from the server as a response
	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		return fmt.Errorf("API error: %d - %s", resp.StatusCode, resp.Body)
	}

	fmt.Println("Communication with backend and DB was established")

	return nil
}

// CreateEventRequest defines the expected JSON payload.
// Description and Location are optional (omitted when empty).
type CreateEventRequest struct {
	Title       string `json:"title"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}

// CreateEvent sends a POST request to create a new event.
// It returns (true, nil) if the server responds with 200 and is_successful: true.
// If the server returns a non-200 status, or the response cannot be parsed, an error is returned.
func (c *Client) CreateEvent(title, start, end, description, location string) (bool, error) {
	// Build the endpoint URL
	endpoint := c.BaseURL + "api/v1/create"
	u, err := url.Parse(endpoint)
	if err != nil {
		return false, fmt.Errorf("invalid base URL: %w", err)
	}

	// Marshal the request body
	reqBody := CreateEventRequest{
		Title:       title,
		Start:       start,
		End:         end,
		Description: description,
		Location:    location,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code (200 is expected, but we also accept 2xx)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false, fmt.Errorf("API error: %d - %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var respData struct {
		IsSuccessful bool `json:"is_successful"`
	}
	if err := json.Unmarshal(body, &respData); err != nil {
		return false, fmt.Errorf("failed to parse response: %w", err)
	}

	// Optionally log success
	fmt.Println("Event creation request processed")

	return respData.IsSuccessful, nil
}
