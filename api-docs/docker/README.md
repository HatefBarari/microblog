# 🐳 Docker Documentation Setup

این پوشه شامل تنظیمات Docker برای اجرای documentation services است.

## 📁 فایل‌های موجود

- **`docker-compose.docs.yml`** - Docker setup برای documentation
- **`README.md`** - این فایل

## 🚀 نحوه استفاده

### 1. اجرای تمام Services
```bash
# اجرای تمام documentation services
docker-compose -f docker-compose.docs.yml up -d

# دسترسی به services:
# - Swagger UI: http://localhost:8080
# - Redoc: http://localhost:8081  
# - Swagger Editor: http://localhost:8082
# - Custom HTML: http://localhost:8083
```

### 2. اجرای Service های جداگانه
```bash
# فقط Swagger UI
docker run -p 8080:8080 -e SWAGGER_JSON=/swagger.yaml -v $(pwd)/../swagger:/swagger swaggerapi/swagger-ui

# فقط Redoc
docker run -p 8081:80 -v $(pwd)/../swagger:/usr/share/nginx/html/swagger.yaml:ro redocly/redoc

# فقط Swagger Editor
docker run -p 8082:8080 -v $(pwd)/../swagger:/swagger.yaml:ro swaggerapi/swagger-editor
```

## 🔧 Services Overview

### Swagger UI (Port: 8080)
- **Image**: `swaggerapi/swagger-ui:latest`
- **Description**: Interactive API documentation
- **Features**: Testing, Authentication, File Upload

### Redoc (Port: 8081)
- **Image**: `redocly/redoc:latest`
- **Description**: Beautiful documentation
- **Features**: Search, Code Examples, Schema Visualization

### Swagger Editor (Port: 8082)
- **Image**: `swaggerapi/swagger-editor:latest`
- **Description**: Online Swagger editor
- **Features**: Editing, Validation, Code Generation

### Nginx (Port: 8083)
- **Image**: `nginx:alpine`
- **Description**: Custom HTML interfaces
- **Features**: Static file serving, Custom routing

## 📊 Docker Compose Configuration

```yaml
version: '3.8'

services:
  swagger-ui:
    image: swaggerapi/swagger-ui:latest
    container_name: microblog-swagger-ui
    ports:
      - "8080:8080"
    environment:
      - SWAGGER_JSON=/swagger.yaml
    volumes:
      - ./swagger.yaml:/swagger.yaml:ro
    networks:
      - microblog-docs

  redoc:
    image: redocly/redoc:latest
    container_name: microblog-redoc
    ports:
      - "8081:80"
    volumes:
      - ./swagger.yaml:/usr/share/nginx/html/swagger.yaml:ro
    networks:
      - microblog-docs

  swagger-editor:
    image: swaggerapi/swagger-editor:latest
    container_name: microblog-swagger-editor
    ports:
      - "8082:8080"
    volumes:
      - ./swagger.yaml:/swagger.yaml:ro
    environment:
      - SWAGGER_FILE=/swagger.yaml
    networks:
      - microblog-docs

  nginx:
    image: nginx:alpine
    container_name: microblog-docs-nginx
    ports:
      - "8083:80"
    volumes:
      - ./swagger.yaml:/usr/share/nginx/html/swagger.yaml:ro
      - ./swagger-ui.html:/usr/share/nginx/html/index.html:ro
      - ./redoc.html:/usr/share/nginx/html/redoc.html:ro
    networks:
      - microblog-docs

networks:
  microblog-docs:
    driver: bridge
```

## 🔧 Management Commands

### اجرا و توقف
```bash
# اجرای services
docker-compose -f docker-compose.docs.yml up -d

# توقف services
docker-compose -f docker-compose.docs.yml down

# مشاهده logs
docker-compose -f docker-compose.docs.yml logs -f

# restart services
docker-compose -f docker-compose.docs.yml restart
```

### بررسی وضعیت
```bash
# وضعیت containers
docker-compose -f docker-compose.docs.yml ps

# استفاده از resources
docker-compose -f docker-compose.docs.yml top

# بررسی networks
docker network ls | grep microblog
```

### Cleanup
```bash
# حذف containers و volumes
docker-compose -f docker-compose.docs.yml down -v

# حذف images
docker-compose -f docker-compose.docs.yml down --rmi all

# حذف networks
docker-compose -f docker-compose.docs.yml down --remove-orphans
```

## 🔧 Customization

### Environment Variables
```yaml
services:
  swagger-ui:
    environment:
      - SWAGGER_JSON=/swagger.yaml
      - SWAGGER_UI_DISPLAY_REQUEST_DURATION=true
      - SWAGGER_UI_TRY_IT_OUT_ENABLED=true
```

### Volume Mounts
```yaml
services:
  swagger-ui:
    volumes:
      - ./swagger.yaml:/swagger.yaml:ro
      - ./custom-theme.css:/usr/share/nginx/html/custom-theme.css:ro
```

### Network Configuration
```yaml
networks:
  microblog-docs:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

## 📚 منابع بیشتر

- [Docker Compose](https://docs.docker.com/compose/)
- [Swagger UI Docker](https://hub.docker.com/r/swaggerapi/swagger-ui)
- [Redoc Docker](https://hub.docker.com/r/redocly/redoc)
- [Nginx Docker](https://hub.docker.com/_/nginx)
