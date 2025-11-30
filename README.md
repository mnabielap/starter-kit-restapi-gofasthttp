# 🚀 Go FastHTTP REST API Starter Kit

[![Go Version](https://img.shields.io/badge/Go-1.22%2B-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![FastHTTP](https://img.shields.io/badge/FastHTTP-High%20Performance-blueviolet?style=for-the-badge)](https://github.com/valyala/fasthttp)
[![GORM](https://img.shields.io/badge/GORM-SQLite%20%7C%20Postgres-lightgrey?style=for-the-badge)](https://gorm.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

A high-performance, production-ready REST API starter kit built with **Go**, **FastHTTP**, and **Clean Architecture**. This project is designed to be database-agnostic (switching easily between **SQLite** and **PostgreSQL**) and fully containerized.

## ✨ Key Features

*   **⚡ Blazing Fast**: Built on [FastHTTP](https://github.com/valyala/fasthttp), up to 10x faster than `net/http`.
*   **🏗️ Clean Architecture**: Clear separation of concerns (Handler -> Service -> Repository -> DB).
*   **🗄️ Database Agnostic**: Uses [GORM](https://gorm.io/) to support **SQLite** (default for dev) and **PostgreSQL** (production).
*   **🔐 Secure Auth**: Complete JWT authentication system (Access & Refresh tokens) with Role-Based Access Control (RBAC).
*   **🐳 Docker Ready**: Multi-stage build Dockerfile for a tiny, secure production image.
*   **📝 Auto Documentation**: Integrated **Swagger/OpenAPI** generated from code comments.
*   **🧪 Scripted Testing**: Includes a suite of Python scripts for API testing (Postman alternative).
*   **🪵 Logging & Config**: Structured logging and environment variable management.

---

## 📂 Project Structure

```text
.
├── cmd/api/main.go            # 🏁 Application Entry Point
├── internal/
│   ├── config/                # ⚙️ Configuration Loader
│   ├── database/              # 🔌 Database Connection (GORM)
│   ├── http/                  # 🌐 HTTP Layer
│   │   ├── handler/           # 🎮 Request Handlers (Controllers)
│   │   ├── middleware/        # 🛡️ Middleware (Auth, Logger)
│   │   └── router/            # 🗺️ Route Definitions
│   ├── model/                 # 📦 Data Models & DB Schemas
│   ├── repository/            # 🗄️ Data Access Layer (Repository Pattern)
│   └── service/               # 🧠 Business Logic Layer
├── pkg/                       # 🛠️ Shared Utilities (Logger, Validator, Response)
├── api_tests/                 # 🧪 Python API Test Scripts
├── Dockerfile                 # 🐳 Docker Build Instructions
├── entrypoint.sh              # 🚀 Docker Startup Script
├── go.mod                     # 📦 Go Dependencies
└── .env                       # 🌍 Environment Variables
```

---

## 🛠️ Prerequisites

*   **Go** (version 1.22 or higher)
*   **Git**
*   **Python 3.x** (for running test scripts)
*   **Docker** (optional, for containerized deployment)

---

## 🚀 Getting Started (Local Development)

We recommended running the project locally first to understand how it works. By default, it uses **SQLite**, so you don't need to install any external database!

### 1. Clone the Repository
```bash
git clone https://github.com/mnabielap/starter-kit-restapi-gofasthttp.git
cd starter-kit-restapi-gofasthttp
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Setup Environment
Copy the example environment file.
```bash
cp .env.example .env
```
*By default, `.env` is configured for **SQLite** (`DB_DRIVER=sqlite`). No changes needed for a quick start.*

### 4. Generate Swagger Docs (Optional)
If you have `swag` installed:
```bash
swag init -g cmd/api/main.go
```

### 5. Run the Application
```bash
go run cmd/api/main.go
```
The server will start at **http://localhost:3000**.

> 📄 **Documentation:** Visit `http://localhost:3000/swagger/index.html` to see the interactive API docs.

---

## 🐳 Running with Docker (Production Mode)

This method runs the application and a **PostgreSQL** database in separate containers, connected via a custom network.

### 1. Create Network & Volumes
We need a network for the containers to talk to each other and volumes to persist data.
```bash
# Create network
docker network create restapi_gofasthttp_network

# Create volumes
docker volume create restapi_gofasthttp_db_volume
docker volume create restapi_gofasthttp_media_volume
```

### 2. Setup Docker Environment
Create a `.env.docker` file. This config tells the app to use Postgres.
```ini
PORT=5005
NODE_ENV=production
DB_DRIVER=postgres
# Hostname 'postgres-container' matches the name we give the DB container below
DB_SOURCE=host=postgres-container user=restapi password=secret dbname=restapidb port=5432 sslmode=disable
JWT_SECRET=complex_secret_key
JWT_ACCESS_EXPIRATION_MINUTES=60
JWT_REFRESH_EXPIRATION_DAYS=30
```

### 3. Build the Image
```bash
docker build -t restapi-gofasthttp-app .
```

### 4. Run PostgreSQL Container
```bash
docker run -d \
  --name postgres-container \
  --network restapi_gofasthttp_network \
  -e POSTGRES_USER=restapi \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=restapidb \
  -v restapi_gofasthttp_db_volume:/var/lib/postgresql/data \
  postgres:16-alpine
```

### 5. Run Application Container
```bash
docker run -d -p 5005:5005 \
  --env-file .env.docker \
  -v restapi_gofasthttp_media_volume:/app/media \
  --network restapi_gofasthttp_network \
  --name restapi-gofasthttp-container \
  restapi-gofasthttp-app
```
Your API is now running at **http://localhost:5005**.

---

## 🕹️ Docker Management Commands

Here are useful commands to manage your containers.

#### 📄 View Logs
See what's happening inside your app.
```bash
docker logs -f restapi-gofasthttp-container
```

#### 🛑 Stop Container
Temporarily stops the application.
```bash
docker stop restapi-gofasthttp-container
```

#### ▶️ Start Container
Restarts a stopped application.
```bash
docker start restapi-gofasthttp-container
```

#### 🗑️ Remove Container
Deletes the container instance (data in volumes remains safe).
```bash
docker stop restapi-gofasthttp-container
docker rm restapi-gofasthttp-container
```

#### 🗂️ View Volumes
```bash
docker volume ls
```

#### ⚠️ Remove Volumes (Danger Zone!)
**WARNING:** This will permanently delete your database and uploaded files.
```bash
docker volume rm restapi_gofasthttp_db_volume
docker volume rm restapi_gofasthttp_media_volume
```

---

## 🧪 API Testing (Python Scripts)

Instead of using Postman, we provide a suite of Python scripts in the `api_tests/` folder. These scripts automatically handle token management (saving/loading `access_token` to `secrets.json`).

### 1. Running Tests
Run the scripts in the following order to simulate a user flow. No arguments are needed.

1.  **Register Admin:**
    ```bash
    python api_tests/A1.auth_register.py
    ```
    *Creates a user and saves tokens to `secrets.json`.*

2.  **Create User:**
    ```bash
    python api_tests/B1.user_create.py
    ```
    *Uses the saved Admin Token to create a new user.*

3.  **Get All Users:**
    ```bash
    python api_tests/B2.user_get_all.py
    ```

4.  **Other Operations:**
    *   `auth_login.py`: Test login flow.
    *   `auth_refresh_tokens.py`: Test token rotation.
    *   `user_update.py`: Update the user created in step 2.
    *   `user_delete.py`: Delete the user.

> **Note:** The `utils.py` script generates detailed logs for every request in `response.json` (or specific output files defined in the scripts).

---

## ⚙️ Environment Variables

| Variable | Description | Example (SQLite) | Example (Postgres) |
| :--- | :--- | :--- | :--- |
| `PORT` | Server Port | `3000` | `5005` |
| `NODE_ENV` | Environment | `development` | `production` |
| `DB_DRIVER` | Database Type | `sqlite` | `postgres` |
| `DB_SOURCE` | DSN / File Path | `app.db` | `host=postgres...` |
| `JWT_SECRET` | Secret for signing tokens | `secret123` | `super_secure_key` |

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).