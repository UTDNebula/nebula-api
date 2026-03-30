package routes

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gin-gonic/gin"
	"github.com/wneessen/go-mail"

	"github.com/UTDNebula/nebula-api/api/controllers"
	"github.com/UTDNebula/nebula-api/api/schema"
)

var emailClient *mail.Client
var smtpUsername string
var emailSendKey string
var emailClientOnce sync.Once

var tasksClient *cloudtasks.Client
var queuePath string
var queueUrl string
var tasksClientOnce sync.Once

func initTasksClient() (*cloudtasks.Client, string, string) {
	// Singleton to prevent multiple clients
	tasksClientOnce.Do(func() {
		qPath := os.Getenv("GCLOUD_EMAIL_QUEUE_PATH")
		qUrl := os.Getenv("GCLOUD_EMAIL_QUEUE_URL")

		if qPath == "" || qUrl == "" {
			log.Println("Cloud Tasks environment variables are not fully configured; skipping email queuing routes")
			return
		}

		ctx := context.Background()
		c, err := cloudtasks.NewClient(ctx)
		if err != nil {
			log.Printf("Failed to create Cloud Tasks client: %v", err)
			return
		}
		tasksClient = c
		queuePath = qPath
		queueUrl = qUrl
	})
	return tasksClient, queuePath, queueUrl
}

func initEmailClient() (*mail.Client, string, string) {
	// Singleton to prevent multiple clients
	emailClientOnce.Do(func() {
		smtpHost := os.Getenv("SMTP_HOST")
		smtpUser := os.Getenv("SMTP_USERNAME")
		smtpPass := os.Getenv("SMTP_PASSWORD")
		sendKey := os.Getenv("EMAIL_SEND_ROUTE_KEY")

		if smtpHost == "" || smtpUser == "" || smtpPass == "" || sendKey == "" {
			log.Println("SMTP environment variables are not fully configured; skipping email routes")
			return
		}

		c, err := mail.NewClient(
			smtpHost,
			mail.WithTLSPortPolicy(mail.TLSMandatory),
			mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
			mail.WithUsername(smtpUser),
			mail.WithPassword(smtpPass),
		)
		if err != nil {
			log.Printf("Failed to create SMTP client: %v", err)
			return
		}
		emailClient = c
		smtpUsername = smtpUser
	})
	return emailClient, smtpUsername, emailSendKey
}

func EmailRoute(router *gin.Engine) {
	client, username, emailSendKey := initEmailClient()
	tClient, qPath, qUrl := initTasksClient()

	if client == nil {
		log.Println("SMTP client not initialized")
	}

	if tClient == nil {
		log.Println("Cloud Tasks client not initialized")
	}

	if client == nil || tClient == nil {
		log.Println("skipping email routes")
		return
	}

	// Restrict with password
	authMiddleware := func(key string, envKey string) gin.HandlerFunc {
		return func(c *gin.Context) {
			secret := c.GetHeader(key)
			expected, exist := os.LookupEnv(envKey)
			if !exist || secret != expected {
				c.AbortWithStatusJSON(http.StatusForbidden, schema.APIResponse[string]{Status: http.StatusForbidden, Message: "error", Data: "Forbidden"})
				return
			}
			c.Next()
		}
	}

	// All routes related to email come here
	emailGroup := router.Group("/email")

	// Pass to next layer
	emailGroup.Use(func(c *gin.Context) {
		c.Set("emailClient", client)
		c.Set("emailUsername", username)
		c.Set("emailSendKey", emailSendKey)
		c.Set("tasksClient", tClient)
		c.Set("queuePath", qPath)
		c.Set("queueUrl", qUrl)
		c.Next()
	})

	emailGroup.OPTIONS("", controllers.Preflight)
	emailGroup.POST("/send", authMiddleware("x-email-send-key", "EMAIL_SEND_ROUTE_KEY"), controllers.SendEmail)
	emailGroup.POST("/queue", authMiddleware("x-email-queue-key", "EMAIL_QUEUE_ROUTE_KEY"), controllers.QueueEmail)
}
