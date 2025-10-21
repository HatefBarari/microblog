# 🔧 Swagger/OpenAPI Documentation

این پوشه شامل فایل‌های مربوط به Swagger/OpenAPI documentation است.

## 📁 فایل‌های موجود

- **`swagger.yaml`** - فایل OpenAPI 3.0 کامل با تمام endpoint ها
- **`SWAGGER_SETUP.md`** - راهنمای استفاده از Swagger
- **`README.md`** - این فایل

## 🚀 نحوه استفاده

### 1. Swagger UI Online
```
https://editor.swagger.io/
```
فایل `swagger.yaml` را کپی و paste کنید

### 2. Swagger UI محلی
```bash
# با Docker
docker run -p 8080:8080 -e SWAGGER_JSON=/swagger.yaml -v $(pwd):/swagger swaggerapi/swagger-ui

# یا با npm
npm install -g swagger-ui-serve
swagger-ui-serve swagger.yaml
```

### 3. Redoc
```bash
# نصب Redoc
npm install -g redoc-cli

# تولید documentation
redoc-cli serve swagger.yaml
```

## 📊 API Coverage

- **Auth Service**: 8 endpoints (احراز هویت)
- **Blog Service**: 12 endpoints (مدیریت محتوا)
- **Media Service**: 5 endpoints (مدیریت فایل)
- **Total**: 25 endpoint کامل

## 🔧 ویژگی‌های Swagger

1. **Complete API Coverage**: تمام endpoint ها
2. **Detailed Schemas**: مدل‌های کامل request/response
3. **Authentication**: JWT Bearer token support
4. **Examples**: مثال‌های واقعی
5. **Error Handling**: تمام کدهای خطا
6. **Interactive Testing**: تست مستقیم از UI

## 📚 منابع بیشتر

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Redoc](https://redoc.ly/)
- [API Design Best Practices](https://swagger.io/resources/articles/best-practices-in-api-design/)
