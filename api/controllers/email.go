package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"github.com/wneessen/go-mail"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

// @Id				sendEmail
// @Router			/email/send [post]
// @Tags			Internal
// @Description	"Send an email via SMTP. This route is restricted to only Nebula Labs internal Projects."
// @Accept			json
// @Produce		json
// @Param			request				body		schema.EmailRequest						true	"Email Request Body"
// @Param			x-email-send-key	header		string									true	"The internal email send key"
// @Success		200					{object}	schema.APIResponse[schema.EmailRequest]	"Email Request Body"
// @Failure		500					{object}	schema.APIResponse[string]				"A string describing the error"
// @Failure		400					{object}	schema.APIResponse[string]				"A string describing the error"
func SendEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var req schema.EmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respond(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	client := c.MustGet("emailClient").(*mail.Client)
	smtpUsername := os.Getenv("SMTP_USERNAME")

	m := mail.NewMsg()
	if err := m.FromFormat(req.From, smtpUsername); err != nil {
		respond(c, http.StatusInternalServerError, "failed to set from address", err.Error())
		return
	}

	if err := m.To(req.To...); err != nil {
		respond(c, http.StatusBadRequest, "invalid to address(es)", err.Error())
		return
	}

	m.Subject(req.Subject)
	m.SetBodyString(mail.TypeTextHTML, req.Body)

	for _, att := range req.Attachments {
		m.AttachReader(att.Name, bytes.NewReader(att.Data))
	}

	for _, emb := range req.Embeds {
		m.EmbedReader(emb.Name, bytes.NewReader(emb.Data), mail.WithFileContentID(emb.Name))
	}

	if err := client.DialAndSendWithContext(ctx, m); err != nil {
		respond(c, http.StatusInternalServerError, "failed to send email", err.Error())
		return
	}

	respond(c, http.StatusOK, "success", req)
}

// @Id				QueueEmail
// @Router			/email/queue [post]
// @Tags			Internal
// @Description	"Queue an email to be sent via SMTP. Multi-recipient emails will be queued as separate emails to avoid bypassing queueing system. This route is restricted to only Nebula Labs internal Projects."
// @Accept			json
// @Produce		json
// @Param			request				body		schema.EmailRequest				true	"Email Request Body"
// @Param			x-email-queue-key	header		string							true	"The internal email queue key"
// @Success		200					{object}	schema.APIResponse[[]string]	"The list of queued task names"
// @Failure		500					{object}	schema.APIResponse[string]		"A string describing the error"
// @Failure		400					{object}	schema.APIResponse[string]		"A string describing the error"
func QueueEmail(c *gin.Context) {
	// Request must be able to bind to email request
	var emailReq schema.EmailRequest
	if err := c.ShouldBindJSON(&emailReq); err != nil {
		respond(c, http.StatusBadRequest, "invalid request payload", err.Error())
		return
	}

	client := c.MustGet("tasksClient").(*cloudtasks.Client)
	queuePath := os.Getenv("GCLOUD_EMAIL_QUEUE_PATH")
	queueUrl := os.Getenv("GCLOUD_EMAIL_QUEUE_URL")

	baseEmailReq := emailReq
	baseEmailReq.To = []string{}

	queuedTasks := []string{}

	numOfRecipients := len(emailReq.To)
	for i, to := range emailReq.To {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		baseEmailReq.To = []string{to}

		body, err := json.Marshal(baseEmailReq)
		if err != nil {
			respond(c, http.StatusInternalServerError, fmt.Sprintf("failed to serialize email request for recipient %s %d/%d", to, i, numOfRecipients), err.Error())
			return
		}

		// Build the Task payload.
		// // https://docs.cloud.google.com/tasks/docs/creating-http-target-tasks
		taskReq := &taskspb.CreateTaskRequest{
			Parent: queuePath,
			Task: &taskspb.Task{
				MessageType: &taskspb.Task_HttpRequest{
					HttpRequest: &taskspb.HttpRequest{
						HttpMethod: taskspb.HttpMethod_POST,
						Url:        queueUrl,
						Headers: map[string]string{
							"x-email-send-key": os.Getenv("EMAIL_SEND_ROUTE_KEY"), // Must get from env bc queue only has x-email-queue-key header
							"x-api-key":        c.GetHeader("x-api-key"),
						},
					},
				},
			},
		}

		// Add a payload message if one is present.
		taskReq.Task.GetHttpRequest().Body = []byte(body)

		task, err := client.CreateTask(ctx, taskReq)
		if err != nil {
			respond(c, http.StatusInternalServerError, fmt.Sprintf("failed to queue email for recipient %s %d/%d", to, i, numOfRecipients), err.Error())
			return
		}

		queuedTasks = append(queuedTasks, task.GetName())
	}

	respond(c, http.StatusOK, "success", queuedTasks)
}
