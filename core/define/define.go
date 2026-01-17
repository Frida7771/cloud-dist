package define

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.RegisteredClaims
}

var JwtKey = os.Getenv("JWT_KEY")
var SendGridAPIKey = os.Getenv("SendGridAPIKey")
var SendGridFromEmail = os.Getenv("SendGridFromEmail")

// InitJWTConfig initializes JWT key from config struct.
// Environment variables take precedence over config file values.
func InitJWTConfig(jwtKey string) {
	if JwtKey == "" && jwtKey != "" {
		JwtKey = jwtKey
	}
	// Fallback to default if still empty (for development only)
	if JwtKey == "" {
		JwtKey = "cloud-dist-key" // Default for development
	}
}

// CodeLength verification code length
var CodeLength = 6

// CodeExpire verification code expiration time (seconds)
var CodeExpire = 300

// AWS S3 configuration
var AWSAccessKeyID = os.Getenv("AWSAccessKeyID")
var AWSSecretAccessKey = os.Getenv("AWSSecretAccessKey")
var S3Bucket = os.Getenv("S3Bucket")
var S3Region = os.Getenv("AWSRegion")
var S3Endpoint = os.Getenv("S3Endpoint")                             // Optional custom endpoint
var S3UseAcceleration = os.Getenv("S3UseAcceleration") == "true"     // Enable S3 Transfer Acceleration

// InitS3Config initializes S3 configuration from config struct.
// Environment variables take precedence over config file values.
func InitS3Config(accessKeyID, secretAccessKey, bucket, region, endpoint string, useAcceleration bool) {
	// Only set if not already set by environment variable
	if AWSAccessKeyID == "" && accessKeyID != "" {
		AWSAccessKeyID = accessKeyID
	}
	if AWSSecretAccessKey == "" && secretAccessKey != "" {
		AWSSecretAccessKey = secretAccessKey
	}
	if S3Bucket == "" && bucket != "" {
		S3Bucket = bucket
	}
	if S3Region == "" && region != "" {
		S3Region = region
	}
	if S3Endpoint == "" && endpoint != "" {
		S3Endpoint = endpoint
	}
	// Enable acceleration if not set by env and config says true
	if os.Getenv("S3UseAcceleration") == "" && useAcceleration {
		S3UseAcceleration = true
	}
	log.Printf("[InitS3Config] Bucket=%s, Region=%s, UseAcceleration=%v", S3Bucket, S3Region, S3UseAcceleration)
}

// InitSendGridConfig initializes SendGrid configuration from config struct.
// Environment variables take precedence over config file values.
func InitSendGridConfig(apiKey, fromEmail string) {
	// Only set if not already set by environment variable
	if SendGridAPIKey == "" && apiKey != "" {
		SendGridAPIKey = apiKey
	}
	if SendGridFromEmail == "" && fromEmail != "" {
		SendGridFromEmail = fromEmail
	}
}

// PageSize default pagination parameter
var PageSize = 20

var Datetime = "2006-01-02 15:04:05"

var TokenExpire = 3600
var RefreshTokenExpire = 7200

// Stripe configuration
var StripeSecretKey = os.Getenv("STRIPE_SECRET_KEY")
var StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

// InitStripeConfig initializes Stripe configuration from config struct.
// Environment variables take precedence over config file values.
func InitStripeConfig(secretKey, webhookSecret string) {
	if StripeSecretKey == "" && secretKey != "" {
		StripeSecretKey = secretKey
	}
	if StripeWebhookSecret == "" && webhookSecret != "" {
		StripeWebhookSecret = webhookSecret
	}
	// Log initialization (partial secret for security)
	if StripeWebhookSecret != "" {
		secretLen := len(StripeWebhookSecret)
		if secretLen > 20 {
			log.Printf("[InitStripeConfig] WebhookSecret initialized: %s...%s (len: %d)",
				StripeWebhookSecret[:10],
				StripeWebhookSecret[secretLen-10:],
				secretLen)
		}
	}
}

// Storage pricing tiers (in bytes)
const (
	Storage10GB  = 10 * 1024 * 1024 * 1024   // 10GB
	Storage50GB  = 50 * 1024 * 1024 * 1024   // 50GB
	Storage100GB = 100 * 1024 * 1024 * 1024  // 100GB
	Storage500GB = 500 * 1024 * 1024 * 1024  // 500GB
	Storage1TB   = 1024 * 1024 * 1024 * 1024 // 1TB
)

// GetStoragePrice returns the price in cents for a given storage amount
func GetStoragePrice(storageBytes int64) int64 {
	// Pricing: $0.10 per GB per month (example pricing, adjust as needed)
	// For one-time purchase, we'll use a different model
	// Example: 10GB = $9.99, 50GB = $39.99, 100GB = $69.99, 500GB = $299.99, 1TB = $499.99
	switch storageBytes {
	case Storage10GB:
		return 999 // $9.99
	case Storage50GB:
		return 3999 // $39.99
	case Storage100GB:
		return 6999 // $69.99
	case Storage500GB:
		return 29999 // $299.99
	case Storage1TB:
		return 49999 // $499.99
	default:
		// Custom pricing: $0.10 per GB (minimum $4.99)
		price := (storageBytes / (1024 * 1024 * 1024)) * 10 // $0.10 per GB
		if price < 499 {
			price = 499 // Minimum $4.99
		}
		return price
	}
}
