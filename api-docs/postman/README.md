# 🚀 Postman Collection

این پوشه شامل فایل‌های مربوط به Postman Collection و Environment است.

## 📁 فایل‌های موجود

- **`Microblog-API.postman_collection.json`** - Collection کامل Postman
- **`Microblog-Environment.postman_environment.json`** - Environment variables
- **`POSTMAN_SETUP.md`** - راهنمای استفاده از Postman
- **`README.md`** - این فایل

## 🚀 نحوه استفاده

### 1. Import Collection
1. فایل `Microblog-API.postman_collection.json` را در Postman import کنید
2. فایل `Microblog-Environment.postman_environment.json` را import کنید
3. Environment "Microblog Environment" را انتخاب کنید

### 2. Environment Variables
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

## ✨ ویژگی‌های Collection

1. **Auto Token Management**: Token ها خودکار ذخیره می‌شوند
2. **Response Validation**: تست‌های خودکار برای هر response
3. **Environment Variables**: مدیریت آسان متغیرها
4. **Complete Workflows**: مثال‌های کامل workflow ها
5. **Organized Structure**: ساختار منظم و قابل فهم

## 🔄 Workflow پیشنهادی

1. **Register User** → دریافت token
2. **Login User** → تایید احراز هویت  
3. **Upload Media** → آپلود تصویر
4. **Create Article** → ایجاد مقاله
5. **Create Comment** → اضافه کردن نظر
6. **Rate Article** → امتیازدهی

## 📊 API Endpoints Summary

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

## 📚 منابع بیشتر

- [Postman Documentation](https://learning.postman.com/)
- [REST API Best Practices](https://restfulapi.net/)
- [JWT Authentication](https://jwt.io/introduction/)
- [Microservices Testing](https://microservices.io/patterns/testing/)
