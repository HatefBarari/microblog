# Microblog Platform

یک پلتفرم وبلاگ مدرن با معماری microservice که شامل مدیریت کاربران، محتوا و رسانه است.

## 🚀 ویژگی‌های کلیدی

- **معماری Microservice**: سرویس‌های مستقل و قابل مقیاس
- **احراز هویت JWT**: سیستم امن احراز هویت و مدیریت نقش‌ها
- **مدیریت محتوا**: سیستم کامل مدیریت مقالات، دسته‌بندی‌ها و نظرات
- **مدیریت رسانه**: آپلود و مدیریت فایل‌های تصویری
- **امتیازدهی**: سیستم امتیازدهی و نظرات
- **API RESTful**: API های استاندارد و مستند
- **تست‌های جامع**: تست‌های unit، integration و end-to-end

## 🏗️ معماری سیستم

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Auth Service  │    │   Blog Service  │    │  Media Service  │
│   Port: 8081    │    │   Port: 8082    │    │   Port: 8083    │
│                 │    │                 │    │                 │
│ • Registration  │    │ • Articles      │    │ • File Upload   │
│ • Login         │    │ • Categories    │    │ • File Storage  │
│ • JWT Tokens    │    │ • Comments      │    │ • Thumbnails    │
│ • User Roles    │    │ • Ratings       │    │ • File Serving  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Shared Lib    │
                    │                 │
                    │ • JWT Utils     │
                    │ • HTTP Utils    │
                    │ • MongoDB       │
                    │ • Validation    │
                    └─────────────────┘
                                 │
                    ┌─────────────────┐
                    │    MongoDB      │
                    │   Port: 27017   │
                    └─────────────────┘
```

## 📦 سرویس‌ها

### 🔐 Auth Service (Port: 8081)

- **وظیفه**: احراز هویت و مدیریت کاربران
- **ویژگی‌ها**:
  - ثبت‌نام و ورود کاربران
  - تولید و اعتبارسنجی JWT tokens
  - مدیریت نقش‌های کاربری (guest, user, manager, admin)
  - ارسال ایمیل تایید
  - رمزگذاری امن رمز عبور

### 📝 Blog Service (Port: 8082)

- **وظیفه**: مدیریت محتوا و وبلاگ
- **ویژگی‌ها**:
  - مدیریت مقالات (ایجاد، ویرایش، حذف)
  - دسته‌بندی‌های درختی
  - سیستم نظرات با تایید ادمین
  - امتیازدهی 1 تا 5 ستاره
  - جستجو و فیلتر پیشرفته
  - آمار بازدید و امتیاز

### 🖼️ Media Service (Port: 8083)

- **وظیفه**: مدیریت فایل‌های رسانه
- **ویژگی‌ها**:
  - آپلود فایل‌های تصویری
  - تولید thumbnail خودکار
  - ذخیره‌سازی فایل‌ها
  - کنترل دسترسی (فقط admin/manager/author)
  - مدیریت metadata فایل‌ها

## 🛠️ تکنولوژی‌ها

### Backend

- **Go 1.25+**: زبان برنامه‌نویسی اصلی
- **Echo Framework**: HTTP framework
- **MongoDB**: پایگاه داده NoSQL
- **JWT**: احراز هویت
- **bcrypt**: رمزگذاری رمز عبور

### Infrastructure

- **Docker**: Containerization
- **Docker Compose**: Orchestration
- **MongoDB**: Database
- **MailHog**: Email testing

### Testing

- **Testify**: Testing framework
- **Mock**: Mock objects
- **Coverage**: Test coverage

## 🚀 نصب و اجرا

### پیش‌نیازها

- Go 1.25+
- Docker & Docker Compose
- Git

### اجرای سریع

```bash
# کلون کردن پروژه
git clone <repository-url>
cd microblog

# اجرای تمام سرویس‌ها
docker-compose up -d

# یا اجرای دستی
make run
```

### اجرای دستی

```bash
# 1. اجرای MongoDB
docker run -d --name mongo -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=root \
  -e MONGO_INITDB_ROOT_PASSWORD=rootpass \
  mongo:7

# 2. اجرای MailHog (برای تست ایمیل)
docker run -d --name mailhog -p 1025:1025 -p 8025:8025 \
  mailhog/mailhog

# 3. اجرای سرویس‌ها
cd auth-service && go run cmd/server/main.go &
cd blog-service && go run cmd/server/main.go &
cd media-service && go run cmd/server/main.go &
```

## 📋 API Documentation

### Auth Service (Port: 8081)

#### ثبت‌نام

```bash
POST /api/v1/auth/register
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### ورود

```bash
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### تازه‌سازی Token

```bash
POST /api/v1/auth/refresh
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Blog Service (Port: 8082)

#### ایجاد مقاله

```bash
POST /api/v1/articles
Authorization: Bearer <token>
{
  "title": "عنوان مقاله",
  "content": "محتوای مقاله...",
  "summary": "خلاصه مقاله",
  "category_id": "cat123",
  "tags": ["تگ1", "تگ2"]
}
```

#### دریافت مقاله

```bash
GET /api/v1/articles/{slug}
```

#### لیست مقالات

```bash
GET /api/v1/articles?page=1&page_size=10
```

#### ایجاد نظر

```bash
POST /api/v1/articles/{id}/comments
Authorization: Bearer <token>
{
  "content": "نظر من"
}
```

#### امتیازدهی

```bash
POST /api/v1/articles/{id}/rate
Authorization: Bearer <token>
{
  "stars": 5
}
```

### Media Service (Port: 8083)

#### آپلود فایل

```bash
POST /api/v1/media/upload
Authorization: Bearer <token> (admin/manager/author)
Content-Type: multipart/form-data
Form Data: file
```

#### دریافت اطلاعات فایل

```bash
GET /api/v1/media/{id}
```

#### لیست فایل‌های کاربر

```bash
GET /api/v1/media?page=1&page_size=20
```

#### حذف فایل

```bash
DELETE /api/v1/media/{id}
Authorization: Bearer <token>
```

## 🧪 تست

### اجرای تست‌ها

```bash
# تست تمام سرویس‌ها
make test

# تست سرویس خاص
cd auth-service && go test ./...
cd blog-service && go test ./...
cd media-service && go test ./...

# تست با coverage
go test -cover ./...
```

### تست‌های موجود

- **Unit Tests**: تست usecase ها و repository ها
- **Integration Tests**: تست HTTP handlers
- **Mock Tests**: استفاده از mock objects
- **End-to-End Tests**: تست کامل چرخه زندگی

## 📁 ساختار پروژه

```
microblog/
├── auth-service/                 # سرویس احراز هویت
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/              # مدل‌های دامنه
│   │   ├── infrastructure/     # زیرساخت
│   │   ├── presenter/          # HTTP handlers
│   │   ├── repository/         # Repository ها
│   │   └── usecase/            # Business logic
│   ├── tests/                  # تست‌ها
│   ├── configs/                # تنظیمات
│   ├── deployments/           # Docker files
│   └── README.md
├── blog-service/                # سرویس وبلاگ
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/             # مدل‌های دامنه
│   │   ├── infrastructure/     # زیرساخت
│   │   ├── presenter/         # HTTP handlers
│   │   ├── repository/        # Repository ها
│   │   └── usecase/           # Business logic
│   ├── tests/                 # تست‌ها
│   ├── configs/               # تنظیمات
│   ├── deployments/          # Docker files
│   └── README.md
├── media-service/              # سرویس رسانه
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── domain/            # مدل‌های دامنه
│   │   ├── infrastructure/    # زیرساخت
│   │   ├── presenter/        # HTTP handlers
│   │   ├── repository/       # Repository ها
│   │   └── usecase/          # Business logic
│   ├── tests/                # تست‌ها
│   ├── configs/              # تنظیمات
│   ├── deployments/         # Docker files
│   └── README.md
├── shared/                    # کدهای مشترک
│   ├── pkg/
│   │   ├── auth/             # JWT utilities
│   │   ├── httputil/         # HTTP utilities
│   │   ├── mongo/            # MongoDB client
│   │   ├── validator/        # Validation
│   │   ├── email/            # Email sender
│   │   └── logger/           # Logger
│   └── tests/
├── deployments/              # Docker Compose
│   ├── docker-compose.yml
│   └── Makefile
└── README.md                 # این فایل
```

## 🔧 تنظیمات

### متغیرهای محیطی

```bash
# Database
DATABASE_URI=mongodb://root:rootpass@localhost:27017

# JWT Secrets
JWT_ACCESS_SECRET=your-access-secret
JWT_REFRESH_SECRET=your-refresh-secret

# Email
EMAIL_FROM=noreply@microblog.com
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-email@gmail.com
EMAIL_SMTP_PASSWORD=your-app-password

# Storage
STORAGE_BASE_PATH=./uploads
STORAGE_BASE_URL=http://localhost:8083

# Logging
LOG_LEVEL=info
```

### Docker Compose

```yaml
version: "3.9"
services:
  mongo:
    image: mongo:7
    restart: always
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: rootpass
    ports:
      - "27017:27017"
    volumes:
      - mongo_data:/data/db

  mailhog:
    image: mailhog/mailhog:latest
    ports:
      - "1025:1025" # smtp
      - "8025:8025" # web ui

  auth-service:
    build: ./auth-service
    ports:
      - "8081:8081"
    environment:
      - DATABASE_URI=mongodb://root:rootpass@mongo:27017
      - JWT_ACCESS_SECRET=your-secret
      - JWT_REFRESH_SECRET=your-refresh-secret

  blog-service:
    build: ./blog-service
    ports:
      - "8082:8082"
    environment:
      - DATABASE_URI=mongodb://root:rootpass@mongo:27017
      - JWT_ACCESS_SECRET=your-secret

  media-service:
    build: ./media-service
    ports:
      - "8083:8083"
    environment:
      - DATABASE_URI=mongodb://root:rootpass@mongo:27017
      - STORAGE_BASE_URL=http://localhost:8083
    volumes:
      - ./uploads:/app/uploads

volumes:
  mongo_data:
```

## 🔒 امنیت

### احراز هویت

- JWT tokens با secret keys
- Access token (15 دقیقه) + Refresh token (24 ساعت)
- رمزگذاری رمز عبور با bcrypt

### کنترل دسترسی

- **Guest**: دسترسی محدود
- **User**: کاربر عادی
- **Manager**: مدیریت محتوا
- **Admin**: دسترسی کامل

### اعتبارسنجی

- اعتبارسنجی ورودی‌ها
- محدودیت اندازه فایل
- محدودیت نوع فایل
- اعتبارسنجی JWT tokens

## 📊 مانیتورینگ

### Health Checks

```bash
# Auth Service
curl http://localhost:8081/health

# Blog Service
curl http://localhost:8082/health

# Media Service
curl http://localhost:8083/health
```

### لاگ‌ها

```bash
# مشاهده لاگ‌های real-time
tail -f auth-service/logs/auth.log
tail -f blog-service/logs/blog.log
tail -f media-service/logs/media.log

# جستجو در لاگ‌ها
grep "ERROR" */logs/*.log
grep "login" */logs/*.log
```

## 🚀 توسعه

### اضافه کردن سرویس جدید

1. ایجاد ساختار سرویس:

```bash
mkdir new-service
cd new-service
go mod init github.com/HatefBarari/microblog-new-service
```

2. کپی کردن ساختار از سرویس موجود
3. به‌روزرسانی docker-compose.yml
4. اضافه کردن به Makefile

### اضافه کردن فیلد جدید

1. ویرایش domain models
2. به‌روزرسانی DTO ها
3. به‌روزرسانی usecase ها
4. به‌روزرسانی repository ها
5. به‌روزرسانی HTTP handlers
6. اضافه کردن تست‌ها

## 🐛 عیب‌یابی

### مشکلات رایج

1. **خطای اتصال به MongoDB**:

   - بررسی اتصال به MongoDB
   - بررسی تنظیمات DATABASE_URI

2. **خطای JWT**:

   - بررسی JWT secrets
   - بررسی اعتبار token

3. **خطای آپلود فایل**:
   - بررسی مجوزهای پوشه uploads
   - بررسی محدودیت‌های فایل

### لاگ‌های مفید

```bash
# مشاهده لاگ‌های تمام سرویس‌ها
tail -f */logs/*.log

# جستجو در لاگ‌ها
grep "ERROR" */logs/*.log
grep "database" */logs/*.log
grep "auth" */logs/*.log
```

## 📈 عملکرد

### بهینه‌سازی

- استفاده از connection pooling برای MongoDB
- Cache برای JWT tokens
- Compression برای HTTP responses
- Optimized database queries

### مقیاس‌پذیری

- Horizontal scaling با load balancer
- Database sharding
- CDN برای فایل‌های رسانه
- Microservice architecture

## 🤝 مشارکت

### راهنمای مشارکت

1. Fork کردن پروژه
2. ایجاد feature branch
3. نوشتن تست‌ها
4. Commit کردن تغییرات
5. Push کردن branch
6. ایجاد Pull Request

### استانداردهای کد

- Go coding standards
- Test coverage > 80%
- Documentation برای API ها
- Clean code principles

## 📄 مجوز

این پروژه تحت مجوز MIT منتشر شده است.

## 📞 پشتیبانی

برای سوالات و پشتیبانی:

- ایجاد Issue در GitHub
- ایمیل: support@microblog.com
- مستندات: [Wiki](https://github.com/username/microblog/wiki)

## 🔄 تغییرات

### نسخه 1.0.0

- پیاده‌سازی اولیه
- Auth Service
- Blog Service
- Media Service
- تست‌های جامع
- مستندات کامل

---

**نکته**: این پروژه در حال توسعه است و ممکن است تغییرات breaking در نسخه‌های آینده داشته باشد.
