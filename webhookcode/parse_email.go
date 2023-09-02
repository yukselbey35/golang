// we define the alert template and use the Go template package to parse and execute it with the email data. The extracted and formatted data is then logged in your desired format

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var emailContent string
	// Extract the email content from the request body (you may need to implement this part)

	// Define the template based on your alert-template.yaml structure
	alertTemplate := `
    Knowledge Base:
    {{ template "__knowledgeBaseLink" .}}
    ______________________________________________________________________________________
    {{ range .Alerts.Firing }}
    {{ if .Labels.alertname}}Alertname:   {{ .Labels.alertname }}{{ end}}
    {{ if .Labels.cluster}}Cluster:   {{ .Labels.cluster }}{{ end}}
    {{ if .Annotations.message}}Details:    {{ .Annotations.message }}{{ end}}
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
	var result string
	err = tmpl.Execute(&result, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}

	// Process the result, which contains the formatted alert data
	log.Printf("Received Alert:\n%s", result)

	w.WriteHeader(http.StatusOK)
}

// Rest of your code...
