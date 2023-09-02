// Added the missing import statements for io.
//
// Implemented the extractEmailContent function to handle the extraction of email content from the request body. You need to implement the actual logic to extract the email content based on your requirements.
//
// Removed the var emailContent string declaration, as email content is now extracted dynamically.
//
// Please note that the extractEmailContent function's implementation depends on how the email content is structured in your incoming request. You need to adapt it to your specific use case to correctly extract the email content.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"text/template"
)

// Mapping of predefined namespaces to their owners' email addresses
var ownerEmails = map[string]string{
	"k8-testing": "owner1@example.com",
	// Add more mappings as needed
}

// SMTP server settings
var smtpServer = "smtp.office365.com"
var smtpPort = 587
var senderEmail = "your_email@example.com"
var senderPassword = "your_email_password"

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	// Extract the email content from the request body
	emailContent, err := extractEmailContent(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to extract email content: %v", err), http.StatusInternalServerError)
		return
	}

	// Define the alert template based on your alert-template.yaml structure
	alertTemplate := `
    Knowledge Base:
    {{ template "__knowledgeBaseLink" . }}
    ______________________________________________________________________________________
    {{ range .Alerts.Firing }}
    {{ if .Labels.alertname}}Alert Name:   {{ .Labels.alertname }}{{ end}}
    {{ if .Labels.cluster}}Cluster:   {{ .Labels.cluster }}{{ end}}
    {{ if .Annotations.message}}Message:    {{ .Annotations.message }}{{ end}}
    {{ if .Labels.namespace}}Namespace:   {{ .Labels.namespace }}{{ end}}
    {{ if .Labels.severity}}Severity:   {{ .Labels.severity }}{{ end}}
    {{ if .Labels.statefulset}}Statefulset:   {{ .Labels.statefulset }}{{ end}}
    {{ if .Labels.instance}}Target:   {{ .Labels.instance }}{{ end}}
    Start Time:   {{ .StartsAt }}
    {{ if .Labels.prometheus}}Prometheus Instance:   {{ .Labels.prometheus }}{{ end}}
    {{ if .Labels.service}}Service:   {{ .Labels.service }}{{ end}}
    {{ if .Labels.job}}Job:   {{ .Labels.job }}{{ end}}
    {{ if .Labels.pod}}Pod:   {{ .Labels.pod }}{{ end}}
    Alert Source:   {{ .GeneratorURL }}
    _____________________________________________________________________________________
    {{ end }}
    {{end}}
    `

	// Parse the template
	tmpl, err := template.New("alertTemplate").Parse(alertTemplate)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse the email content using the template
	var data map[string]interface{}
	err = json.Unmarshal([]byte(emailContent), &data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to unmarshal email content: %v", err), http.StatusInternalServerError)
		return
	}

	// Execute the template with the email data
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	// Process the result, which contains the formatted alert data
	logHumanReadableAlert(result.String())

	// Forward critical alerts to namespace owner's email
	alerts := data["Alerts"].([]interface{})
	for _, alert := range alerts {
		alertMap := alert.(map[string]interface{})
		severity := alertMap["Labels"].(map[string]interface{})["severity"].(string)
		namespace := alertMap["Labels"].(map[string]interface{})["namespace"].(string)

		if severity == "critical" && isPredefinedNamespace(namespace) {
			namespaceOwnerEmail := ownerEmails[namespace]
			if namespaceOwnerEmail != "" {
				// Send email to namespace owner
				if err := sendEmail(result.String(), namespaceOwnerEmail); err != nil {
					log.Printf("Error sending email: %v", err)
				}
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func extractEmailContent(body io.Reader) (string, error) {
	// Implement the extraction logic from the request body
	// Return the extracted email content as a string
}

func sendEmail(content, recipientEmail string) error {
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	toEmail := recipientEmail
	subject := "Alert from Alertmanager"
	emailBody := content

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + emailBody)

	err := smtp.SendMail(smtpServer+":"+fmt.Sprint(smtpPort), auth, senderEmail, []string{toEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}

func isPredefinedNamespace(namespace string) bool {
	_, exists := ownerEmails[namespace]
	return exists
}

func logHumanReadableAlert(alertData string) {
	// Split the email content into lines
	lines := strings.Split(alertData, "\n")

	// Initialize variables for capturing alert details
	var alertname, cluster, message, severity, namespace, target, startTime, job, pod, service string

	// Iterate through lines and capture details
	for _, line := range lines {
		if strings.Contains(line, "Alert Name:") {
			alertname = strings.TrimSpace(strings.TrimPrefix(line, "Alert Name:"))
		} else if strings.Contains(line, "Cluster:") {
			cluster = strings.TrimSpace(strings.TrimPrefix(line, "Cluster:"))
		} else if strings.Contains(line, "Message:") {
			message = strings.TrimSpace(strings.TrimPrefix(line, "Message:"))
		} else if strings.Contains(line, "Severity:") {
			severity = strings.TrimSpace(strings.TrimPrefix(line, "Severity:"))
		} else if strings.Contains(line, "Namespace:") {
			namespace = strings.TrimSpace(strings.TrimPrefix(line, "Namespace:"))
		} else if strings.Contains(line, "Target:") {
			target = strings.TrimSpace(strings.TrimPrefix(line, "Target:"))
		} else if strings.Contains(line, "Start Time:") {
			startTime = strings.TrimSpace(strings.TrimPrefix(line, "Start Time:"))
		} else if strings.Contains(line, "Job:") {
			job = strings.TrimSpace(strings.TrimPrefix(line, "Job:"))
		} else if strings.Contains(line, "Pod:") {
			pod = strings.TrimSpace(strings.TrimPrefix(line, "Pod:"))
		} else if strings.Contains(line, "Service:") {
			service = strings.TrimSpace(strings.TrimPrefix(line, "Service:"))
		}
	}

	// Log the captured alert details
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

func main() {
	http.HandleFunc("/webhook", webhookHandler)
	port := "5000"
	log.Printf("Webhook server is listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
