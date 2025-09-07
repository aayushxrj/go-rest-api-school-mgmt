# Helpful commands
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config openssl.cnf
```

# External Libs used  
```
go get github.com/go-sql-driver/mysql
go get github.com/joho/godotenv
go get golang.org/x/crypto/argon2
```

# Setting up MariaDB on WSL Local (PORT 3306)

```
sudo apt install mariadb-server mariadb-client -y
```

```
sudo service mariadb start
sudo service mariadb status
```

```
sudo mysql_secure_installation
```

```
sudo mariadb -u root -p
```

```
CREATE DATABASE IF NOT EXISTS school;
USE school;

CREATE TABLE IF NOT EXISTS teachers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    class VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    INDEX (email)
) AUTO_INCREMENT=100;
```
```
CREATE INDEX idx_class ON teachers(class);

CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    class VARCHAR(255) NOT NULL,
    INDEX (email),
    FOREIGN KEY (class) REFERENCES teachers(class)
) AUTO_INCREMENT=100;
```
```
CREATE TABLE IF NOT EXISTS execs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    password_changed_at VARCHAR(255),
    user_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    password_reset_token VARCHAR(255),
    inactive_status BOOLEAN NOT NULL,
    role VARCHAR(50) NOT NULL,
    INDEX idx_email (email),
    INDEX idx_username (username)
);
```
```
ALTER TABLE execs ADD COLUMN password_token_expires VARCHAR(255);
```

# Swagger

```
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/http-swagger
```
```
swag init -g cmd/api/server.go
```

Update
```
w.Header().Set("Content-Security-Policy", "default-src 'self'")
```
to
```
w.Header().Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline'")
```

```
// @title School Management API
// @version 1.0
// @description REST API for managing school data.
// @schemes https
// @host localhost:3000
// @BasePath /
package main
```

```
https://localhost:3000/swagger/index.html
```

# TODO

- Reduce boilerplate code using ([GORM](https://gorm.io/docs/index.html)) (Manual to Automatic Query Generation)

- Implement data validation using [Validator](https://github.com/go-playground/validator)