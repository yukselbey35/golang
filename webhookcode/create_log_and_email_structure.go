package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var alertData map[string]interface{} // Use a map for JSON parsing
	if err := json.NewDecoder(r.Body).Decode(&alertData); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode JSON data: %v", err), http.StatusBadRequest)
		return
	}

	// Log the received alert data in a human-readable format
	logHumanReadableAlert(alertData)

	// Check if the alert is critical and the namespace is in the predefined list
	severity, _ := alertData["severity"].(string) // Get severity as a string
	namespace, _ := alertData["details"].(map[string]interface{})["namespace"].(string)
	if severity == "critical" && isPredefinedNamespace(namespace) {
		// Get the owner's email for the namespace
		namespaceOwnerEmail := ownerEmails[namespace]
		if namespaceOwnerEmail != "" {
			// Send email to namespace owner
			if err := sendEmail(alertData, namespaceOwnerEmail); err != nil {
				log.Printf("Error sending email: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

// ... (rest of your code)

func logHumanReadableAlert(alertData map[string]interface{}) {
	alertname := alertData["alertname"]
	cluster := alertData["cluster"]
	message := alertData["details"].(map[string]interface{})["annotations_message"]
	severity := alertData["severity"]
	namespace := alertData["details"].(map[string]interface{})["namespace"]
	target := alertData["details"].(map[string]interface{})["target"]
	startTime := alertData["details"].(map[string]interface{})["start_time"]
	job := alertData["details"].(map[string]interface{})["job"]
	pod := alertData["details"].(map[string]interface{})["pod"]
	service := alertData["details"].(map[string]interface{})["service"]

	log.Printf("Received Alert:\n"+
		"Alert Name: %s\n"+
		"Cluster: %s\n"+
		"Message: %s\n"+
		"Severity: %s\n"+
		"Namespace: %s\n"+
		"Target: %s\n"+
		"Start Time: %s\n"+
		"Job: %s\n"+
		"Pod: %s\n"+
		"Service: %s\n",
		alertname, cluster, message, severity, namespace, target,
		startTime, job, pod, service)
}

// ... (rest of your code)
