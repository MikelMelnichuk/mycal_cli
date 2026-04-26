package api

import (
	"net/http"

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

func (c *Client) GetDayEvent(targetDateStr string, all bool, after string) ([]models.Event, error) {
	return nil, nil
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
