# CloudDist

A lightweight cloud storage backend built with Go, Gin, and GORM.

## Prerequisites

- Go 1.23+
- MySQL
- Redis
- AWS S3 (or S3-compatible storage like MinIO)

## Quick Start

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Setup Database

```bash
mysql -u root -p < setup_db.sql
```

### 3. Configure

Edit `configs/config.yaml` with your MySQL and Redis connection details.

### 4. Set Environment Variables (Optional)

```bash
export AWSAccessKeyID=your-access-key
export AWSSecretAccessKey=your-secret-key
export S3Bucket=your-bucket-name
export AWSRegion=us-east-1
export SendGridAPIKey=your-sendgrid-key  # For email verification
```

### 5. Run

```bash
go run ./cmd/cloud-dist/main.go -config configs/config.yaml
```

Service runs on `http://0.0.0.0:8888`

## Features

- **User Management**: Registration, login, logout with JWT authentication
- **File Storage**: Upload, download, organize files with S3 backend
- **File Management**: Create folders, rename, move, delete files
- **Sharing**: Public file sharing with expiration
- **Friend System**: Add friends, send friend requests, share files with friends
- **Token Management**: JWT-based auth with refresh tokens and blacklist support


## Architecture

- **Framework**: Gin web framework
- **ORM**: GORM
- **Storage**: AWS S3
- **Cache**: Redis (for verification codes and token blacklist)
- **Auth**: JWT tokens with Redis blacklist

## License

MIT
