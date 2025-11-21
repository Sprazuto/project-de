# Project DE Docker Setup

This document provides comprehensive instructions for containerizing and deploying the Project DE (Project Dashboard Executive) application using Docker and microservices architecture. The application serves as a government procurement performance monitoring system with real-time dashboards and data visualization.

## Overview

The Project DE application is deployed as a microservices architecture consisting of:

- **API Service**: Go/Gin backend handling authentication, business logic, and data operations
- **Web Service**: Vue.js frontend serving the user interface
- **Database Service**: PostgreSQL for primary data storage
- **Cache Service**: Redis for session management and caching
- **Reverse Proxy**: Nginx for load balancing and SSL termination

## Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   Nginx Proxy   │    │   Web Service   │
│   (Port 80/443) │◄──►│   Vue.js App    │
│                 │    │   (Port 80)     │
└─────────────────┘    └─────────────────┘
         │
         ▼
┌─────────────────┐    ┌─────────────────┐
│   API Service   │    │  Database       │
│   Go/Gin App    │◄──►│  PostgreSQL     │
│   (Port 8080)   │    │   Redis Cache   │
└─────────────────┘    └─────────────────┘
```

## Prerequisites

- Docker Engine 20.10+
- Docker Compose 2.0+
- Git
- At least 4GB RAM available
- 10GB free disk space

## Quick Start

### 1. Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd project-de

# Copy environment files
cp .env.example .env

# Edit environment variables
nano .env
```

### 2. Environment Configuration

Edit `docker/.env`:

```env
# Application Configuration
COMPOSE_PROJECT_NAME=project_de
APP_ENV=production

# Database Configuration
POSTGRES_DB=project_de_db
POSTGRES_USER=project_de_user
POSTGRES_PASSWORD=your_secure_db_password_here
POSTGRES_HOST=database
POSTGRES_PORT=5432

# Redis Configuration
REDIS_HOST=cache
REDIS_PORT=6379
REDIS_PASSWORD=your_secure_redis_password_here

# API Service Configuration
ENV=production
API_HOST=0.0.0.0
API_PORT=8080
DB_USER=project_de_user
DB_PASS=your_secure_db_password_here
DB_NAME=project_de_db
DB_HOST=database
DB_PORT=5432
REDIS_HOST=cache
REDIS_PORT=6379
JWT_ACCESS_SECRET=your_jwt_access_secret_here
JWT_REFRESH_SECRET=your_jwt_refresh_secret_here

# Web Service Configuration
VITE_API_BASE_URL=http://api:8080/api
VITE_APP_TITLE=Project Dashboard Executive
NODE_ENV=production

# SSL Configuration (Optional)
SSL_CERT_PATH=/etc/ssl/certs/project.crt
SSL_KEY_PATH=/etc/ssl/private/project.key
```

### 3. Launch Application

```bash
# Build and start all services
docker-compose up -d --build

# View logs
docker-compose logs -f

# Check service status
docker-compose ps
```

### 4. Access Application

- **Web Interface**: http://localhost
- **API Documentation**: http://localhost/api/swagger/index.html
- **Health Check**: http://localhost/api/health

## Service Configuration

### API Service (Go/Gin)

**Dockerfile** (`api/Dockerfile`):

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/db ./db
CMD ["./main"]
```

**Configuration**:

- Environment variables loaded from `.env`
- Database migrations run on startup
- Admin user seeded automatically
- Health check endpoint available

### Web Service (Vue.js)

**Dockerfile** (`web/Dockerfile`):

```dockerfile
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

**Nginx Configuration** (`web/nginx.conf`):

```nginx
server {
    listen 80;
    server_name localhost;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://api:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_For;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;
}
```

### Database Service (PostgreSQL)

**Configuration**:

- Persistent data volume: `sijagur_postgres_data`
- Automatic database initialization
- Connection pooling via PgBouncer (optional)

### Cache Service (Redis)

**Configuration**:

- Persistent data volume: `sijagur_redis_data`
- Password authentication enabled
- Memory limits and eviction policies configured

### Reverse Proxy (Nginx)

**Features**:

- SSL/TLS termination
- Load balancing across multiple API instances
- Static file serving with caching
- Security headers and CORS configuration

## Docker Compose Configuration

### Main Compose File (`docker-compose.yml`)

```yaml
version: "3.8"

services:
  # Reverse Proxy
  proxy:
    image: nginx:alpine
    container_name: project_de_proxy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/ssl/certs:ro
    depends_on:
      - web
      - api
    networks:
      - project_de_network
    restart: unless-stopped

  # Web Service (Vue.js Frontend)
  web:
    build:
      context: ./web
      dockerfile: Dockerfile
    container_name: project_de_web
    environment:
      - VITE_API_BASE_URL=http://api:8080/api
      - VITE_APP_TITLE=Project Dashboard Executive
      - NODE_ENV=production
    depends_on:
      - api
    networks:
      - project_de_network
    restart: unless-stopped

  # API Service (Go/Gin Backend)
  api:
    build:
      context: ./api
      dockerfile: Dockerfile
    container_name: project_de_api
    environment:
      - ENV=production
      - API_HOST=0.0.0.0
      - API_PORT=8080
    env_file:
      - .env
    depends_on:
      - database
      - cache
    networks:
      - project_de_network
    restart: unless-stopped

  # PostgreSQL Database
  database:
    image: postgres:15-alpine
    container_name: project_de_database
    environment:
      - POSTGRES_DB=project_de_db
      - POSTGRES_USER=project_de_user
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./api/db/database.sql:/docker-entrypoint-initdb.d/01-init.sql:ro
      - ./api/db/db_sijagur.sql:/docker-entrypoint-initdb.d/02-sijagur.sql:ro
    networks:
      - project_de_network
    restart: unless-stopped

  # Redis Cache
  cache:
    image: redis:7-alpine
    container_name: project_de_cache
    command: redis-server --requirepass ${REDIS_PASSWORD}
    volumes:
      - redis_data:/data
    networks:
      - project_de_network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  project_de_network:
    driver: bridge
```

## Environment Management

### Development Environment

Create `docker-compose.dev.yml`:

```yaml
version: "3.8"

services:
  api:
    build:
      context: ./api
      dockerfile: Dockerfile.dev
    volumes:
      - ./api:/app
    environment:
      - ENV=development
      - API_PORT=8080
    command: go run main.go

  web:
    build:
      context: ./web
      dockerfile: Dockerfile.dev
    volumes:
      - ./web:/app
      - /app/node_modules
    ports:
      - "5173:5173"
    environment:
      - NODE_ENV=development
    command: npm run dev -- --host 0.0.0.0

  database:
    ports:
      - "5432:5432"

  cache:
    ports:
      - "6379:6379"
```

### Production Environment

```yaml
# docker/docker-compose.prod.yml
version: "3.8"

services:
  api:
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: "1.0"
          memory: 1G
        reservations:
          cpus: "0.5"
          memory: 512M

  web:
    deploy:
      replicas: 2

  proxy:
    ports:
      - "80:80"
      - "443:443"
```

## Networking

### Internal Network

All services communicate through `project_de_network`:

- **web** ↔ **api**: HTTP communication for API calls using useApi composable
- **api** ↔ **database**: PostgreSQL connections for data persistence
- **api** ↔ **cache**: Redis connections for session management
- **proxy** ↔ **web**: Static file serving with caching
- **proxy** ↔ **api**: API proxying with load balancing

### External Access

- **Port 80/443**: Public access through Nginx proxy
- **Health checks**: Internal service monitoring endpoints
- **Database**: Not exposed externally (security best practice)

## Security Best Practices

### Container Security

```dockerfile
# Use non-root user
FROM alpine:latest
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001
USER nextjs

# Minimize attack surface
RUN apk add --no-cache libc6-compat
```

### Environment Security

```bash
# Generate secure secrets
openssl rand -base64 32

# Store secrets securely
echo "JWT_SECRET=$(openssl rand -base64 32)" >> .env
echo "DB_PASSWORD=$(openssl rand -base64 16)" >> .env
```

### Network Security

```yaml
# docker-compose.yml
services:
  api:
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
```

## Scaling and Performance

### Horizontal Scaling

```bash
# Scale API service
docker-compose up -d --scale api=3

# Scale web service
docker-compose up -d --scale web=2
```

### Load Balancing

```nginx
upstream api_backend {
    server api-1:8080;
    server api-2:8080;
    server api-3:8080;
}

server {
    location /api {
        proxy_pass http://api_backend;
    }
}
```

### Database Optimization

```yaml
database:
  deploy:
    resources:
      limits:
        cpus: "2.0"
        memory: 4G
      reservations:
        cpus: "1.0"
        memory: 2G
```

## Monitoring and Logging

### Health Checks

```yaml
services:
  api:
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### Logging

```yaml
services:
  api:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Monitoring Stack

```yaml
# docker-compose.monitoring.yml
version: "3.8"

services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./monitoring/prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    volumes:
      - grafana_data:/var/lib/grafana

  loki:
    image: grafana/loki
    volumes:
      - ./monitoring/loki-config.yml:/etc/loki/local-config.yaml
```

## Backup and Recovery

### Database Backup

```bash
# Backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
docker exec sijagur_database pg_dump -U sijagur_user sijagur_db > backup_$DATE.sql
```

### Automated Backups

```yaml
# docker-compose.backup.yml
version: "3.8"

services:
  backup:
    image: postgres:15-alpine
    volumes:
      - ./backups:/backups
      - postgres_data:/var/lib/postgresql/data:ro
    command: >
      bash -c "
      pg_dump -h database -U sijagur_user sijagur_db > /backups/backup_$(date +%Y%m%d_%H%M%S).sql
      "
```

## Troubleshooting

### Common Issues

#### Service Won't Start

```bash
# Check logs
docker-compose logs <service_name>

# Check service status
docker-compose ps

# Restart service
docker-compose restart <service_name>
```

#### Database Connection Issues

```bash
# Check database connectivity
docker-compose exec api nc -zv database 5432

# Check database logs
docker-compose logs database
```

#### API Not Responding

```bash
# Test API endpoint
curl http://localhost/api/health

# Check API logs
docker-compose logs api
```

### Performance Issues

```bash
# Monitor resource usage
docker stats

# Check container logs for errors
docker-compose logs --tail=100 -f
```

### Network Issues

```bash
# Test inter-service communication
docker-compose exec web ping api

# Check network configuration
docker network ls
docker network inspect sijagur_sijagur_network
```

## Deployment Strategies

### Blue-Green Deployment

```bash
# Deploy new version
docker-compose -f docker-compose.green.yml up -d

# Switch traffic (update nginx config)
docker-compose -f docker-compose.blue.yml down

# Rename green to blue
mv docker-compose.green.yml docker-compose.blue.yml
```

### Rolling Updates

```bash
# Update with zero downtime
docker-compose up -d --no-deps api

# Check health
curl http://localhost/api/health

# Scale back if needed
docker-compose up -d --scale api=3
```

### CI/CD Integration

```yaml
# .github/workflows/deploy.yml
name: Deploy to Production
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build and push images
        run: |
          echo ${{ secrets.DOCKER_PASSWORD }} | docker login -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
          docker build -t myregistry/project-de-api:${{ github.sha }} ./api
          docker build -t myregistry/project-de-web:${{ github.sha }} ./web
          docker push myregistry/project-de-api:${{ github.sha }}
          docker push myregistry/project-de-web:${{ github.sha }}

      - name: Deploy to server
        run: |
          ssh user@server << EOF
            cd /opt/project-de
            sed -i 's|image:.*api.*|image: myregistry/project-de-api:${{ github.sha }}|' docker-compose.yml
            sed -i 's|image:.*web.*|image: myregistry/project-de-web:${{ github.sha }}|' docker-compose.yml
            docker-compose pull
            docker-compose up -d
          EOF
```

## Maintenance

### Updates

```bash
# Update all images
docker-compose pull

# Update specific service
docker-compose pull api
docker-compose up -d api

# Update with zero downtime
docker-compose up -d --no-deps api
```

### Cleanup

```bash
# Remove unused images
docker image prune -f

# Remove unused volumes
docker volume prune -f

# Remove stopped containers
docker container prune -f
```

### Log Rotation

```bash
# Rotate logs
docker-compose exec api logrotate /etc/logrotate.d/api

# View log files
docker-compose exec api tail -f /var/log/api/app.log
```

## API Integration Architecture

The Project DE application implements a sophisticated API integration system that seamlessly connects the Vue.js frontend with the Go/Gin backend:

### Frontend Integration (Vue.js + useAuth/useApi)

- **useAuth Composable**: Manages JWT authentication state and token refresh
- **useApi Composable**: Handles HTTP requests with automatic authentication headers
- **useDashboard Composable**: Orchestrates complex data fetching for procurement dashboards
- **Reactive State Management**: Vue 3 Composition API for real-time UI updates

### Backend Services (Go/Gin)

- **Authentication Service**: JWT token generation and validation
- **Data Service**: Procurement performance data aggregation
- **API Gateway**: RESTful endpoints for frontend consumption
- **Database Layer**: PostgreSQL with GORM for data persistence
- **Cache Layer**: Redis for session management and performance optimization

### Data Flow Architecture

```
Frontend (Vue.js)
    ↓ useApi composable
API Service (Go/Gin)
    ↓ Business Logic
Database (PostgreSQL) + Cache (Redis)
    ↓ Data Processing
Frontend Dashboard (Real-time Updates)
```

## Conclusion

This Docker setup provides a production-ready, scalable, and secure deployment for the Project DE application. The microservices architecture ensures maintainability, while Docker containers provide consistency across development and production environments.

The integration leverages modern Vue 3 patterns with useAuth and useApi composables for seamless API communication, JWT authentication, and real-time dashboard data management for government procurement performance monitoring.

For additional support or questions, refer to the API_INTEGRATION_GUIDE.md, FRONTEND_DOCUMENTATION.md, and API_DOCUMENTATION.md files or contact the development team.
