package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/harsshhan/RA2311026010086/logging_middleware"
	"github.com/harsshhan/RA2311026010086/vehicle_scheduling/api"
	"github.com/harsshhan/RA2311026010086/vehicle_scheduling/scheduler"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		loggingmiddleware.Log("backend", "main", "Could not load .env file", "warn")
	}

	r := gin.Default()

	r.GET("/schedule", func(c *gin.Context) {
		baseURL := os.Getenv("API_URL")

		authToken := os.Getenv("AUTH_TOKEN")
		if authToken == "" {
			loggingmiddleware.Log("backend", "handler", "AUTH_TOKEN is missing", "warn")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing."})
			return
		}

		depots, err := api.FetchDepots(baseURL, authToken)
		if err != nil {
			loggingmiddleware.Log("backend", "handler", fmt.Sprintf("Error fetching depots: %v", err), "error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch depots", "details": err.Error()})
			return
		}

		totalHours := 0
		for _, d := range depots {
			totalHours += d.MechanicHours
		}
		loggingmiddleware.Log("backend", "handler", fmt.Sprintf("Total Available Mechanic Hours (Budget): %d", totalHours), "info")

		vehicles, err := api.FetchVehicles(baseURL, authToken)
		if err != nil {
			loggingmiddleware.Log("backend", "handler", fmt.Sprintf("Error fetching vehicles: %v", err), "error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch vehicles", "details": err.Error()})
			return
		}
		loggingmiddleware.Log("backend", "handler", fmt.Sprintf("Total Vehicles Requiring Maintenance: %d", len(vehicles)), "info")

		maxImpact, selectedTasks := scheduler.Optimize(totalHours, vehicles)

		resultMsg := fmt.Sprintf("Optimization Results: Max Impact: %d, Vehicles Selected: %d", maxImpact, len(selectedTasks))
		loggingmiddleware.Log("backend", "handler", resultMsg, "info")

		c.JSON(http.StatusOK, gin.H{
			"total_mechanic_hours": totalHours,
			"total_vehicles":       len(vehicles),
			"max_impact":           maxImpact,
			"vehicles_selected":    len(selectedTasks),
			"selected_task_ids":    selectedTasks,
		})
	})
	loggingmiddleware.Log("backend", "main", "Starting server on :8080", "info")
	r.Run(":8080")
}
