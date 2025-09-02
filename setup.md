# External Libs used  
```
go get github.com/go-sql-driver/mysql
go get github.com/joho/godotenv
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