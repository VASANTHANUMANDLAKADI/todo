# Todo API

A RESTful Todo Management API built with Go and PostgreSQL. The application can be run locally or deployed using Docker Compose.

## Features

- Create, retrieve, update, complete, and delete todos.
- PostgreSQL database integration.
- HTTP API health and service information endpoint.
- Automated Go tests.
- Multi-stage Docker build with tests executed before compilation.

## Technologies

- Go
- PostgreSQL
- Docker
- Docker Compose

## Prerequisites

Install the following tools:

- Git
- Docker Engine
- Docker Compose (Docker Compose v2 or compatible)

Verify the Docker installation:

```bash
docker --version
docker compose version
```

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/VASANTHANUMANDLAKADI/todo.git
cd todo
```

Run these instructions after the Docker changes have been merged into the repository's default branch. Otherwise, check out the branch containing the Docker files.

### 2. Configure environment variables

Create your local environment file from the example:

```bash
cp .env.example .env
```

Open `.env` and change `POSTGRES_PASSWORD` to a secure password.

Update `DATABASE_URL` to use the same username, password, and database name. Keep `db` as the hostname because it identifies the PostgreSQL service inside the Docker Compose network.

Do not commit `.env` or share real credentials.

### 3. Build and start the application

From the project root, run:

```bash
docker compose up --build -d
```

This command builds the Go API image, runs the Go tests during the image build, and starts the API and PostgreSQL containers in the background.

Check the container status:

```bash
docker compose ps
```

Wait until PostgreSQL is healthy and the API is running.

### 4. Verify the API

Check the service information endpoint:

```bash
curl -i http://localhost:8080/info
```

A successful response returns HTTP `200 OK`, with JSON similar to:

```json
{
  "service": "todo-api",
  "version": "1.0.0",
  "status": "running"
}
```

Retrieve todos:

```bash
curl -i http://localhost:8080/get/todos
```

The response contains the todos stored in PostgreSQL. An empty database may initially return an empty list or `null`, depending on the API implementation.

### 5. View application logs

View API logs:

```bash
docker compose logs api
```

Follow the logs continuously:

```bash
docker compose logs -f api
```

Press `Ctrl+C` to stop following the logs. This does not stop the containers.

### 6. Stop and restart the application

Stop the containers:

```bash
docker compose stop
```

Start the existing containers again:

```bash
docker compose start
```

Stop and remove the containers and Compose network:

```bash
docker compose down
```

The named PostgreSQL data volume is retained when using `docker compose down`, so database data can persist.

**Warning:** `docker compose down -v` removes the Compose-managed database volume and its stored data. Use it only when you intentionally want to delete that data.

## Run Tests

To run the Go tests locally, install a compatible Go version and run:

```bash
go test -count=1 ./...
```

The Docker build also runs the tests before compiling the application.

## Jenkins CI/CD Pipeline

The project uses Jenkins to automatically build, test, analyze, and deploy the Todo API when changes are pushed to the configured GitHub branch.

### Pipeline Prerequisites

- Jenkins installed and running.
- Git and Go available to the Jenkins process.
- Docker Engine and Docker Compose installed.
- Jenkins Pipeline, Git, JUnit, and GitHub integration plugins configured as needed.
- A GitHub webhook configured to trigger the Jenkins job on push events.
- PostgreSQL configuration and credentials managed securely outside source control.

### Pipeline Stages

The Jenkins pipeline performs the following steps:

1. **Checkout** — retrieves the source code from GitHub.
2. **Check Go Version** — verifies the Go installation.
3. **Download Dependencies** — downloads Go module dependencies.
4. **Compile** — runs `go build ./...`.
5. **Unit Tests** — runs Go tests and generates a JUnit-compatible XML report.
6. **Static Analysis** — runs `go vet ./...`.
7. **Check Go Formatting** — checks that Go files are formatted with `gofmt`.
8. **Build Todo API** — compiles the application binary.
9. **Archive Build Artifact** — stores the `todo-api` binary as a Jenkins build artifact.
10. **Build Docker Image** — builds the application image.
11. **Deploy with Docker Compose** — updates the API container using Docker Compose.
12. **Smoke Test** — checks the `/info` endpoint to verify that the application responds.

The pipeline stops when a required stage fails, preventing later deployment stages from running after a failure.

### GitHub Integration

The Jenkins job is configured to build the `Create-Todo` branch. A GitHub push webhook triggers the pipeline when new changes are pushed to that branch.

Webhook authentication is configured in Jenkins. Keep webhook secrets and database credentials out of the Git repository.

### Verify a Jenkins Build

1. Open the Jenkins dashboard.
2. Select the `todo-api-ci` job.
3. Open the latest build.
4. Confirm that the build status is `SUCCESS`.
5. Open **Test Result** to review passed, failed, and skipped tests.
6. Open **Build Artifacts** to access the archived `todo-api` binary.
7. Confirm that the smoke test for `http://localhost:8080/info` passes.

A successful pipeline indicates that the required build, test, analysis, artifact, deployment, and smoke-test stages completed successfully.

## Docker Files

- `Dockerfile` — builds and tests the Go application, then creates a runtime image.
- `compose.yaml` — configures the API and PostgreSQL services.
- `.dockerignore` — excludes unnecessary files and environment secrets from the Docker build context.
- `.env.example` — documents the required environment variables with example values.
- `db/init.sql` — initializes the database table when PostgreSQL creates a new database volume.

## API Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/info` | Service information |
| POST | `/create/todo` | Create a todo |
| GET | `/get/todos` | Retrieve todos |
| PUT | `/update/todo` | Update a todo |
| POST | `/completed/todo` | Mark a todo completed |
| DELETE | `/delete/todo` | Delete a todo |

The API listens on port `8080`.

## Troubleshooting

Check container status:

```bash
docker compose ps -a
```

Inspect the logs:

```bash
docker compose logs
```

If the API cannot connect to PostgreSQL, verify the `.env` values, ensure the database service is healthy, and confirm that `DATABASE_URL` uses the Compose service hostname `db`.

If port `8080` is already in use, stop the conflicting process or change the host-side port mapping in `compose.yaml`.
