# Broadcast API
A powerful and scalable email broadcasting system for managing campaigns, contacts, and sending both transactional and bulk emails.

![Architecture Diagram](https://github.com/user-attachments/assets/095b4a36-0976-4496-8ee2-c6b015df689b)

## Overview

Broadcast API is a comprehensive email management solution designed to handle various types of email communications. It provides robust functionality for creating and managing email campaigns, contact lists, and sending both individual transactional emails and bulk broadcasts.

## Features

- **User Authentication**: Secure registration and login system
- **Campaign Management**: Create, retrieve, update, and delete email campaigns
- **Contact Management**: Organize and manage your contact lists
- **Broadcast Management**: Schedule and send email broadcasts to targeted recipients
- **Email Types**:
  - Transactional emails (individual, event-triggered)
  - Test emails (for verification purposes)
  - Campaign emails (to specific segments)
  - Bulk emails (to large recipient lists)
- **Health Monitoring**: System health checks and diagnostics

## API Endpoints

### Authentication
- `POST /register` - Create a new user account
- `POST /login` - Authenticate and receive access token

### Health
- `GET /healthz` - Check system health status

### Campaigns
- `POST /create Campaigns` - Create a new email campaign
- `GET /get Campaigns` - List all campaigns
- `GET /get Campaign` - Retrieve a specific campaign by ID
- `DEL /delete Campaign` - Remove a campaign

### Contacts
- `POST /create Contact` - Add a new contact
- `POST /get Contacts` - List all contacts (with optional filtering)
- `GET /get Contact` - Retrieve a specific contact by ID
- `PUT /update Contact` - Update contact information
- `DEL /delete Contact` - Remove a contact

### Broadcasts
- `POST /create Broadcast` - Create a new email broadcast
- `GET /get Broadcast` - Retrieve a specific broadcast by ID
- `GET /get Broadcasts` - List all broadcasts (with optional filtering)
- `PUT /update Broadcast` - Update broadcast details
- `POST /send Broadcast` - Execute sending of a broadcast
- `DEL /delete Broadcast` - Remove a broadcast

### Email Operations
- `POST /send Test Email` - Send a test email
- `POST /send Transactional Email` - Send a single transactional email
- `POST /process Email Job` - Process email job from queue
- `POST /send Campaign Emails` - Send emails for a specific campaign
- `POST /send Bulk Emails` - Send emails to a large list of recipients

## Architecture

The application follows a clean architecture pattern with clear separation of concerns:

- **API Layer**: HTTP handlers and route definitions
- **Service Layer**: Business logic implementation
- **Repository Layer**: Data access and persistence
- **Worker Layer**: Background processing for email sending
- **Scheduler**: Handles timed and recurring tasks

## Technology Stack

- **Backend**: Go (Golang)
- **Database**: PostgreSQL
- **Email Delivery**: SMTP integration
- **Containerization**: Docker
- **Orchestration**: Kubernetes and Helm
- **Local Development**: Kind (Kubernetes in Docker)

## Installation

### Prerequisites
- Go 1.22 or later
- Docker
- kubectl (Kubernetes command-line tool)
- Helm
- Kind (for local Kubernetes development)

### Using Docker Compose (Development)
```bash
# Clone the repository
git clone https://github.com/MdSadiqMd/Broadcast-API.git
cd broadcast-api

# Start the application using Docker Compose
docker-compose up -d
```

### Using Kubernetes and Helm (Production/Staging)
```bash
# Clone the repository
git clone https://github.com/MdSadiqMd/Broadcast-API.git
cd broadcast-api

# Start a local Kubernetes cluster
kind create cluster --name broadcast-cluster

# Build and load the Docker image
docker build -t broadcast-api:latest .
kind load docker-image broadcast-api:latest --name broadcast-cluster

# Deploy using Helm
helm upgrade --install broadcast-api ./helm/broadcast-api -f ./helm/broadcast-api/values.secret.yaml

# Access the application (port forwarding)
kubectl port-forward svc/broadcast-api 3000:80

# Delete everything when done
helm uninstall broadcast-api
kind delete cluster --name broadcast-cluster
```

### Manual Setup
```bash
# Clone the repository
git clone https://github.com/MdSadiqMd/Broadcast-API.git
cd broadcast-api

# Install dependencies
go mod download

# Configure the application
cp pkg/config/config.example.yaml pkg/config/config.yaml
# Edit config.yaml with your settings

# Run the application
go run cmd/server/main.go
```

## Configuration

The application is configured via `pkg/config/config.yaml`. Key configuration options include:

- Database connection details
- SMTP server settings
- Authentication settings
- Worker and scheduler configurations

When deploying with Kubernetes, configuration is managed through Helm values and Kubernetes ConfigMaps/Secrets.

## Kubernetes Deployment

The application can be deployed to any Kubernetes cluster using Helm. The Helm chart in the `helm/broadcast-api` directory contains all necessary Kubernetes manifests:

- Deployment
- Service
- ConfigMap
- Secrets
- Ingress (for external access)

### Key Kubernetes Commands

```bash
# Check pod status
kubectl get pods

# View application logs
kubectl logs -l app=broadcast-api

# Access the application (port forwarding)
kubectl port-forward svc/broadcast-api 3000:80

# Update after code changes
# 1. Rebuild image
docker build -t broadcast-api:latest .
# 2. Load to cluster
kind load docker-image broadcast-api:latest --name broadcast-cluster
# 3. Update deployment
helm upgrade broadcast-api ./helm/broadcast-api -f ./helm/broadcast-api/values.secret.yaml

# Delete everything when done
helm uninstall broadcast-api
kind delete cluster --name broadcast-cluster
```

## Development

### Project Structure
```
.
├── cmd                     # Application entry points
│   └── server              
├── internal                # Private application code
│   ├── api                 # API handlers and routes
│   ├── models              # Data models
│   ├── repositories        # Data access layer
│   ├── scheduler           # Scheduling components
│   ├── services            # Business logic
│   └── workers             # Background processing
├── pkg                     # Public libraries
│   ├── config              # Configuration handling
│   ├── email               # Email utilities
│   └── utils               # General utilities
├── helm                    # Helm chart for Kubernetes deployment
│   └── broadcast-api       # Application-specific chart
├── k8s                     # Kubernetes manifests (alternative to Helm)
│   ├── manifests           # Resource definitions
│   └── secrets             # Secret management scripts
└── Dockerfile              # Container definition
```

### Build and Test
```bash
# Build the application
go build -o broadcast-api cmd/server/main.go

# Build Docker image
docker build -t broadcast-api:latest .
```

## License

This project is licensed under the terms found in the LICENSE file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
