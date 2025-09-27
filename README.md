# 📘 School Management System (Go + MySQL)

A simple **School Management System** built with **Go (net/http)** and **MySQL**.  
It provides APIs for managing **students, teachers, and execs (admins/staff)** along with authentication and password management.  

---

## 🚀 Features
- REST API server with Go (`net/http`).
- MySQL-backed data storage.
- CRUD operations for Students, Teachers, and Execs.
- Authentication: login, logout, password reset.
- Graceful shutdown handling (`SIGINT`, `SIGTERM`).
- Modular, clean architecture (`internal` & `pkg` separation).

---

## 🛠️ Tech Stack
- **Go 1.24+**
- **MySQL 8+**

---

## ⚙️ Setup & Run

### 1. Clone the repo
```bash
git clone https://github.com/your-username/school-management-system.git
cd school-management-system

2. Setup MySQL Database

Install MySQL and start it.

Create a database:  
CREATE DATABASE school_mgmt;


Import schema:

mysql -u root -p school_mgmt < schema.sql

3. Configure Environment Variables

Create a .env file in the root:

DB_USER=root
DB_PASS=root
DB_NAME=school_mgmt
DB_HOST=localhost
DB_PORT=3306

4. Run the App
go run command/main.go

📂 Project Structure
REST_API_GO/
│── command/                  # Entrypoint (main.go)
│── internal/
│   ├── api/
│   │   ├── handlers/         # Request handlers
│   │   ├── middlewares/      # Middlewares
│   │   └── router/           # API routes
│   ├── models/               # Database models
│   └── repositories/
│       └── sqlconn/          # MySQL connection & queries
│── pkg/
│   └── utils/                # Helper utilities
│── schema.sql                # MySQL schema
│── go.mod / go.sum           # Dependencies
│── .env                      # Config (ignored in git)
│── .gitignore
│── openssl.cnf               # SSL config (if used)

✅ API Endpoints

👨‍🎓 Students

GET /students → Get all students

GET /students/{id} → Get one student

POST /students → Add new student

PATCH /students → Update student

DELETE /students/{id} → Delete student

GET /teacher/{id}/students → Get students by teacher ID


👨‍🏫 Teachers

GET /teachers → Get all teachers

GET /teachers/{id} → Get one teacher

POST /teachers → Add new teacher

PATCH /teachers → Update teacher

DELETE /teachers/{id} → Delete teacher


🏫 Execs (Admins / Staff)

GET /execs → Get all execs

GET /execs/{id} → Get one exec

POST /execs → Add new exec

PATCH /execs → Update exec

DELETE /execs/{id} → Delete exec

POST /execs/{id}/updatepassword → Update exec password

POST /execs/login → Exec login

POST /execs/logout → Exec logout

POST /execs/forgotpassword → Forgot password (send reset code)

GET /execs/resetpassword/reset/{resetcode} → Reset password with reset code