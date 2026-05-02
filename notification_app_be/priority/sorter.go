package priority

import (
	"sort"
	"time"

	loggingmiddleware "github.com/harsshhan/RA2311026010086/logging_middleware"
	"github.com/harsshhan/RA2311026010086/notification_app_be/models"
)

func typeWeight(notifType string) int {
	switch notifType {
	case "Placement":
		return 3
	case "Result":
		return 2
	case "Event":
		return 1
	default:
		return 0
	}
}

func GetTopN(notifications []models.Notification, n int) []models.Notification {
	loggingmiddleware.Log("backend", "priority", "Sorting notifications by priority", "info")

	sort.SliceStable(notifications, func(i, j int) bool {
		wi := typeWeight(notifications[i].Type)
		wj := typeWeight(notifications[j].Type)

		if wi != wj {
			return wi > wj
		}

		ti, errI := time.Parse("2006-01-02 15:04:05", notifications[i].Timestamp)
		tj, errJ := time.Parse("2006-01-02 15:04:05", notifications[j].Timestamp)

		if errI != nil || errJ != nil {
			return notifications[i].Timestamp > notifications[j].Timestamp
		}

		return ti.After(tj)
	})

	if n > len(notifications) {
		n = len(notifications)
	}

	loggingmiddleware.Log("backend", "priority", "Priority sorting complete", "info")
	return notifications[:n]
}
