package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	_ "github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"github.com/wneessen/go-mail"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

type EmailRequest struct {
	From    string `json:"from" binding:"required"`
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// @Id				sendEmail
// @Router			/email [post]
// @Description	"Send an email via SMTP"
// @Accept			json
// @Produce		json
// @Param			request	body		EmailRequest				true	"Email Request Body"
// @Success		200		{object}	schema.APIResponse[string]	"Email sent successfully"
// @Failure		500		{object}	schema.APIResponse[string]	"A string describing the error"
// @Failure		400		{object}	schema.APIResponse[string]	"A string describing the error"
func SendEmail(c *gin.Context) {

	var req EmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	smtpHost := os.Getenv("SMTP_HOST") // TODO: We should be using env.go for this
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpFrom := os.Getenv("SMTP_FROM")

	if smtpHost == "" || smtpFrom == "" {
		respond(c, http.StatusInternalServerError, "SMTP configuration is missing", "SMTP environment variables are not fully configured")
		return
	}

	m := mail.NewMsg()
	if err := m.FromFormat(req.From, smtpFrom); err != nil {
		respond(c, http.StatusInternalServerError, "failed to set from address", err.Error())
		return
	}

	if err := m.To(req.To); err != nil {
		respond(c, http.StatusBadRequest, "invalid to address", err.Error())
		return
	}

	m.Subject(req.Subject)
	m.SetBodyString(mail.TypeTextPlain, req.Body)

	client, err := mail.NewClient(
		smtpHost,
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
		mail.WithUsername(smtpUser),
		mail.WithPassword(smtpPass),
	)

	if err != nil {
		respond(c, http.StatusInternalServerError, "failed to setup SMTP client", err.Error())
		return
	}

	if err := client.DialAndSend(m); err != nil {
		respond(c, http.StatusInternalServerError, "failed to send email", err.Error())
		return
	}

	respond(c, http.StatusOK, "success", "Email sent successfully")
}

func QueueEmail(c *gin.Context) {

	// Request must be able to bind to email request
	if err := c.ShouldBindJSON(&EmailRequest{}); err != nil {
		respond(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	var body []byte
	c.Request.Body.Read(body)

	queuePath := os.Getenv("QUEUE_PATH")
	url := os.Getenv("EMAIL_URL") // TODO: Consider a different name

	_, err := createHTTPTask(queuePath, url, body)

	if err != nil {
		respond(c, http.StatusInternalServerError, "failed to queue email", err.Error())
		return
	}

	respond(c, http.StatusOK, "success", "Email queued successfully") // TODO: Change the response

}

// createHTTPTask creates a new task with a HTTP target then adds it to a Queue.
func createHTTPTask(queuePath string, url string, body []byte) (*taskspb.Task, error) {

	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
				},
			},
		},
	}

	// Add a payload message if one is present.
	req.Task.GetHttpRequest().Body = []byte(body)

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %w", err)
	}

	return createdTask, nil
}

// createHTTPTaskWithToken constructs a task with a authorization token
// and HTTP target then adds it to a Queue.
func createHTTPTaskWithToken(projectID, locationID, queueID, url, email, message string) (*taskspb.Task, error) {
	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2
	ctx := context.Background()
	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &taskspb.CreateTaskRequest{
		Parent: queuePath,
		Task: &taskspb.Task{
			// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#HttpRequest
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        url,
					AuthorizationHeader: &taskspb.HttpRequest_OidcToken{
						OidcToken: &taskspb.OidcToken{
							ServiceAccountEmail: email,
						},
					},
				},
			},
		},
	}

	// Add a payload message if one is present.
	req.Task.GetHttpRequest().Body = []byte(message)

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %w", err)
	}

	return createdTask, nil
}
