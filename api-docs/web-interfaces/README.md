# 🌐 Web Interfaces

این پوشه شامل فایل‌های مربوط به رابط‌های وب برای documentation است.

## 📁 فایل‌های موجود

- **`swagger-ui.html`** - Swagger UI محلی
- **`redoc.html`** - Redoc documentation
- **`package.json`** - npm scripts برای documentation
- **`README.md`** - این فایل

## 🚀 نحوه استفاده

### 1. Swagger UI محلی
```bash
# باز کردن فایل HTML
open swagger-ui.html

# یا با HTTP server
npx http-server . -p 8080 -o
```

### 2. Redoc
```bash
# باز کردن فایل HTML
open redoc.html

# یا با HTTP server
npx http-server . -p 8081 -o
```

### 3. npm Scripts
```bash
# نصب dependencies
npm install

# اجرای Swagger UI
npm run docs:swagger

# اجرای Redoc
npm run docs:redoc

# ساخت documentation
npm run docs:build

# اعتبارسنجی Swagger
npm run docs:validate

# تولید documentation
npm run docs:generate

# اجرای HTTP server
npm run serve
```

## 🔧 ویژگی‌های Web Interfaces

### Swagger UI
- **Interactive Testing**: تست مستقیم API ها
- **Request/Response Examples**: مثال‌های کامل
- **Authentication Support**: پشتیبانی از JWT
- **File Upload**: آپلود فایل برای تست
- **Responsive Design**: طراحی واکنش‌گرا

### Redoc
- **Beautiful Documentation**: documentation زیبا
- **Search Functionality**: قابلیت جستجو
- **Code Examples**: مثال‌های کد
- **Schema Visualization**: نمایش schema ها
- **Mobile Friendly**: سازگار با موبایل

## 📊 npm Scripts

| Script | Description |
|--------|-------------|
| `docs:swagger` | اجرای Swagger UI |
| `docs:redoc` | اجرای Redoc |
| `docs:build` | ساخت documentation |
| `docs:validate` | اعتبارسنجی Swagger |
| `docs:generate` | تولید documentation |
| `serve` | اجرای HTTP server |
| `start` | اجرای HTTP server |

## 🔧 Customization

### Swagger UI Configuration
```javascript
const ui = SwaggerUIBundle({
  url: './swagger.yaml',
  dom_id: '#swagger-ui',
  deepLinking: true,
  presets: [
    SwaggerUIBundle.presets.apis,
    SwaggerUIStandalonePreset
  ],
  plugins: [
    SwaggerUIBundle.plugins.DownloadUrl
  ],
  layout: "StandaloneLayout"
});
```

### Redoc Configuration
```html
<redoc spec-url='./swagger.yaml'></redoc>
```

## 📚 منابع بیشتر

- [Swagger UI](https://swagger.io/tools/swagger-ui/)
- [Redoc](https://redoc.ly/)
- [HTTP Server](https://www.npmjs.com/package/http-server)
- [npm Scripts](https://docs.npmjs.com/cli/v7/commands/npm-run-script)
