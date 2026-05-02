package loggingmiddleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/harsshhan/RA2311026010086/logging_middleware/models"
)


func Log(stack,pkg,message,level string) {
	logData := models.LogRequest{
		Stack: stack,
		Package: pkg,
		Message: message,
		Level: level,
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		fmt.Println("Error marshalling log data:", err)
		return
	}

	fmt.Println("Logs:",logData)
	
	apiUrl := os.Getenv("API_URL")
	endpoint := apiUrl + "evaluation-service/logs"
	
	resp, err := http.Post(endpoint,"application/json",bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending log:", err)
		return
	}
	defer resp.Body.Close()}



