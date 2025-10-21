# Auth Service

سرویس احراز هویت و مدیریت کاربران برای پروژه microblog که مسئول ثبت‌نام، ورود، احراز هویت و مدیریت نقش‌های کاربری است.

## ویژگی‌ها

- **ثبت‌نام کاربر**: ثبت‌نام کاربران جدید با رمزگذاری رمز عبور
- **ورود کاربر**: احراز هویت با ایمیل و رمز عبور
- **JWT Tokens**: تولید و اعتبارسنجی access و refresh tokens
- **مدیریت نقش‌ها**: نقش‌های مختلف (guest, user, manager, admin)
- **تایید ایمیل**: ارسال ایمیل تایید برای کاربران جدید
- **امنیت**: رمزگذاری رمز عبور با bcrypt

## API Endpoints

### ثبت‌نام کاربر
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### ورود کاربر
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### تازه‌سازی Token
```
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

Response:
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### تایید ایمیل
```
GET /verify?token=verification_token

Response:
{
  "success": true,
  "data": {
    "message": "email verified successfully"
  }
}
```

### ارسال مجدد ایمیل تایید
```
POST /api/v1/resend-verification
Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "message": "verification email sent"
  }
}
```

### فراموشی رمز عبور
```
POST /forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}

Response:
{
  "success": true,
  "data": {
    "message": "password reset email sent"
  }
}
```

### بازیابی رمز عبور
```
POST /reset-password
Content-Type: application/json

{
  "token": "reset_token",
  "new_password": "newpassword123"
}

Response:
{
  "success": true,
  "data": {
    "message": "password reset successfully"
  }
}
```

## مدل‌های داده

### User
```go
type User struct {
    ID           string    `bson:"_id,omitempty"`
    Email        string    `bson:"email"`
    PasswordHash string    `bson:"password_hash"`
    Role         Role      `bson:"role"`
    Verified     bool      `bson:"verified"`
    CreatedAt    time.Time `bson:"created_at"`
}
```

### Role
```go
type Role string

const (
    RoleGuest   Role = "guest"
    RoleUser    Role = "user"
    RoleManager Role = "manager"
    RoleAdmin   Role = "admin"
)
```

## تنظیمات

### متغیرهای محیطی

```bash
# Server
SERVER_PORT=8081
SERVER_HOST=0.0.0.0

# Database
DATABASE_URI=mongodb://root:rootpass@localhost:27017
DATABASE_NAME=microblog_auth

# JWT
JWT_ACCESS_SECRET=your-access-secret
JWT_REFRESH_SECRET=your-refresh-secret
JWT_ACCESS_TTL_MIN=15
JWT_REFRESH_TTL_HOUR=24

# Email
EMAIL_FROM=noreply@microblog.com
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-email@gmail.com
EMAIL_SMTP_PASSWORD=your-app-password

# Logging
LOG_LEVEL=info
LOG_FILE=logs/auth.log
```

### فایل تنظیمات (config.yaml)

```yaml
server:
  port: "8081"
  host: "0.0.0.0"

database:
  uri: "mongodb://root:rootpass@localhost:27017"
  database: "microblog_auth"

jwt:
  access_secret: "your-access-secret"
  refresh_secret: "your-refresh-secret"
  access_ttl_min: 15
  refresh_ttl_hour: 24

email:
  from: "noreply@microblog.com"
  smtp_host: "smtp.gmail.com"
  smtp_port: 587
  smtp_username: "your-email@gmail.com"
  smtp_password: "your-app-password"

log:
  level: "info"
  file: "logs/auth.log"
```

## نصب و اجرا

### پیش‌نیازها

- Go 1.25+
- MongoDB 7+
- SMTP Server (برای ارسال ایمیل)

### اجرای محلی

```bash
# کلون کردن پروژه
git clone <repository-url>
cd microblog/auth-service

# نصب وابستگی‌ها
go mod download

# اجرای MongoDB
docker run -d --name mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=rootpass mongo:7

# اجرای MailHog (برای تست ایمیل)
docker run -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog

# اجرای سرویس
go run cmd/server/main.go
```

### اجرا با Docker

```bash
# ساخت image
docker build -f deployments/Dockerfile -t microblog-auth .

# اجرای container
docker run -d --name auth-service -p 8081:8081 \
  -e DATABASE_URI=mongodb://root:rootpass@mongo:27017 \
  -e JWT_ACCESS_SECRET=your-secret \
  -e JWT_REFRESH_SECRET=your-refresh-secret \
  microblog-auth
```

## تست

### اجرای تست‌ها

```bash
# اجرای تمام تست‌ها
go test ./...

# اجرای تست‌ها با coverage
go test -cover ./...

# اجرای تست‌های خاص
go test ./tests/http_handler_test.go
```

### تست‌های موجود

- **Unit Tests**: تست usecase ها و repository ها
- **Integration Tests**: تست HTTP handlers
- **Mock Tests**: استفاده از mock objects برای تست

## ساختار پروژه

```
auth-service/
├── cmd/
│   └── server/
│       └── main.go              # نقطه ورود برنامه
├── configs/
│   └── config.yaml              # تنظیمات پیش‌فرض
├── deployments/
│   └── Dockerfile               # Docker configuration
├── internal/
│   ├── domain/
│   │   ├── user.go              # مدل کاربر
│   │   ├── role.go              # نقش‌های کاربری
│   │   └── repository.go        # interface های repository
│   ├── infrastructure/
│   │   ├── config.go           # مدیریت تنظیمات
│   │   ├── echo_server.go      # HTTP server
│   │   └── logger.go            # مدیریت لاگ
│   ├── presenter/
│   │   └── http_handler.go      # HTTP handlers
│   ├── repository/
│   │   └── mongo_user_repo.go   # MongoDB repository
│   └── usecase/
│       ├── dto.go               # Data Transfer Objects
│       ├── user_usecase.go     # Business logic
│       └── email_usecase.go    # مدیریت ایمیل
├── tests/
│   ├── http_handler_test.go    # تست HTTP handlers
│   └── logs/
│       └── test.log            # لاگ‌های تست
├── logs/
│   └── auth.log                # لاگ‌های اصلی
├── go.mod                      # Go modules
└── README.md                   # این فایل
```

## امنیت

### رمزگذاری رمز عبور

- استفاده از bcrypt برای رمزگذاری رمز عبور
- Salt خودکار برای هر رمز عبور
- مقاوم در برابر حملات brute force

### JWT Tokens

- **Access Token**: کوتاه‌مدت (15 دقیقه) برای احراز هویت
- **Refresh Token**: بلندمدت (24 ساعت) برای تازه‌سازی
- امضای دیجیتال با secret key

### نقش‌های کاربری

- **Guest**: دسترسی محدود
- **User**: کاربر عادی
- **Manager**: مدیریت محتوا
- **Admin**: دسترسی کامل

## مانیتورینگ

### Health Check

```bash
curl http://localhost:8081/health
```

### لاگ‌ها

لاگ‌ها در فایل `logs/auth.log` ذخیره می‌شوند و شامل:

- درخواست‌های HTTP
- عملیات ثبت‌نام و ورود
- خطاهای احراز هویت
- ارسال ایمیل

## توسعه

### اضافه کردن نقش جدید

1. ویرایش `role.go`:
```go
const (
    RoleGuest   Role = "guest"
    RoleUser    Role = "user"
    RoleManager Role = "manager"
    RoleAdmin   Role = "admin"
    RoleModerator Role = "moderator"  // نقش جدید
)
```

2. به‌روزرسانی middleware ها و authorization logic

### اضافه کردن فیلد جدید به User

1. ویرایش `user.go`:
```go
type User struct {
    ID           string    `bson:"_id,omitempty"`
    Email        string    `bson:"email"`
    PasswordHash string    `bson:"password_hash"`
    Role         Role      `bson:"role"`
    Verified     bool      `bson:"verified"`
    FirstName    string    `bson:"first_name"`    // فیلد جدید
    LastName     string    `bson:"last_name"`      // فیلد جدید
    CreatedAt    time.Time `bson:"created_at"`
}
```

2. به‌روزرسانی DTO ها و usecase ها

## عیب‌یابی

### مشکلات رایج

1. **خطای اتصال به MongoDB**:
   - بررسی اتصال به MongoDB
   - بررسی تنظیمات DATABASE_URI

2. **خطای JWT**:
   - بررسی JWT secrets
   - بررسی اعتبار token

3. **خطای ارسال ایمیل**:
   - بررسی تنظیمات SMTP
   - بررسی مجوزهای ایمیل

### لاگ‌های مفید

```bash
# مشاهده لاگ‌های real-time
tail -f logs/auth.log

# جستجو در لاگ‌ها
grep "ERROR" logs/auth.log
grep "login" logs/auth.log
grep "register" logs/auth.log
```

## مثال‌های استفاده

### ثبت‌نام کاربر جدید

```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### ورود کاربر

```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### استفاده از Token در درخواست‌ها

```bash
curl -X GET http://localhost:8082/api/v1/articles \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
