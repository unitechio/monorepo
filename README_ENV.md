# Environment Configuration Guide

Dự án này hỗ trợ cấu hình qua các file `.env` cho các môi trường khác nhau.

## Cách hoạt động

1. Hệ thống sẽ tự động load file `.env.{ENV}` dựa trên biến môi trường `ENV`
2. Nếu không tìm thấy file `.env.{ENV}`, sẽ fallback về file `.env`
3. Các biến môi trường hệ thống sẽ override các giá trị trong file `.env`

## Thiết lập

### 1. Tạo file `.env` hoặc `.env.development` trong thư mục root hoặc `api/`

Copy từ template dưới đây và điều chỉnh các giá trị:

```env
# Environment
ENV=development

# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=oss_monorepo
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=10

# MinIO Configuration
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET=oss-monorepo

# JWT Configuration
JWT_SECRET=change-this-secret-key-in-production
JWT_ACCESS_TOKEN_EXPIRE=15m
JWT_REFRESH_TOKEN_EXPIRE=168h

# OTP Configuration
OTP_EXPIRE=30s
OTP_LENGTH=6

# 2FA Configuration
TWO_FA_ISSUER=OSS-Monorepo

# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM_EMAIL=noreply@oss-monorepo.com
SMTP_FROM_NAME=OSS Monorepo

# CORS Configuration (comma-separated for multiple values)
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_ALLOW_CREDENTIALS=true

# Rate Limit Configuration
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_DURATION=1m

# File Upload Configuration
MAX_UPLOAD_SIZE=10485760
ALLOWED_FILE_TYPES=jpg,jpeg,png,gif,pdf,docx,xlsx,mp4

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout

# Pagination Configuration
DEFAULT_PAGE_SIZE=20
MAX_PAGE_SIZE=100
```

### 2. Sử dụng các môi trường khác nhau

#### Development
```bash
# Set ENV=development hoặc không set (mặc định)
export ENV=development
make dev
```

#### Production
```bash
export ENV=production
make build-prod
```

### 3. Các file được tìm kiếm theo thứ tự:

1. `.env.{ENV}` (ví dụ: `.env.development`, `.env.production`)
2. `.env`
3. Environment variables từ hệ thống
4. Default values trong code

## Lưu ý

- **KHÔNG commit** các file `.env` chứa thông tin nhạy cảm vào git
- File `.env.example` có thể được commit để làm template
- Các biến môi trường hệ thống luôn có độ ưu tiên cao nhất

## Mapping các biến

Tất cả các biến trong file `.env` được map trực tiếp với các field trong `config.go`:

- `SERVER_HOST` → `Config.Server.Host`
- `DB_HOST` → `Config.Database.Host`
- `JWT_SECRET` → `Config.JWT.Secret`
- ... và tương tự cho tất cả các config khác

Xem file `api/internal/config/config.go` để biết đầy đủ các biến được hỗ trợ.

