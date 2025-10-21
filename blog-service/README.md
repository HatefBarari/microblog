# Blog Service

سرویس مدیریت محتوا و وبلاگ برای پروژه microblog که مسئول مدیریت مقالات، دسته‌بندی‌ها، نظرات و امتیازدهی است.

## ویژگی‌ها

- **مدیریت مقالات**: ایجاد، ویرایش، حذف و لیست مقالات
- **دسته‌بندی‌ها**: مدیریت دسته‌بندی‌های درختی
- **نظرات**: سیستم نظرات با تایید ادمین
- **امتیازدهی**: سیستم امتیازدهی 1 تا 5 ستاره
- **جستجو و فیلتر**: جستجو بر اساس دسته‌بندی، تگ و وضعیت
- **آمار**: شمارش بازدید و میانگین امتیاز

## API Endpoints

### مقالات

#### ایجاد مقاله
```
POST /api/v1/articles
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "عنوان مقاله",
  "content": "محتوای مقاله...",
  "summary": "خلاصه مقاله",
  "category_id": "cat123",
  "tags": ["تگ1", "تگ2"],
  "cover_url": "http://localhost:8083/uploads/cover.jpg"
}

Response:
{
  "success": true,
  "data": {
    "id": "article123",
    "author_id": "user123",
    "title": "عنوان مقاله",
    "slug": "عنوان-مقاله",
    "summary": "خلاصه مقاله",
    "content": "محتوای مقاله...",
    "cover_url": "http://localhost:8083/uploads/cover.jpg",
    "status": "draft",
    "category_id": "cat123",
    "tags": ["تگ1", "تگ2"],
    "view_count": 0,
    "rating_avg": 0,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### دریافت مقاله
```
GET /api/v1/articles/{slug}

Response:
{
  "success": true,
  "data": {
    "id": "article123",
    "author_id": "user123",
    "title": "عنوان مقاله",
    "slug": "عنوان-مقاله",
    "summary": "خلاصه مقاله",
    "content": "محتوای مقاله...",
    "cover_url": "http://localhost:8083/uploads/cover.jpg",
    "status": "approved",
    "category_id": "cat123",
    "tags": ["تگ1", "تگ2"],
    "view_count": 15,
    "rating_avg": 4.5,
    "created_at": "2024-01-01T00:00:00Z",
    "published_at": "2024-01-01T12:00:00Z"
  }
}
```

#### لیست مقالات
```
GET /api/v1/articles?page=1&page_size=10&category_id=cat123&status=approved

Response:
{
  "success": true,
  "data": {
    "list": [
      {
        "id": "article123",
        "title": "عنوان مقاله",
        "slug": "عنوان-مقاله",
        "summary": "خلاصه مقاله",
        "status": "approved",
        "view_count": 15,
        "rating_avg": 4.5,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1
  }
}
```

#### ویرایش مقاله
```
PUT /api/v1/articles/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "عنوان جدید",
  "content": "محتوای جدید...",
  "summary": "خلاصه جدید",
  "category_id": "cat456",
  "tags": ["تگ جدید"]
}
```

#### حذف مقاله
```
DELETE /api/v1/articles/{id}
Authorization: Bearer <token>

Response: 204 No Content
```

### دسته‌بندی‌ها

#### لیست درخت دسته‌بندی‌ها
```
GET /api/v1/categories/tree

Response:
{
  "success": true,
  "data": [
    {
      "id": "cat1",
      "name": "تکنولوژی",
      "slug": "تکنولوژی",
      "children": [
        {
          "id": "cat2",
          "name": "برنامه‌نویسی",
          "slug": "برنامه-نویسی",
          "parent_id": "cat1"
        }
      ]
    }
  ]
}
```

#### ایجاد دسته‌بندی
```
POST /api/v1/categories
Authorization: Bearer <token> (admin/manager)
Content-Type: application/json

{
  "name": "دسته‌بندی جدید",
  "parent_id": "cat123"
}
```

### نظرات

#### لیست نظرات مقاله
```
GET /api/v1/articles/{id}/comments

Response:
{
  "success": true,
  "data": [
    {
      "id": "comment123",
      "article_id": "article123",
      "author_id": "user123",
      "content": "نظر کاربر",
      "status": "approved",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### ایجاد نظر
```
POST /api/v1/articles/{id}/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "نظر من",
  "parent_id": "comment123"  // اختیاری - برای پاسخ
}
```

#### تایید/رد نظر
```
PUT /api/v1/comments/{id}/status
Authorization: Bearer <token> (admin/manager)
Content-Type: application/json

{
  "status": "approved"  // یا "rejected"
}
```

### امتیازدهی

#### امتیازدهی به مقاله
```
POST /api/v1/articles/{id}/rate
Authorization: Bearer <token>
Content-Type: application/json

{
  "stars": 5
}
```

#### حذف امتیاز
```
DELETE /api/v1/articles/{id}/rate
Authorization: Bearer <token>

Response: 204 No Content
```

## مدل‌های داده

### Article
```go
type Article struct {
    ID          string        `bson:"_id,omitempty"`
    AuthorID    string        `bson:"author_id"`
    Title       string        `bson:"title"`
    Slug        string        `bson:"slug"`
    Summary     string        `bson:"summary"`
    Content     string        `bson:"content"`
    CoverURL    string        `bson:"cover_url"`
    Status      ArticleStatus `bson:"status"`
    CategoryID  string        `bson:"category_id"`
    Tags        []string      `bson:"tags"`
    ViewCount   int64         `bson:"view_count"`
    RatingAvg   float64       `bson:"rating_avg"`
    CreatedAt   time.Time     `bson:"created_at"`
    UpdatedAt   time.Time     `bson:"updated_at"`
    PublishedAt *time.Time    `bson:"published_at"`
}
```

### Category
```go
type Category struct {
    ID       string `bson:"_id,omitempty"`
    Name     string `bson:"name"`
    Slug     string `bson:"slug"`
    ParentID string `bson:"parent_id"`
}
```

### Comment
```go
type Comment struct {
    ID        string        `bson:"_id,omitempty"`
    ArticleID string        `bson:"article_id"`
    ParentID  string        `bson:"parent_id"`
    AuthorID  string        `bson:"author_id"`
    Content   string        `bson:"content"`
    Status    CommentStatus `bson:"status"`
    CreatedAt time.Time     `bson:"created_at"`
}
```

### Rating
```go
type Rating struct {
    ID        string    `bson:"_id,omitempty"`
    UserID    string    `bson:"user_id"`
    TargetID  string    `bson:"target_id"`
    Type      string    `bson:"type"`
    Stars     int       `bson:"stars"`
    CreatedAt time.Time `bson:"created_at"`
}
```

## وضعیت‌های مقاله

- **draft**: پیش‌نویس
- **pending**: در انتظار تایید
- **approved**: تایید شده
- **rejected**: رد شده
- **archived**: آرشیو شده

## وضعیت‌های نظر

- **pending**: در انتظار تایید
- **approved**: تایید شده
- **rejected**: رد شده

## تنظیمات

### متغیرهای محیطی

```bash
# Server
SERVER_PORT=8082
SERVER_HOST=0.0.0.0

# Database
DATABASE_URI=mongodb://root:rootpass@localhost:27017
DATABASE_NAME=microblog_blog

# JWT
JWT_ACCESS_SECRET=your-access-secret
JWT_REFRESH_SECRET=your-refresh-secret

# Logging
LOG_LEVEL=info
LOG_FILE=logs/blog.log
```

### فایل تنظیمات (config.yaml)

```yaml
server:
  port: "8082"
  host: "0.0.0.0"

database:
  uri: "mongodb://root:rootpass@localhost:27017"
  database: "microblog_blog"

jwt:
  access_secret: "your-access-secret"
  refresh_secret: "your-refresh-secret"

log:
  level: "info"
  file: "logs/blog.log"
```

## نصب و اجرا

### پیش‌نیازها

- Go 1.25+
- MongoDB 7+
- Auth Service (برای احراز هویت)

### اجرای محلی

```bash
# کلون کردن پروژه
git clone <repository-url>
cd microblog/blog-service

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
docker build -f deployments/Dockerfile -t microblog-blog .

# اجرای container
docker run -d --name blog-service -p 8082:8082 \
  -e DATABASE_URI=mongodb://root:rootpass@mongo:27017 \
  -e JWT_ACCESS_SECRET=your-secret \
  microblog-blog
```

## تست

### اجرای تست‌ها

```bash
# اجرای تمام تست‌ها
go test ./...

# اجرای تست‌ها با coverage
go test -cover ./...

# اجرای تست‌های خاص
go test ./tests/article_usecase_test.go
go test ./tests/http_handler_test.go
go test ./tests/integration_test.go
```

### تست‌های موجود

- **Unit Tests**: تست usecase ها و repository ها
- **Integration Tests**: تست HTTP handlers
- **Mock Tests**: استفاده از mock objects برای تست
- **End-to-End Tests**: تست کامل چرخه زندگی مقاله

## ساختار پروژه

```
blog-service/
├── cmd/
│   └── server/
│       └── main.go              # نقطه ورود برنامه
├── configs/
│   └── config.yaml              # تنظیمات پیش‌فرض
├── deployments/
│   └── Dockerfile               # Docker configuration
├── internal/
│   ├── domain/
│   │   ├── article.go           # مدل مقاله
│   │   ├── category.go          # مدل دسته‌بندی
│   │   ├── comment.go           # مدل نظر
│   │   ├── rating.go            # مدل امتیاز
│   │   └── repository.go        # interface های repository
│   ├── infrastructure/
│   │   ├── config.go           # مدیریت تنظیمات
│   │   ├── echo_server.go      # HTTP server
│   │   └── logger.go           # مدیریت لاگ
│   ├── presenter/
│   │   └── http_handler.go     # HTTP handlers
│   ├── repository/
│   │   └── mongo_repos.go      # MongoDB repositories
│   └── usecase/
│       ├── dto.go              # Data Transfer Objects
│       ├── article_usecase.go # Business logic مقالات
│       ├── category_usecase.go # Business logic دسته‌بندی‌ها
│       ├── comment_usecase.go  # Business logic نظرات
│       └── rating_usecase.go   # Business logic امتیازها
├── tests/
│   ├── article_usecase_test.go # تست usecase مقالات
│   ├── category_usecase_test.go # تست usecase دسته‌بندی‌ها
│   ├── comment_usecase_test.go # تست usecase نظرات
│   ├── rating_usecase_test.go  # تست usecase امتیازها
│   ├── http_handler_test.go    # تست HTTP handlers
│   └── integration_test.go     # تست‌های integration
├── logs/
│   └── blog.log                # لاگ‌های اصلی
├── go.mod                      # Go modules
└── README.md                   # این فایل
```

## امنیت

### کنترل دسترسی

- **مقالات**: فقط نویسنده می‌تواند ویرایش/حذف کند
- **دسته‌بندی‌ها**: فقط admin/manager می‌تواند ایجاد کند
- **نظرات**: تایید ادمین برای نمایش
- **امتیازدهی**: هر کاربر یک امتیاز

### اعتبارسنجی

- اعتبارسنجی ورودی‌ها
- محدودیت طول محتوا
- محدودیت تعداد تگ‌ها
- اعتبارسنجی امتیاز (1-5)

## مانیتورینگ

### Health Check

```bash
curl http://localhost:8082/health
```

### لاگ‌ها

لاگ‌ها در فایل `logs/blog.log` ذخیره می‌شوند و شامل:

- درخواست‌های HTTP
- عملیات CRUD مقالات
- عملیات نظرات و امتیازدهی
- خطاهای سیستم

## توسعه

### اضافه کردن فیلد جدید به Article

1. ویرایش `article.go`:
```go
type Article struct {
    // فیلدهای موجود...
    Featured    bool      `bson:"featured"`     // فیلد جدید
    ExpiresAt   *time.Time `bson:"expires_at"`  // فیلد جدید
}
```

2. به‌روزرسانی DTO ها و usecase ها

### اضافه کردن فیلتر جدید

1. ویرایش `ListFilter`:
```go
type ListFilter struct {
    // فیلترهای موجود...
    Featured    *bool     // فیلتر جدید
    DateFrom    *time.Time // فیلتر جدید
    DateTo      *time.Time // فیلتر جدید
}
```

2. به‌روزرسانی repository و usecase

## عیب‌یابی

### مشکلات رایج

1. **خطای اتصال به MongoDB**:
   - بررسی اتصال به MongoDB
   - بررسی تنظیمات DATABASE_URI

2. **خطای احراز هویت**:
   - بررسی JWT token
   - بررسی اتصال به auth-service

3. **خطای slug تکراری**:
   - بررسی منطق تولید slug
   - بررسی unique constraint

### لاگ‌های مفید

```bash
# مشاهده لاگ‌های real-time
tail -f logs/blog.log

# جستجو در لاگ‌ها
grep "ERROR" logs/blog.log
grep "article" logs/blog.log
grep "comment" logs/blog.log
```

## مثال‌های استفاده

### ایجاد مقاله جدید

```bash
curl -X POST http://localhost:8082/api/v1/articles \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "مقاله جدید",
    "content": "محتوای مقاله...",
    "summary": "خلاصه مقاله",
    "category_id": "cat123",
    "tags": ["تکنولوژی", "برنامه‌نویسی"]
  }'
```

### دریافت مقاله

```bash
curl http://localhost:8082/api/v1/articles/عنوان-مقاله
```

### ایجاد نظر

```bash
curl -X POST http://localhost:8082/api/v1/articles/article123/comments \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "نظر من در مورد این مقاله"
  }'
```

### امتیازدهی

```bash
curl -X POST http://localhost:8082/api/v1/articles/article123/rate \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "stars": 5
  }'
```
