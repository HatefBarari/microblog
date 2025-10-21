# Postman Collection Setup Guide

این راهنما نحوه استفاده از Postman Collection برای تست API های Microblog Platform را توضیح می‌دهد.

## 📥 نصب و راه‌اندازی

### 1. Import Collection

1. فایل `Microblog-API.postman_collection.json` را در Postman import کنید
2. فایل `Microblog-Environment.postman_environment.json` را به عنوان Environment import کنید
3. Environment "Microblog Environment" را انتخاب کنید

### 2. تنظیم Environment Variables

در Environment، متغیرهای زیر را تنظیم کنید:

```json
{
  "auth_base_url": "http://localhost:8081",
  "blog_base_url": "http://localhost:8082", 
  "media_base_url": "http://localhost:8083",
  "access_token": "",
  "refresh_token": "",
  "user_id": "",
  "article_id": "",
  "media_id": "",
  "comment_id": "",
  "category_id": ""
}
```

## 🚀 نحوه استفاده

### مرحله 1: راه‌اندازی سرویس‌ها

```bash
# اجرای MongoDB
docker run -d --name mongo -p 27017:27017 \
  -e MONGO_INITDB_ROOT_USERNAME=root \
  -e MONGO_INITDB_ROOT_PASSWORD=rootpass \
  mongo:7

# اجرای MailHog (برای تست ایمیل)
docker run -d --name mailhog -p 1025:1025 -p 8025:8025 \
  mailhog/mailhog

# اجرای سرویس‌ها
cd auth-service && go run cmd/server/main.go &
cd blog-service && go run cmd/server/main.go &
cd media-service && go run cmd/server/main.go &
```

### مرحله 2: تست Authentication

1. **Register User**: درخواست ثبت‌نام کاربر جدید
2. **Login User**: ورود کاربر (token ها خودکار ذخیره می‌شوند)
3. **Get Current User**: دریافت اطلاعات کاربر فعلی

### مرحله 3: تست Blog Service

1. **Create Article**: ایجاد مقاله جدید
2. **Get Article**: دریافت مقاله
3. **List Articles**: لیست مقالات
4. **Create Comment**: ایجاد نظر
5. **Rate Article**: امتیازدهی

### مرحله 4: تست Media Service

1. **Upload Media**: آپلود فایل (فقط admin/manager/author)
2. **Get Media Info**: دریافت اطلاعات فایل
3. **List User Media**: لیست فایل‌های کاربر
4. **Serve Media**: سرو فایل

## 🔄 Workflow Examples

### Complete User Registration Flow

1. **Register User** → دریافت access_token و refresh_token
2. **Verify Email** → کلیک روی لینک در ایمیل (MailHog UI: http://localhost:8025)
3. **Login User** → ورود مجدد
4. **Get Current User** → تایید احراز هویت

### Content Creation Flow

1. **Upload Media** → آپلود تصویر
2. **Create Article** → ایجاد مقاله با cover_url
3. **Get Article** → مشاهده مقاله
4. **Create Comment** → اضافه کردن نظر
5. **Rate Article** → امتیازدهی

## 🛠️ ویژگی‌های Collection

### Auto Token Management

Collection به صورت خودکار token ها را از response ها استخراج و ذخیره می‌کند:

```javascript
// Auto-extract tokens from login response
if (pm.response && pm.response.json && pm.response.data) {
    const data = pm.response.json.data;
    if (data.access_token) {
        pm.collectionVariables.set('access_token', data.access_token);
    }
    if (data.refresh_token) {
        pm.collectionVariables.set('refresh_token', data.refresh_token);
    }
}
```

### Response Validation

هر درخواست شامل تست‌های خودکار است:

```javascript
// Test response status
pm.test('Status code is successful', function () {
    pm.expect(pm.response.code).to.be.oneOf([200, 201, 204]);
});

// Test response format
pm.test('Response has success field', function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('success');
});
```

### Environment Variables

متغیرهای محیطی برای مدیریت آسان:

- `{{auth_base_url}}` - آدرس Auth Service
- `{{blog_base_url}}` - آدرس Blog Service  
- `{{media_base_url}}` - آدرس Media Service
- `{{access_token}}` - JWT Access Token
- `{{refresh_token}}` - JWT Refresh Token

## 📋 API Endpoints Summary

### Auth Service (Port: 8081)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/register` | ثبت‌نام کاربر | ❌ |
| POST | `/login` | ورود کاربر | ❌ |
| POST | `/auth/refresh` | تازه‌سازی token | ❌ |
| GET | `/verify` | تایید ایمیل | ❌ |
| POST | `/forgot-password` | فراموشی رمز عبور | ❌ |
| POST | `/reset-password` | بازیابی رمز عبور | ❌ |
| POST | `/api/v1/resend-verification` | ارسال مجدد ایمیل تایید | ✅ |
| GET | `/api/v1/me` | اطلاعات کاربر فعلی | ✅ |

### Blog Service (Port: 8082)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/articles` | ایجاد مقاله | ✅ |
| GET | `/api/v1/articles/{slug}` | دریافت مقاله | ❌ |
| GET | `/api/v1/articles` | لیست مقالات | ❌ |
| PUT | `/api/v1/articles/{id}` | ویرایش مقاله | ✅ |
| DELETE | `/api/v1/articles/{id}` | حذف مقاله | ✅ |
| GET | `/api/v1/categories/tree` | درخت دسته‌بندی‌ها | ❌ |
| POST | `/api/v1/categories` | ایجاد دسته‌بندی | ✅ (admin/manager) |
| GET | `/api/v1/articles/{id}/comments` | نظرات مقاله | ❌ |
| POST | `/api/v1/articles/{id}/comments` | ایجاد نظر | ✅ |
| PUT | `/api/v1/comments/{id}/status` | تایید/رد نظر | ✅ (admin/manager) |
| POST | `/api/v1/articles/{id}/rate` | امتیازدهی | ✅ |
| DELETE | `/api/v1/articles/{id}/rate` | حذف امتیاز | ✅ |

### Media Service (Port: 8083)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/media/upload` | آپلود فایل | ✅ (admin/manager/author) |
| GET | `/api/v1/media/{id}` | اطلاعات فایل | ❌ |
| GET | `/api/v1/media` | لیست فایل‌های کاربر | ✅ |
| DELETE | `/api/v1/media/{id}` | حذف فایل | ✅ |
| GET | `/media/{filename}` | سرو فایل | ❌ |

## 🔧 تنظیمات پیشرفته

### Custom Headers

برای درخواست‌های خاص، header های زیر را اضافه کنید:

```javascript
// Content-Type for JSON
"Content-Type": "application/json"

// Authorization for protected endpoints
"Authorization": "Bearer {{access_token}}"

// Custom headers
"X-Request-ID": "unique-request-id"
```

### File Upload

برای آپلود فایل:

1. Method را روی `POST` تنظیم کنید
2. Body را روی `form-data` تنظیم کنید
3. Key را `file` و Type را `File` انتخاب کنید
4. فایل مورد نظر را انتخاب کنید

### Query Parameters

برای فیلتر کردن نتایج:

```
GET /api/v1/articles?page=1&page_size=10&status=approved&category_id=cat123
```

## 🐛 عیب‌یابی

### مشکلات رایج

1. **Connection Refused**: سرویس‌ها در حال اجرا نیستند
2. **401 Unauthorized**: Token منقضی شده یا نامعتبر
3. **403 Forbidden**: دسترسی کافی ندارید
4. **404 Not Found**: Endpoint یا resource یافت نشد

### لاگ‌ها

```bash
# مشاهده لاگ‌های سرویس‌ها
tail -f auth-service/logs/auth.log
tail -f blog-service/logs/blog.log  
tail -f media-service/logs/media.log

# مشاهده ایمیل‌ها در MailHog
open http://localhost:8025
```

### تست Connectivity

```bash
# تست اتصال به سرویس‌ها
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

## 📚 منابع بیشتر

- [Postman Documentation](https://learning.postman.com/)
- [REST API Best Practices](https://restfulapi.net/)
- [JWT Authentication](https://jwt.io/introduction/)
- [Microservices Testing](https://microservices.io/patterns/testing/)

---

**نکته**: این Collection برای تست و توسعه طراحی شده است. برای production، تنظیمات امنیتی مناسب را اعمال کنید.
