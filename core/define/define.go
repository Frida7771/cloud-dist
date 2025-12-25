package define

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaim struct {
	Id       int
	Identity string
	Name     string
	jwt.RegisteredClaims
}

var JwtKey = "cloud-disk-key"
var SendGridAPIKey = os.Getenv("SendGridAPIKey")
var SendGridFromEmail = os.Getenv("SendGridFromEmail")

// CodeLength verification code length
var CodeLength = 6

// CodeExpire verification code expiration time (seconds)
var CodeExpire = 300

// AWS S3 configuration
var AWSAccessKeyID = os.Getenv("AWSAccessKeyID")
var AWSSecretAccessKey = os.Getenv("AWSSecretAccessKey")
var S3Bucket = os.Getenv("S3Bucket")
var S3Region = os.Getenv("AWSRegion")
var S3Endpoint = os.Getenv("S3Endpoint") // Optional custom endpoint

// InitS3Config initializes S3 configuration from config struct.
// Environment variables take precedence over config file values.
func InitS3Config(accessKeyID, secretAccessKey, bucket, region, endpoint string) {
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
}

// PageSize default pagination parameter
var PageSize = 20

var Datetime = "2006-01-02 15:04:05"

var TokenExpire = 3600
var RefreshTokenExpire = 7200
