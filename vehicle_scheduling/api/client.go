package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/harsshhan/RA2311026010086/vehicle_scheduling/models"
)

func FetchDepots(baseURL, token string) ([]models.Depot, error) {
	req, err := http.NewRequest("GET", baseURL+"evaluation-service/depots", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var data models.DepotsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Depots, nil
}

func FetchVehicles(baseURL, token string) ([]models.Vehicle, error) {
	req, err := http.NewRequest("GET", baseURL+"evaluation-service/vehicles", nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	var data models.VehiclesResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Vehicles, nil
}
