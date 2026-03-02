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
var smtpFromAddr string
var emailClientOnce sync.Once

var tasksClient *cloudtasks.Client
var queuePath string
var queueUrl string
var tasksClientOnce sync.Once

func initTasksClient() (*cloudtasks.Client, string, string) {
	tasksClientOnce.Do(func() {
		qPath := os.Getenv("QUEUE_PATH")
		qUrl := os.Getenv("EMAIL_URL") // TODO: Consider a different name

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

func initEmailClient() (*mail.Client, string) {
	emailClientOnce.Do(func() {
		smtpHost := os.Getenv("SMTP_HOST") // TODO: use lookupenv instead
		smtpUser := os.Getenv("SMTP_USERNAME")
		smtpPass := os.Getenv("SMTP_PASSWORD")
		smtpFrom := os.Getenv("SMTP_FROM")

		if smtpHost == "" || smtpUser == "" || smtpPass == "" || smtpFrom == "" {
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
		smtpFromAddr = smtpFrom
	})
	return emailClient, smtpFromAddr
}

func EmailRoute(router *gin.Engine) {
	client, fromAddr := initEmailClient()
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

	// Rescrict with password
	authMiddleware := func(c *gin.Context) {
		secret := c.GetHeader("x-email-key")
		expected, exist := os.LookupEnv("EMAIL_ROUTE_KEY")
		if !exist || secret != expected {
			c.AbortWithStatusJSON(http.StatusForbidden, schema.APIResponse[string]{Status: http.StatusForbidden, Message: "error", Data: "Forbidden"})
			return
		}
		c.Next()
	}

	// All routes related to email come here
	emailGroup := router.Group("/email")

	// Pass to next layer
	emailGroup.Use(func(c *gin.Context) {
		c.Set("emailClient", client)
		c.Set("emailFrom", fromAddr)
		c.Set("tasksClient", tClient)
		c.Set("queuePath", qPath)
		c.Set("queueUrl", qUrl)
		c.Next()
	})

	// Use auth
	emailGroup.Use(authMiddleware)

	emailGroup.OPTIONS("", controllers.Preflight)
	emailGroup.POST("/send", controllers.SendEmail)
	emailGroup.POST("/queue", controllers.QueueEmail)
}
