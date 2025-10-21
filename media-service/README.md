# Media Service

سرویس مدیریت رسانه برای پروژه microblog که مسئول آپلود، ذخیره و مدیریت فایل‌های رسانه (عکس، ویدیو، صدا) است.

## ویژگی‌ها

- **آپلود فایل**: آپلود فایل‌های رسانه با محدودیت نوع و اندازه
- **مدیریت دسترسی**: فقط ادمین، مدیر و نویسنده می‌توانند فایل آپلود کنند
- **تولید thumbnail**: تولید تصویر کوچک برای فایل‌های تصویری
- **مدیریت فایل**: حذف و لیست فایل‌های آپلود شده
- **ذخیره‌سازی**: ذخیره فایل‌ها در سیستم فایل محلی
- **پایگاه داده**: ذخیره اطلاعات فایل‌ها در MongoDB

## API Endpoints

### آپلود فایل
```
POST /api/v1/media/upload
Content-Type: multipart/form-data

Form Data:
- file: فایل آپلود شده

Headers:
- Authorization: Bearer <token>

Response:
{
  "success": true,
  "data": {
    "media": {
      "id": "media123",
      "uploader_id": "user123",
      "filename": "test.jpg",
      "original_name": "test.jpg",
      "mime_type": "image/jpeg",
      "size": 1024,
      "type": "image",
      "status": "processed",
      "url": "http://localhost:8083/uploads/test.jpg",
      "thumbnail_url": "http://localhost:8083/uploads/thumb_test.jpg",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "url": "http://localhost:8083/uploads/test.jpg"
  }
}
```

### دریافت اطلاعات فایل
```
GET /api/v1/media/{id}

Response:
{
  "success": true,
  "data": {
    "id": "media123",
    "uploader_id": "user123",
    "filename": "test.jpg",
    "original_name": "test.jpg",
    "mime_type": "image/jpeg",
    "size": 1024,
    "type": "image",
    "status": "processed",
    "url": "http://localhost:8083/uploads/test.jpg",
    "thumbnail_url": "http://localhost:8083/uploads/thumb_test.jpg",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### لیست فایل‌های کاربر
```
GET /api/v1/media?page=1&page_size=20

Response:
{
  "success": true,
  "data": {
    "media": [
      {
        "id": "media123",
        "uploader_id": "user123",
        "filename": "test.jpg",
        "original_name": "test.jpg",
        "mime_type": "image/jpeg",
        "size": 1024,
        "type": "image",
        "status": "processed",
        "url": "http://localhost:8083/uploads/test.jpg",
        "thumbnail_url": "http://localhost:8083/uploads/thumb_test.jpg",
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1,
    "page": 1,
    "size": 1
  }
}
```

### حذف فایل
```
DELETE /api/v1/media/{id}

Response: 204 No Content
```

### سرو فایل
```
GET /media/{filename}

Response: Redirect to file URL
```

## تنظیمات

### متغیرهای محیطی

```bash
# Server
SERVER_PORT=8083
SERVER_HOST=0.0.0.0

# Database
DATABASE_URI=mongodb://root:rootpass@localhost:27017
DATABASE_NAME=microblog_media

# Storage
STORAGE_BASE_PATH=./uploads
STORAGE_BASE_URL=http://localhost:8083

# Media
MEDIA_MAX_FILE_SIZE=10485760  # 10MB
MEDIA_ALLOWED_TYPES=image/jpeg,image/png,image/gif,image/webp,image/bmp

# Logging
LOG_LEVEL=info
LOG_FILE=logs/media.log
```

### فایل تنظیمات (config.yaml)

```yaml
server:
  port: "8083"
  host: "0.0.0.0"

database:
  uri: "mongodb://root:rootpass@localhost:27017"
  database: "microblog_media"

storage:
  type: "file"
  base_path: "./uploads"
  base_url: "http://localhost:8083"

media:
  max_file_size: 10485760  # 10MB
  allowed_types:
    - "image/jpeg"
    - "image/png"
    - "image/gif"
    - "image/webp"
    - "image/bmp"
  thumbnail_size: 300

log:
  level: "info"
  file: "logs/media.log"
```

## نصب و اجرا

### پیش‌نیازها

- Go 1.25+
- MongoDB 7+
- Docker (اختیاری)

### اجرای محلی

```bash
# کلون کردن پروژه
git clone <repository-url>
cd microblog/media-service

# نصب وابستگی‌ها
go mod download

# اجرای MongoDB
docker run -d --name mongo -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=rootpass mongo:7

# اجرای سرویس
go run cmd/server/main.go
```

### اجرا با Docker

```bash
# ساخت image
docker build -f deployments/Dockerfile -t microblog-media .

# اجرای container
docker run -d --name media-service -p 8083:8083 \
  -e DATABASE_URI=mongodb://root:rootpass@mongo:27017 \
  -e STORAGE_BASE_URL=http://localhost:8083 \
  microblog-media
```

## تست

### اجرای تست‌ها

```bash
# اجرای تمام تست‌ها
go test ./...

# اجرای تست‌ها با coverage
go test -cover ./...

# اجرای تست‌های خاص
go test ./tests/media_usecase_test.go
```

### تست‌های موجود

- **Unit Tests**: تست usecase ها و repository ها
- **Integration Tests**: تست HTTP handlers
- **Mock Tests**: استفاده از mock objects برای تست

## ساختار پروژه

```
media-service/
├── cmd/
│   └── server/
│       └── main.go              # نقطه ورود برنامه
├── configs/
│   └── config.yaml              # تنظیمات پیش‌فرض
├── deployments/
│   └── Dockerfile               # Docker configuration
├── internal/
│   ├── domain/
│   │   └── media.go             # مدل‌های دامنه
│   ├── infrastructure/
│   │   ├── config.go           # مدیریت تنظیمات
│   │   ├── echo_server.go      # HTTP server
│   │   ├── file_storage.go     # ذخیره‌سازی فایل
│   │   └── logger.go           # مدیریت لاگ
│   ├── presenter/
│   │   └── http_handler.go     # HTTP handlers
│   ├── repository/
│   │   └── mongo_media_repo.go # MongoDB repository
│   └── usecase/
│       ├── dto.go              # Data Transfer Objects
│       └── media_usecase.go    # Business logic
├── tests/
│   ├── media_usecase_test.go   # تست usecase
│   └── http_handler_test.go    # تست HTTP handlers
├── go.mod                       # Go modules
└── README.md                    # این فایل
```

## امنیت

### محدودیت‌های دسترسی

- فقط کاربران با نقش `admin`، `manager` یا `author` می‌توانند فایل آپلود کنند
- کاربران فقط می‌توانند فایل‌های خود را حذف کنند
- اعتبارسنجی نوع فایل و اندازه فایل

### محدودیت‌های فایل

- حداکثر اندازه فایل: 10MB (قابل تنظیم)
- انواع فایل مجاز: تصاویر (JPEG, PNG, GIF, WebP, BMP)
- تولید نام فایل منحصر به فرد برای جلوگیری از تداخل

## مانیتورینگ

### Health Check

```bash
curl http://localhost:8083/health
```

### لاگ‌ها

لاگ‌ها در فایل `logs/media.log` ذخیره می‌شوند و شامل:

- درخواست‌های HTTP
- عملیات آپلود و حذف فایل
- خطاهای سیستم
- اطلاعات عملکرد

## توسعه

### اضافه کردن نوع فایل جدید

1. ویرایش `config.yaml`:
```yaml
media:
  allowed_types:
    - "image/jpeg"
    - "image/png"
    - "video/mp4"  # اضافه کردن ویدیو
```

2. به‌روزرسانی `getMediaType` در `media_usecase.go`

### اضافه کردن storage جدید

پیاده‌سازی interface `MediaStorage`:

```go
type CustomStorage struct {
    // implementation
}

func (s *CustomStorage) Save(ctx context.Context, filename string, data []byte) (string, error) {
    // implementation
}
```

## عیب‌یابی

### مشکلات رایج

1. **خطای اتصال به MongoDB**:
   - بررسی اتصال به MongoDB
   - بررسی تنظیمات DATABASE_URI

2. **خطای آپلود فایل**:
   - بررسی مجوزهای پوشه uploads
   - بررسی محدودیت‌های فایل

3. **خطای دسترسی**:
   - بررسی JWT token
   - بررسی نقش کاربر

### لاگ‌های مفید

```bash
# مشاهده لاگ‌های real-time
tail -f logs/media.log

# جستجو در لاگ‌ها
grep "ERROR" logs/media.log
grep "upload" logs/media.log
```
