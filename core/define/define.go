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

// CodeLength 验证码长度
var CodeLength = 6

// CodeExpire 验证码过期时间（s）
var CodeExpire = 300

// AWS S3 配置
var AWSAccessKeyID = os.Getenv("AWSAccessKeyID")
var AWSSecretAccessKey = os.Getenv("AWSSecretAccessKey")
var S3Bucket = os.Getenv("S3Bucket")
var S3Region = os.Getenv("AWSRegion")
var S3Endpoint = os.Getenv("S3Endpoint") // 可选自定义 Endpoint

// PageSize 分页的默认参数
var PageSize = 20

var Datetime = "2006-01-02 15:04:05"

var TokenExpire = 3600
var RefreshTokenExpire = 7200
