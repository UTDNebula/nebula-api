package controllers

import (
	"net/http"
	"os"

	_ "github.com/UTDNebula/nebula-api/api/schema"
	"github.com/gin-gonic/gin"
	"github.com/wneessen/go-mail"
)

type EmailRequest struct {
	From    string `json:"from" binding:"required"`
	To      string `json:"to" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// @Id				sendEmail
// @Router			/email [post]
// @Description		"Send an email via SMTP"
// @Accept			json
// @Produce			json
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
