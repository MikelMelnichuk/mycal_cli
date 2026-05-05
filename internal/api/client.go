package api

import (
	"encoding/json"
	"fmt"
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
	fmt.Printf("GetDayEvents got targetDate: %s, all: %t, after: %q\n", targetDate, all, after)

	// Get the endpoint
	endpoint := c.BaseURL + "/events/day"
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
	return nil, nil
}

func (c *Client) GetEventByID(eventID string) (*models.Event, error) {
	return nil, nil
}

func (c *Client) HealthCheck() error {
	return nil
}
