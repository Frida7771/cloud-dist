package helper

import (
	"bytes"
	"cloud-disk/core/define"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func GenerateToken(id int, identity, name string, second int) (string, error) {
	uc := define.UserClaim{
		Id:       id,
		Identity: identity,
		Name:     name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(second))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenString, err := token.SignedString([]byte(define.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyzeToken
// Token 解析
func AnalyzeToken(token string) (*define.UserClaim, error) {
	uc := new(define.UserClaim)
	claims, err := jwt.ParseWithClaims(token, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(define.JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return uc, errors.New("token is invalid")
	}
	return uc, err
}

// MailSendCode
// 邮箱验证码发送
func MailSendCode(emailAddr, code string) error {
	apiKey := define.SendGridAPIKey
	if apiKey == "" {
		return errors.New("SendGrid API key is not configured")
	}

	from := sgmail.NewEmail("CloudDisk", "frida16571@gmail.com")
	to := sgmail.NewEmail("", emailAddr)
	subject := "CloudDisk 验证码"
	plain := "你的验证码为：" + code
	html := fmt.Sprintf("你的验证码为：<h1>%s</h1>", code)
	message := sgmail.NewSingleEmail(from, subject, to, plain, html)

	client := sendgrid.NewSendClient(apiKey)
	resp, err := client.Send(message)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid error: status=%d body=%s", resp.StatusCode, resp.Body)
	}
	return nil
}

func RandCode() string {
	s := "1234567890"
	code := ""
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < define.CodeLength; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

func UUID() string {
	return uuid.NewV4().String()
}

var (
	s3Client     *s3.Client
	s3ClientOnce sync.Once
	s3ClientErr  error
)

func getS3Client(ctx context.Context) (*s3.Client, error) {
	s3ClientOnce.Do(func() {
		if define.S3Bucket == "" {
			s3ClientErr = errors.New("S3Bucket is not configured")
			return
		}

		loadOpts := []func(*config.LoadOptions) error{
			config.WithRegion(define.S3Region),
		}
		if define.AWSAccessKeyID != "" && define.AWSSecretAccessKey != "" {
			creds := credentials.NewStaticCredentialsProvider(define.AWSAccessKeyID, define.AWSSecretAccessKey, "")
			loadOpts = append(loadOpts, config.WithCredentialsProvider(creds))
		}

		cfg, err := config.LoadDefaultConfig(ctx, loadOpts...)
		if err != nil {
			s3ClientErr = err
			return
		}

		s3Client = s3.NewFromConfig(cfg, func(o *s3.Options) {
			if define.S3Endpoint != "" {
				o.BaseEndpoint = aws.String(define.S3Endpoint)
				o.UsePathStyle = true
			}
		})
	})

	return s3Client, s3ClientErr
}

// S3ObjectURL 返回对象的完整访问地址
func S3ObjectURL(key string) string {
	base := strings.TrimSuffix(define.S3Endpoint, "/")
	if base != "" {
		if !strings.HasPrefix(base, "http") {
			base = "https://" + base
		}
		return fmt.Sprintf("%s/%s/%s", base, define.S3Bucket, key)
	}

	region := define.S3Region
	if region == "" {
		region = "us-east-1"
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", define.S3Bucket, region, key)
}

// S3Upload 上传文件到 AWS S3
func S3Upload(r *http.Request) (string, error) {
	ctx := context.Background()
	client, err := getS3Client(ctx)
	if err != nil {
		return "", err
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	key := "cloud-disk/" + UUID() + path.Ext(fileHeader.Filename)
	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(define.S3Bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		return "", err
	}
	return S3ObjectURL(key), nil
}

// MultipartPart 上传的分片信息
type MultipartPart struct {
	PartNumber int32
	ETag       string
}

// S3InitPart 分片上传初始化
func S3InitPart(ext string) (string, string, error) {
	ctx := context.Background()
	client, err := getS3Client(ctx)
	if err != nil {
		return "", "", err
	}

	key := "cloud-disk/" + UUID() + ext
	resp, err := client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket: aws.String(define.S3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", "", err
	}
	return key, aws.ToString(resp.UploadId), nil
}

// S3PartUpload 分片上传
func S3PartUpload(r *http.Request) (string, error) {
	ctx := context.Background()
	client, err := getS3Client(ctx)
	if err != nil {
		return "", err
	}

	key := r.PostForm.Get("key")
	uploadID := r.PostForm.Get("upload_id")
	partNumber, err := strconv.Atoi(r.PostForm.Get("part_number"))
	if err != nil {
		return "", err
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, file); err != nil {
		return "", err
	}

	resp, err := client.UploadPart(ctx, &s3.UploadPartInput{
		Bucket:        aws.String(define.S3Bucket),
		Key:           aws.String(key),
		UploadId:      aws.String(uploadID),
		PartNumber:    aws.Int32(int32(partNumber)),
		Body:          bytes.NewReader(buf.Bytes()),
		ContentLength: aws.Int64(int64(buf.Len())),
	})
	if err != nil {
		return "", err
	}
	return strings.Trim(aws.ToString(resp.ETag), "\""), nil
}

// S3PartUploadComplete 分片上传完成
func S3PartUploadComplete(key, uploadID string, parts []MultipartPart) error {
	ctx := context.Background()
	client, err := getS3Client(ctx)
	if err != nil {
		return err
	}

	completed := make([]s3types.CompletedPart, 0, len(parts))
	for _, part := range parts {
		pn := part.PartNumber
		completed = append(completed, s3types.CompletedPart{
			ETag:       aws.String(part.ETag),
			PartNumber: aws.Int32(pn),
		})
	}

	_, err = client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(define.S3Bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &s3types.CompletedMultipartUpload{
			Parts: completed,
		},
	})
	return err
}
