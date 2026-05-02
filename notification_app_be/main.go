package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	loggingmiddleware "github.com/harsshhan/RA2311026010086/logging_middleware"
	"github.com/harsshhan/RA2311026010086/notification_app_be/api"
	"github.com/harsshhan/RA2311026010086/notification_app_be/priority"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		loggingmiddleware.Log("backend", "main", "Could not load .env file", "warn")
		fmt.Fprintln(os.Stderr, "error: failed to load .env file")
	}

	r := gin.Default()

	r.GET("/priority-inbox", func(c *gin.Context) {
		baseURL := os.Getenv("API_URL")
		authToken := os.Getenv("AUTH_TOKEN")

		if authToken == "" {
			loggingmiddleware.Log("backend", "handler", "AUTH_TOKEN is missing", "warn")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is missing."})
			return
		}

		n := 10
		if qn := c.Query("n"); qn != "" {
			if parsed, err := strconv.Atoi(qn); err == nil && parsed > 0 {
				n = parsed
			}
		}
		loggingmiddleware.Log("backend", "handler", "Requested top "+strconv.Itoa(n)+" priority notifications", "info")

		notifications, err := api.FetchNotifications(baseURL, authToken)
		if err != nil {
			loggingmiddleware.Log("backend", "handler", "Failed to fetch notifications: "+err.Error(), "error")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to fetch notifications",
				"details": err.Error(),
			})
			return
		}
		loggingmiddleware.Log("backend", "handler", "Total notifications received: "+strconv.Itoa(len(notifications)), "info")

		topN := priority.GetTopN(notifications, n)

		loggingmiddleware.Log("backend", "handler", "Returning top "+strconv.Itoa(len(topN))+" priority notifications", "info")

		c.JSON(http.StatusOK, gin.H{
			"total_notifications": len(notifications),
			"top_n":               len(topN),
			"priority_inbox":      topN,
		})
	})

	loggingmiddleware.Log("backend", "main", "Starting server on :8081", "info")
	r.Run(":8082")
}
