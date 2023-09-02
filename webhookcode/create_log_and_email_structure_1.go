	// Log the received alert data in a human-readable format
	logHumanReadableAlert(alertData)

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





The createEmailBody function creates the body of the email with the extracted information.

// This sendEmail function takes the alert interface as input and extracts the required fields (alertname, cluster, namespace, message, and severity) from the alert interface using type assertions. It then composes the email subject and body using this extracted information.
func sendEmail(alert interface{}, recipientEmail string) error {
	alertMap := alert.(map[string]interface{})
	annotations := alertMap["annotations"].(map[string]interface{})
	alertname := annotations["alertname"].(string)
	cluster := alertMap["labels"].(map[string]interface{})["cluster"].(string)
	namespace := alertMap["labels"].(map[string]interface{})["namespace"].(string)
	message := annotations["message"].(string)
	severity := alertMap["labels"].(map[string]interface{})["severity"].(string)

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpServer)

	toEmail := recipientEmail
	subject := fmt.Sprintf("Alert: %s - %s in Namespace: %s", severity, alertname, namespace)
	emailBody := createEmailBody(alertname, cluster, namespace, message)

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + emailBody)

	err := smtp.SendMail(smtpServer+":"+fmt.Sprint(smtpPort), auth, senderEmail, []string{toEmail}, msg)
	if err != nil {
		return err
	}

	return nil
}

func createEmailBody(alertname, cluster, namespace, message string) string {
	return fmt.Sprintf(`
Alert Name: %s
Cluster: %s
Namespace: %s
Message: %s
`,
		alertname, cluster, namespace, message)
}
