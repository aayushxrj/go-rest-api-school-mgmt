# 🏫 School Management REST API

[![Go Version](https://img.shields.io/badge/Go-1.25.0-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A comprehensive REST API for managing school data built with Go. This API enables administrative staff to efficiently manage students, teachers, and executive staff members with robust authentication, security features, and complete CRUD operations.

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Project Structure](#-project-structure)
- [Getting Started](#-getting-started)
- [API Documentation](#-api-documentation)
- [Security Features](#-security-features)
- [Database Schema](#-database-schema)
- [Environment Variables](#-environment-variables)
- [Contributing](#-contributing)

## ✨ Features

### Core Functionality
- **Complete CRUD Operations** for Students, Teachers, and Executives
- **Bulk Operations** for efficient data management
- **Class Management** with teacher-student relationships
- **Advanced Filtering & Sorting** on all list endpoints
- **JWT-based Authentication** with secure token management
- **Password Management** (reset, forgot password, update password)
- **User Deactivation** capabilities

### Security & Performance
- **HTTPS/TLS** with HTTP/2 support
- **JWT Authentication** middleware
- **Rate Limiting** to prevent abuse
- **XSS Protection** with input sanitization
- **Security Headers** (HSTS, CSP, X-Frame-Options, etc.)
- **CORS** configuration
- **Response Compression** (gzip)
- **Response Time Tracking**

### Developer Experience
- **Swagger/OpenAPI Documentation** for easy API exploration
- **Modular Architecture** for maintainability
- **Comprehensive Error Handling**
- **Clean Code Structure** following Go best practices

## 🛠 Tech Stack

- **Language**: Go 1.25.0
- **Database**: MySQL
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Documentation**: Swagger/OpenAPI (swaggo)
- **Security**: 
  - bcrypt for password hashing
  - bluemonday for XSS protection
  - Custom middleware suite
- **Email**: go-mail for password reset emails
- **TLS/HTTPS**: Built-in Go crypto/tls
- **HTTP/2**: golang.org/x/net/http2

## 📁 Project Structure

```
go-rest-api-school-mgmt/
├── cmd/
│   └── api/
│       ├── cert.pem              # TLS certificate
│       ├── key.pem               # TLS private key
│       └── server.go             # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/             # HTTP request handlers
│   │   │   ├── execs.go
│   │   │   ├── students.go
│   │   │   ├── teachers.go
│   │   │   ├── helpers.go
│   │   │   └── root.go
│   │   ├── middlewares/          # HTTP middlewares
│   │   │   ├── jwt_middleware.go
│   │   │   ├── rate_limiter.go
│   │   │   ├── security_headers.go
│   │   │   ├── compression.go
│   │   │   ├── cors.go
│   │   │   ├── sanitize.go
│   │   │   └── ...
│   │   └── router/               # Route definitions
│   │       ├── router.go
│   │       ├── execs_router.go
│   │       ├── students_router.go
│   │       └── teachers_router.go
│   ├── models/                   # Data models
│   │   ├── exec.go
│   │   ├── student.go
│   │   └── teacher.go
│   └── repository/
│       └── sqlconnect/           # Database layer
│           ├── sqlconfig.go
│           ├── execs_crud.go
│           ├── students_crud.go
│           └── teachers_crud.go
├── pkg/
│   └── utils/                    # Utility functions
│       ├── jwt.go
│       ├── password.go
│       ├── error_handler.go
│       ├── authorize_user.go
│       └── database_utils.go
├── docs/                         # Swagger documentation
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── data/                         # JSON data files
│   ├── execs_data.json
│   ├── students_data.json
│   └── teachers_data.json
├── proto/                        # Protocol buffers (if used)
│   └── main.proto
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Getting Started

### Prerequisites

- Go 1.25.0 or higher
- MySQL 8.0+
- OpenSSL (for generating TLS certificates)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/aayushxrj/go-rest-api-school-mgmt.git
   cd go-rest-api-school-mgmt
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   API_PORT=3000
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=school_management
   JWT_SECRET=your_jwt_secret_key
   EMAIL_HOST=smtp.gmail.com
   EMAIL_PORT=587
   EMAIL_USER=your_email@gmail.com
   EMAIL_PASSWORD=your_email_password
   ```

4. **Set up the database**
   ```sql
   CREATE DATABASE school_management;
   USE school_management;
   
   -- Create tables (run the SQL schema from your database setup)
   ```

5. **Generate TLS certificates** (if not already present)
   ```bash
   openssl req -x509 -newkey rsa:4096 -keyout cmd/api/key.pem -out cmd/api/cert.pem -days 365 -nodes -config openssl.cnf
   ```

6. **Run the application**
   ```bash
   go run cmd/api/server.go
   ```

7. **Access the API**
   - API Base URL: `https://localhost:3000`
   - Swagger Documentation: `https://localhost:3000/swagger/index.html`

## 📚 API Documentation

### Authentication

All endpoints (except login and password reset) require JWT authentication via the `Authorization` header:
```
Authorization: Bearer <your_jwt_token>
```

### Executives Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/execs` | Get list of executives with filtering & sorting |
| POST | `/execs` | Create a new executive |
| PATCH | `/execs` | Bulk update executives |
| GET | `/execs/{id}` | Get a specific executive |
| PATCH | `/execs/{id}` | Update a specific executive |
| DELETE | `/execs/{id}` | Delete a specific executive |
| POST | `/execs/login` | Login (returns JWT token) |
| POST | `/execs/logout` | Logout |
| POST | `/execs/forgotpassword` | Request password reset |
| POST | `/execs/resetpassword/reset/{resetcode}` | Reset password with token |
| PATCH | `/execs/{id}/updatepassword` | Update password |

### Students Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/students` | Get list of students with filtering & sorting |
| POST | `/students` | Create a new student |
| PATCH | `/students` | Bulk update students |
| DELETE | `/students` | Bulk delete students |
| GET | `/students/{id}` | Get a specific student |
| PUT | `/students/{id}` | Replace a specific student |
| PATCH | `/students/{id}` | Update a specific student |
| DELETE | `/students/{id}` | Delete a specific student |

### Teachers Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/teachers` | Get list of teachers with filtering & sorting |
| POST | `/teachers` | Create a new teacher |
| PATCH | `/teachers` | Bulk update teachers |
| DELETE | `/teachers` | Bulk delete teachers |
| GET | `/teachers/{id}` | Get a specific teacher |
| PUT | `/teachers/{id}` | Replace a specific teacher |
| PATCH | `/teachers/{id}` | Update a specific teacher |
| DELETE | `/teachers/{id}` | Delete a specific teacher |
| GET | `/teachers/{id}/students` | Get students taught by a teacher |
| GET | `/teachers/{id}/studentcount` | Get student count for a teacher |

### Query Parameters

Most GET endpoints support:
- **Filtering**: `?first_name=John&class=10A`
- **Sorting**: `?sortBy=last_name&sortOrder=asc`
- **Pagination**: `?limit=10&offset=0`

### Example Requests

**Login:**
```bash
curl -X POST https://localhost:3000/execs/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "secure_password"
  }'
```

**Get Students (with filtering):**
```bash
curl -X GET "https://localhost:3000/students?class=10A&sortBy=last_name&sortOrder=asc" \
  -H "Authorization: Bearer <your_jwt_token>"
```

**Create a Teacher:**
```bash
curl -X POST https://localhost:3000/teachers \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane",
    "last_name": "Doe",
    "email": "jane.doe@school.com",
    "class": "10A",
    "subject": "Mathematics"
  }'
```

## 🔒 Security Features

### Authentication & Authorization
- **JWT Tokens**: Secure, stateless authentication
- **Password Hashing**: bcrypt with appropriate cost factor
- **Token Expiration**: Automatic token invalidation
- **Role-based Access**: Different permissions for different user types

### Security Middleware Stack
1. **CORS**: Configurable cross-origin resource sharing
2. **Rate Limiting**: Prevents API abuse (5 requests per minute)
3. **Response Time Tracking**: Performance monitoring
4. **Security Headers**:
   - Strict-Transport-Security (HSTS)
   - Content-Security-Policy (CSP)
   - X-Frame-Options
   - X-Content-Type-Options
   - X-XSS-Protection
5. **Compression**: Gzip compression for responses
6. **HPP Protection**: HTTP Parameter Pollution prevention
7. **XSS Middleware**: Input sanitization using bluemonday

### HTTPS/TLS
- Minimum TLS version: 1.2
- HTTP/2 support enabled
- Strong cipher suites

## 🗄 Database Schema

### Students Table
```sql
CREATE TABLE students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    class VARCHAR(10) NOT NULL
);
```

### Teachers Table
```sql
CREATE TABLE teachers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    class VARCHAR(10) NOT NULL,
    subject VARCHAR(50) NOT NULL
);
```

### Executives Table
```sql
CREATE TABLE execs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    password_changed_at TIMESTAMP NULL,
    user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    password_reset_token VARCHAR(255) NULL,
    password_token_expires TIMESTAMP NULL,
    inactive_status BOOLEAN DEFAULT FALSE,
    role VARCHAR(20) NOT NULL
);
```

## 🌍 Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `API_PORT` | Port for the API server | `3000` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `3306` |
| `DB_USER` | Database username | `root` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `school_management` |
| `JWT_SECRET` | Secret key for JWT signing | `your_secret_key` |
| `EMAIL_HOST` | SMTP server host | `smtp.gmail.com` |
| `EMAIL_PORT` | SMTP server port | `587` |
| `EMAIL_USER` | Email address for sending | `noreply@school.com` |
| `EMAIL_PASSWORD` | Email password/app password | `app_password` |

## 🧪 Testing

Run tests with:
```bash
go test ./... -v
```

Run tests with coverage:
```bash
go test ./... -cover
```

## 📝 Best Practices Implemented

- ✅ **Modularity**: Clear separation of concerns
- ✅ **Documentation**: Swagger/OpenAPI integration
- ✅ **Error Handling**: Comprehensive error responses
- ✅ **Security**: Multiple layers of security
- ✅ **Logging**: Structured logging for debugging
- ✅ **Configuration**: Environment-based configuration
- ✅ **Code Quality**: Following Go best practices
- ✅ **API Design**: RESTful conventions

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

Please ensure your code:
- Follows Go formatting guidelines (`go fmt`)
- Includes appropriate tests
- Updates documentation as needed
- Follows the existing code structure

## 📄 License

This project is licensed under the MIT License