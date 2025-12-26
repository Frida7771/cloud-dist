package handler

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/svc"

	"github.com/gin-gonic/gin"
)

func StoragePurchaseWebhookHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[WebhookHandler] ===== WEBHOOK REQUEST RECEIVED =====")
		log.Printf("[WebhookHandler] Method: %s, Path: %s, IP: %s",
			c.Request.Method, c.Request.URL.Path, c.ClientIP())
		log.Printf("[WebhookHandler] Headers: %v", c.Request.Header)

		// Read the request body (must read before any other processing)
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Printf("[WebhookHandler] ERROR: Failed to read request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		// Create a new ReadCloser from the payload for potential reuse
		c.Request.Body = io.NopCloser(bytes.NewReader(payload))

		log.Printf("[WebhookHandler] Payload size: %d bytes", len(payload))
		if len(payload) == 0 {
			log.Printf("[WebhookHandler] ERROR: Empty payload")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty request body"})
			return
		}

		// Get the Stripe signature from headers
		signature := c.GetHeader("Stripe-Signature")
		if signature == "" {
			log.Printf("[WebhookHandler] ERROR: Missing Stripe-Signature header")
			log.Printf("[WebhookHandler] All headers: %v", c.Request.Header)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Stripe-Signature header"})
			return
		}

		log.Printf("[WebhookHandler] Stripe-Signature header present (length: %d), processing webhook...", len(signature))

		// Process the webhook
		l := logic.NewStoragePurchaseWebhookLogic(c.Request.Context(), svcCtx)
		if err := l.ProcessWebhook(payload, signature); err != nil {
			log.Printf("[WebhookHandler] ERROR: Webhook processing failed: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("[WebhookHandler] Webhook processed successfully")
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
