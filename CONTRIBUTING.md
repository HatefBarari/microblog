# 🤝 Contributing to Microblog Platform

از مشارکت شما در پروژه Microblog Platform خوشحالیم! این راهنما به شما کمک می‌کند تا به راحتی در پروژه مشارکت کنید.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## 📜 Code of Conduct

این پروژه و تمام مشارکت‌کنندگان آن متعهد به احترام به تمام افراد و ایجاد محیطی بدون آزار و اذیت برای همه هستند.

## 🚀 Getting Started

### Prerequisites

- Go 1.21 یا بالاتر
- MongoDB 7.0 یا بالاتر
- Docker و Docker Compose
- Git

### Development Setup

1. **Fork کنید**
   ```bash
   # روی GitHub روی دکمه Fork کلیک کنید
   ```

2. **Clone کنید**
   ```bash
   git clone https://github.com/YOUR_USERNAME/microblog.git
   cd microblog
   ```

3. **Remote اضافه کنید**
   ```bash
   git remote add upstream https://github.com/HatefBarari/microblog.git
   ```

4. **Dependencies نصب کنید**
   ```bash
   # MongoDB
   docker run -d --name mongo -p 27017:27017 mongo:7
   
   # MailHog (برای تست ایمیل)
   docker run -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog
   ```

## 🔧 Development Setup

### 1. Environment Variables

```bash
# کپی کردن فایل environment
cp deployments/.env.example deployments/.env

# ویرایش متغیرهای محیطی
nano deployments/.env
```

### 2. Running Services

```bash
# اجرای تمام سرویس‌ها
cd deployments
docker-compose up -d

# یا اجرای جداگانه
cd auth-service && go run cmd/server/main.go &
cd blog-service && go run cmd/server/main.go &
cd media-service && go run cmd/server/main.go &
```

### 3. Testing

```bash
# اجرای تمام تست‌ها
make test

# یا تست جداگانه
cd auth-service && go test ./...
cd blog-service && go test ./...
cd media-service && go test ./...
```

## 📝 Making Changes

### 1. Branch ایجاد کنید

```bash
# از main branch شروع کنید
git checkout main
git pull upstream main

# branch جدید ایجاد کنید
git checkout -b feature/your-feature-name
# یا
git checkout -b fix/your-bug-fix
# یا
git checkout -b docs/your-documentation
```

### 2. تغییرات را اعمال کنید

```bash
# فایل‌ها را ویرایش کنید
# ...

# تغییرات را stage کنید
git add .

# commit کنید
git commit -m "feat: add new feature"
# یا
git commit -m "fix: resolve bug in authentication"
# یا
git commit -m "docs: update API documentation"
```

### 3. Push کنید

```bash
git push origin feature/your-feature-name
```

## 🔄 Pull Request Process

### 1. Pull Request ایجاد کنید

1. به GitHub repository بروید
2. روی "New Pull Request" کلیک کنید
3. branch خود را انتخاب کنید
4. توضیحات کاملی بنویسید

### 2. Pull Request Template

```markdown
## 📝 Description
توضیح کوتاه از تغییرات

## 🔧 Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## 🧪 Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## 📋 Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
```

### 3. Review Process

- **Code Review**: حداقل یک نفر باید کد را review کند
- **Testing**: تمام تست‌ها باید pass کنند
- **Documentation**: documentation باید به‌روز باشد

## 📏 Coding Standards

### Go Code Style

```go
// ✅ Good
func (uc *UserUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    if err := uc.validateUser(req); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    user := &User{
        ID:    generateID(),
        Name:  req.Name,
        Email: req.Email,
    }
    
    if err := uc.repo.Save(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to save user: %w", err)
    }
    
    return user, nil
}

// ❌ Bad
func (uc *UserUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    if uc.validateUser(req) != nil {
        return nil, errors.New("validation failed")
    }
    user := &User{ID: generateID(), Name: req.Name, Email: req.Email}
    if uc.repo.Save(ctx, user) != nil {
        return nil, errors.New("failed to save user")
    }
    return user, nil
}
```

### Naming Conventions

- **Packages**: lowercase, single word
- **Functions**: PascalCase for public, camelCase for private
- **Variables**: camelCase
- **Constants**: UPPER_CASE
- **Interfaces**: -er suffix (e.g., UserRepository, EmailSender)

### Error Handling

```go
// ✅ Good
if err != nil {
    return fmt.Errorf("failed to process request: %w", err)
}

// ❌ Bad
if err != nil {
    return err
}
```

## 🧪 Testing

### Unit Tests

```go
func TestUserUseCase_CreateUser(t *testing.T) {
    // Arrange
    mockRepo := &MockUserRepository{}
    uc := NewUserUseCase(mockRepo)
    
    req := CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Act
    user, err := uc.CreateUser(context.Background(), req)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "John Doe", user.Name)
    assert.Equal(t, "john@example.com", user.Email)
}
```

### Integration Tests

```go
func TestUserIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Test user creation flow
    // ...
}
```

### Test Coverage

```bash
# اجرای تست‌ها با coverage
go test -cover ./...

# گزارش coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📚 Documentation

### Code Documentation

```go
// UserUseCase handles user-related business logic
type UserUseCase struct {
    repo  UserRepository
    email EmailSender
    log   *zap.Logger
}

// CreateUser creates a new user with the provided information
// It validates the input, generates a unique ID, and saves the user
func (uc *UserUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
    // Implementation
}
```

### API Documentation

- تمام endpoint ها باید در `api-docs/swagger/swagger.yaml` مستند شوند
- مثال‌های request/response باید کامل باشند
- Error codes و messages باید مشخص باشند

### README Files

- هر سرویس باید README.md داشته باشد
- راهنمای نصب و اجرا باید کامل باشد
- مثال‌های استفاده باید موجود باشند

## 🐛 Bug Reports

### Bug Report Template

```markdown
## 🐛 Bug Description
توضیح کوتاه از مشکل

## 🔄 Steps to Reproduce
1. Go to '...'
2. Click on '...'
3. See error

## 🎯 Expected Behavior
توضیح رفتار مورد انتظار

## 📱 Environment
- OS: [e.g. Windows 10]
- Go Version: [e.g. 1.21]
- Browser: [e.g. Chrome 91]

## 📷 Screenshots
اگر مربوط است، screenshot اضافه کنید

## 📋 Additional Context
اطلاعات اضافی
```

## ✨ Feature Requests

### Feature Request Template

```markdown
## 🚀 Feature Description
توضیح کوتاه از feature درخواستی

## 💡 Motivation
چرا این feature مفید است؟

## 📋 Detailed Description
توضیح کامل از feature

## 🎯 Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3

## 📷 Mockups
اگر مربوط است، mockup اضافه کنید
```

## 📞 Getting Help

- **GitHub Issues**: برای bug reports و feature requests
- **Discussions**: برای سوالات و بحث‌ها
- **Email**: برای مسائل خصوصی

## 📄 License

با مشارکت در این پروژه، شما موافقت می‌کنید که کد شما تحت مجوز MIT منتشر شود.

---

**متشکریم از مشارکت شما! 🙏**
