Ecom Go Project

This is a simple e-commerce backend built in Go (Gin + MySQL), with JWT auth, Docker, and GitHub Actions CI/CD.

📂 Project Structure

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

🚀 How to run locally (Docker)

1️⃣ Clone the repo

git clone <repo-url>
cd ecom

2️⃣ Start app + DB with Docker Compose

docker-compose up --build

👉 This:

Builds the Go app container (via Dockerfile)

Spins up MySQL (exposed on 3307)

App accessible on: http://localhost:8080

⚠ DB credentials

DB_USER=root
DB_PASSWORD=kulkarni11
DB_NAME=ecom
DB_HOST=db
DB_PORT=3306

MySQL will listen on localhost:3307 for external clients.

🛠 Local development (without Docker)

If you want to run without Docker:

1️⃣ Make sure MySQL is running locally (on 127.0.0.1:3306 or adjust config)2️⃣ Load .env variables or set them manually3️⃣ Run:

go run ./cmd/main.go

🐳 Docker Compose summary

docker-compose.yml

Used for production (EC2 / server)

Uses pre-built ECR image

docker-compose.override.yml

Used for local development

Builds app image locally from Dockerfile

🚀 CI/CD: GitHub Actions + EC2

✅ On push to main:

Builds Go app Docker image

Pushes to Amazon ECR

SSH into EC2 → Pulls latest image → Restarts with docker-compose

📂 Workflow: .github/workflows/deploy.yml

Secrets needed:

AWS_ACCESS_KEY_ID

AWS_SECRET_ACCESS_KEY

AWS_REGION

AWS_ACCOUNT_ID

EC2_HOST

EC2_SSH_KEY (private SSH key for EC2)

🔑 Useful commands

💡 Build and run locally:

docker-compose up --build

💡 Stop containers:

docker-compose down

💡 Check logs:

docker-compose logs -f

💡 Run Go tests:

go test ./...

⚡ Deployment on EC2

✅ Copy docker-compose.yml to your EC2 machine✅ Run:

docker-compose pull
docker-compose up -d

🔥 Notes

MySQL data is persisted in db_data volume.

You can connect to MySQL from a tool like MySQL Workbench on localhost:3307 (root / kulkarni11).

📝 TODOs for production hardening

Use secrets management for DB_PASSWORD + JWT_SECRET

Enable HTTPS with reverse proxy (Nginx / ALB)

Add database migrations (e.g. golang-migrate)
