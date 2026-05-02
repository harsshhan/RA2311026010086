package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	loggingmiddleware "github.com/harsshhan/RA2311026010086/logging_middleware"
	"github.com/harsshhan/RA2311026010086/notification_app_be/models"
)

func FetchNotifications(baseURL, token string) ([]models.Notification, error) {
	endpoint := baseURL + "evaluation-service/notifications"
	loggingmiddleware.Log("backend", "api", "Fetching notifications from: "+endpoint, "info")

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		loggingmiddleware.Log("backend", "api", "Failed to create request: "+err.Error(), "error")
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		loggingmiddleware.Log("backend", "api", "Request failed: "+err.Error(), "error")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		errMsg := "Unexpected status: " + resp.Status + " body: " + string(bodyBytes)
		loggingmiddleware.Log("backend", "api", errMsg, "error")
		return nil, fmt.Errorf(errMsg)
	}

	var data models.NotificationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		loggingmiddleware.Log("backend", "api", "Failed to decode response: "+err.Error(), "error")
		return nil, err
	}

	loggingmiddleware.Log("backend", "api", "Fetched notifications successfully", "info")
	return data.Notifications, nil
}
