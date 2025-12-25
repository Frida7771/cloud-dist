# CloudDisk

> Lightweight cloud drive backend built with Gin + GORM.

## Getting Started

```bash
# Install dependencies
go mod tidy
# Start the service
go run ./cmd/cloud-disk -config configs/config.yaml
```

### Environment Variables

```bash
# AWS S3
export AWSAccessKeyID=AKIAxxx
export AWSSecretAccessKey=xxxx
export S3Bucket=your-bucket-name
export AWSRegion=us-east-1
# Optional: custom endpoint (S3-compatible)
# export S3Endpoint=https://s3.amazonaws.com

# SendGrid
export SendGridAPIKey=SG.xxxxxx
```

- AWS S3 Console: https://s3.console.aws.amazon.com/s3/home  
- AWS S3 Docs: https://docs.aws.amazon.com/s3/index.html

## Features

- **User Module**
  - Password login
  - Authorization refresh
  - Email registration
  - User detail & quota
- **Storage Pool**
  - Central repository: upload, instant upload, multipart upload, AWS S3 integration
  - Personal repository: link files, list files, rename, create folders, delete, move
- **Share Module**
  - Create share record
  - Retrieve resource detail
  - Save shared resource to personal space
