# Swagger/OpenAPI Documentation Setup

این راهنما نحوه استفاده از Swagger documentation برای Microblog Platform API را توضیح می‌دهد.

## 📋 فایل‌های موجود

- **`swagger.yaml`** - فایل OpenAPI 3.0 کامل
- **`SWAGGER_SETUP.md`** - این راهنما

## 🚀 نحوه استفاده

### 1. Swagger UI Online

1. به [Swagger Editor](https://editor.swagger.io/) بروید
2. فایل `swagger.yaml` را کپی و paste کنید
3. از "Try it out" برای تست API ها استفاده کنید

### 2. Swagger UI Local

#### نصب Swagger UI

```bash
# با Docker
docker run -p 8080:8080 -e SWAGGER_JSON=/swagger.yaml -v $(pwd):/swagger swaggerapi/swagger-ui

# یا با npm
npm install -g swagger-ui-serve
swagger-ui-serve swagger.yaml
```

#### دسترسی به UI

```
http://localhost:8080
```

### 3. Redoc

```bash
# نصب Redoc
npm install -g redoc-cli

# تولید documentation
redoc-cli serve swagger.yaml
```

### 4. Postman Import

1. فایل `swagger.yaml` را در Postman import کنید
2. Collection خودکار ایجاد می‌شود

## 🔧 تنظیمات پیشرفته

### Environment Variables

برای تست با سرویس‌های واقعی، متغیرهای زیر را تنظیم کنید:

```yaml
servers:
  - url: http://localhost:8081
    description: Auth Service (Development)
  - url: http://localhost:8082
    description: Blog Service (Development)
  - url: http://localhost:8083
    description: Media Service (Development)
```

### Authentication Setup

1. در Swagger UI، روی "Authorize" کلیک کنید
2. JWT token را وارد کنید:
   ```
   Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
   ```

### File Upload Testing

برای تست آپلود فایل:

1. Endpoint `/api/v1/media/upload` را انتخاب کنید
2. "Try it out" را کلیک کنید
3. فایل را انتخاب کنید
4. "Execute" را کلیک کنید

## 📚 API Documentation Features

### 1. Complete Service Coverage

- **Auth Service**: تمام endpoint های احراز هویت
- **Blog Service**: مدیریت مقالات، دسته‌بندی‌ها، نظرات
- **Media Service**: آپلود و مدیریت فایل‌ها

### 2. Detailed Schemas

```yaml
Article:
  type: object
  properties:
    id: string
    title: string
    content: string
    status: enum
    # ... تمام فیلدها
```

### 3. Request/Response Examples

```yaml
requestBody:
  content:
    application/json:
      schema:
        type: object
        example:
          title: "My First Article"
          content: "This is the content..."
```

### 4. Error Handling

```yaml
responses:
  "400":
    description: Bad request
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/Error"
```

## 🎯 API Endpoints Summary

### Auth Service (Port: 8081)

| Method | Endpoint                      | Description            | Auth Required |
| ------ | ----------------------------- | ---------------------- | ------------- |
| POST   | `/register`                   | ثبت‌نام کاربر          | ❌            |
| POST   | `/login`                      | ورود کاربر             | ❌            |
| POST   | `/auth/refresh`               | تازه‌سازی token        | ❌            |
| GET    | `/verify`                     | تایید ایمیل            | ❌            |
| POST   | `/forgot-password`            | فراموشی رمز عبور       | ❌            |
| POST   | `/reset-password`             | بازیابی رمز عبور       | ❌            |
| POST   | `/api/v1/resend-verification` | ارسال مجدد ایمیل تایید | ✅            |
| GET    | `/api/v1/me`                  | اطلاعات کاربر فعلی     | ✅            |

### Blog Service (Port: 8082)

| Method | Endpoint                         | Description       | Auth Required      |
| ------ | -------------------------------- | ----------------- | ------------------ |
| POST   | `/api/v1/articles`               | ایجاد مقاله       | ✅                 |
| GET    | `/api/v1/articles/{slug}`        | دریافت مقاله      | ❌                 |
| GET    | `/api/v1/articles`               | لیست مقالات       | ❌                 |
| PUT    | `/api/v1/articles/{id}`          | ویرایش مقاله      | ✅                 |
| DELETE | `/api/v1/articles/{id}`          | حذف مقاله         | ✅                 |
| GET    | `/api/v1/categories/tree`        | درخت دسته‌بندی‌ها | ❌                 |
| POST   | `/api/v1/categories`             | ایجاد دسته‌بندی   | ✅ (admin/manager) |
| GET    | `/api/v1/articles/{id}/comments` | نظرات مقاله       | ❌                 |
| POST   | `/api/v1/articles/{id}/comments` | ایجاد نظر         | ✅                 |
| PUT    | `/api/v1/comments/{id}/status`   | تایید/رد نظر      | ✅ (admin/manager) |
| POST   | `/api/v1/articles/{id}/rate`     | امتیازدهی         | ✅                 |
| DELETE | `/api/v1/articles/{id}/rate`     | حذف امتیاز        | ✅                 |

### Media Service (Port: 8083)

| Method | Endpoint               | Description         | Auth Required             |
| ------ | ---------------------- | ------------------- | ------------------------- |
| POST   | `/api/v1/media/upload` | آپلود فایل          | ✅ (admin/manager/author) |
| GET    | `/api/v1/media/{id}`   | اطلاعات فایل        | ❌                        |
| GET    | `/api/v1/media`        | لیست فایل‌های کاربر | ✅                        |
| DELETE | `/api/v1/media/{id}`   | حذف فایل            | ✅                        |
| GET    | `/media/{filename}`    | سرو فایل            | ❌                        |

## 🔐 Authentication & Authorization

### JWT Token Structure

```json
{
  "uid": "user123",
  "role": "author",
  "exp": 1640995200,
  "iat": 1640908800
}
```

### User Roles & Permissions

| Role        | Permissions                                          |
| ----------- | ---------------------------------------------------- |
| **Guest**   | Read articles, categories, comments                  |
| **User**    | Guest + Create comments, rate articles               |
| **Author**  | User + Create/edit/delete own articles, upload media |
| **Manager** | Author + Manage categories, moderate comments        |
| **Admin**   | All permissions                                      |

### Security Headers

```http
Authorization: Bearer <jwt-token>
Content-Type: application/json
```

## 📊 Response Formats

### Success Response

```json
{
  "success": true,
  "data": {
    "id": "article123",
    "title": "My Article",
    "content": "..."
  }
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": 400,
    "message": "Bad request",
    "details": "Validation failed"
  }
}
```

## 🧪 Testing Workflows

### 1. User Registration Flow

1. **POST** `/register` → دریافت tokens
2. **GET** `/verify?token=...` → تایید ایمیل
3. **POST** `/login` → ورود مجدد
4. **GET** `/api/v1/me` → تایید احراز هویت

### 2. Content Creation Flow

1. **POST** `/api/v1/media/upload` → آپلود تصویر
2. **POST** `/api/v1/articles` → ایجاد مقاله
3. **GET** `/api/v1/articles/{slug}` → مشاهده مقاله
4. **POST** `/api/v1/articles/{id}/comments` → اضافه کردن نظر
5. **POST** `/api/v1/articles/{id}/rate` → امتیازدهی

### 3. Content Management Flow

1. **GET** `/api/v1/articles` → لیست مقالات
2. **PUT** `/api/v1/articles/{id}` → ویرایش مقاله
3. **GET** `/api/v1/categories/tree` → مشاهده دسته‌بندی‌ها
4. **POST** `/api/v1/categories` → ایجاد دسته‌بندی جدید

## 🔧 Customization

### Adding New Endpoints

```yaml
/api/v1/new-endpoint:
  post:
    tags:
      - New Feature
    summary: New endpoint description
    requestBody:
      content:
        application/json:
          schema:
            type: object
            properties:
              field:
                type: string
    responses:
      "200":
        description: Success
```

### Adding New Schemas

```yaml
components:
  schemas:
    NewModel:
      type: object
      properties:
        id:
          type: string
        name:
          type: string
```

### Environment Configuration

```yaml
servers:
  - url: https://api.microblog.com
    description: Production
  - url: https://staging-api.microblog.com
    description: Staging
  - url: http://localhost:8081
    description: Development
```

## 📈 Monitoring & Analytics

### Health Checks

```bash
# Auth Service
curl http://localhost:8081/health

# Blog Service
curl http://localhost:8082/health

# Media Service
curl http://localhost:8083/health
```

### API Metrics

- Response times
- Error rates
- Request volumes
- User activity

## 🚀 Deployment

### Production Setup

1. **Load Balancer**: برای توزیع ترافیک
2. **SSL/TLS**: برای امنیت
3. **Rate Limiting**: برای جلوگیری از abuse
4. **Monitoring**: برای نظارت بر عملکرد

### Docker Deployment

```yaml
version: "3.8"
services:
  swagger-ui:
    image: swaggerapi/swagger-ui
    ports:
      - "8080:8080"
    environment:
      - SWAGGER_JSON=/swagger.yaml
    volumes:
      - ./swagger.yaml:/swagger.yaml
```

## 📚 منابع بیشتر

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Redoc](https://redoc.ly/)
- [API Design Best Practices](https://swagger.io/resources/articles/best-practices-in-api-design/)

---

**نکته**: این documentation برای development و testing طراحی شده است. برای production، تنظیمات امنیتی مناسب را اعمال کنید.
