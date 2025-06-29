# Ecom Go Project (Dockerized)

This is a simple e-commerce backend built in Go (Gin + MySQL), fully Dockerized for easy deployment and local development. It includes JWT auth, Docker, and GitHub Actions CI/CD.

## 📂 Project Structure

```
ecom/
├── cmd/
│ └── main.go # Entry point
│
├── api/ # API server setup
│
├── config/ # Env configuration loader
│
├── db/ # DB connection setup
│
├── service/
│ ├── auth/ # JWT & password helpers
│ ├── products/ # Product service & routes
│ └── user/ # User service & routes
│
├── types/ # Types, interfaces
├── utils/ # Utility functions
│
├── Dockerfile # Go app Dockerfile
├── docker-compose.yml # For production / EC2
├── docker-compose.override.yml # For local dev (optional)
└── .github/workflows/deploy.yml # GitHub Actions CI/CD
```

## 🚀 How to Run Locally (Dockerized)

### 1️⃣ Clone the Repo

```bash
git clone <repo-url>
cd ecom
```

### 2️⃣ Start App + DB with Docker Compose

```bash
docker-compose up --build
```

### 👉 This:

- Builds the Go app container (via Dockerfile)
- Spins up MySQL (exposed on 3307)
- App accessible on: [http://localhost:8080](http://localhost:8080)

### ⚠ DB Credentials (via Environment Variables)

```plaintext
DB_USER=root
DB_PASSWORD=kulkarni11
DB_NAME=ecom
DB_HOST=db
DB_PORT=3306
```

MySQL will listen on `localhost:3307` for external clients.

## 🛠 Local Development (Without Docker)

If you prefer running without Docker:

1️⃣ Ensure MySQL is running locally (on `127.0.0.1:3306` or adjust config).  
2️⃣ Load `.env` variables or set them manually.  
3️⃣ Run:

```bash
go run ./cmd/server/main.go
```

## 🐳 Docker Compose Summary

- **`docker-compose.yml`**: Used for production (EC2 / server).  
  Uses pre-built ECR image.

- **`docker-compose.override.yml`**: Used for local development.  
  Builds app image locally from Dockerfile.

## 🚀 CI/CD: GitHub Actions + EC2

### ✅ On Push to `main`:

- Builds Go app Docker image.
- Pushes to Amazon ECR.
- SSH into EC2 → Pulls latest image → Restarts with Docker Compose.

### 📂 Workflow: `.github/workflows/deploy.yml`

#### Secrets Needed:

- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `AWS_REGION`
- `AWS_ACCOUNT_ID`
- `EC2_HOST`
- `EC2_SSH_KEY` (private SSH key for EC2)

## 🔑 Useful Commands

### 💡 Build and Run Locally:

```bash
docker-compose up --build
```

### 💡 Stop Containers:

```bash
docker-compose down
```

### 💡 Check Logs:

```bash
docker-compose logs -f
```

### 💡 Run Go Tests:

```bash
go test ./...
```

## ⚡ Deployment on EC2

✅ Copy `docker-compose.yml` to your EC2 machine.  
✅ Run:

```bash
docker-compose pull
docker-compose up -d
```

## 🔥 Notes

- MySQL data is persisted in `db_data` volume.
- You can connect to MySQL from a tool like MySQL Workbench on `localhost:3307` (root / kulkarni11).

## 📝 TODOs for Production Hardening

- Use secrets management for `DB_PASSWORD` + `JWT_SECRET`.
- Enable HTTPS with reverse proxy (Nginx / ALB).
- Add database migrations (e.g., golang-migrate).
- Use Different DB like RDS (right now we are running mysql in EC2 insatnce only)
- RBAC
- Session Management

## Bug Fix

- Generating token in local and put in PROD also working ( Need to have Secrets Seperatly for Local env and PROD)
