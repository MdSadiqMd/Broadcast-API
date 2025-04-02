# Broadcast-API

Building this  
https://resend.com/blog/broadcast-api  
https://resend.com/docs/api-reference/audiences/create-audience

Architecture as of now
![image](https://github.com/user-attachments/assets/15322162-06fb-4f59-99b1-a556871642e2)

folder structure for now
```
broadcast-clone/
├── cmd/
│   └── server/
│       └── main.go             # Application entry point
├── config/
│   ├── config.go               # Configuration loading and validation
│   └── config.yaml             # Default configuration file
├── internal/
│   ├── api/                    # HTTP API layer
│   │   ├── handlers/           # Request handlers
│   │   │   ├── campaigns.go
│   │   │   ├── subscribers.go
│   │   │   ├── templates.go
│   │   │   └── lists.go
│   │   ├── middleware/         # Custom middleware
│   │   │   ├── auth.go
│   │   │   ├── logging.go
│   │   │   └── ratelimit.go
│   │   └── routes.go           # Route definitions
│   ├── models/                 # Database models
│   │   ├── campaign.go
│   │   ├── subscriber.go
│   │   ├── list.go
│   │   ├── template.go
│   │   └── message.go
│   ├── repositories/           # Database access layer
│   │   ├── campaign_repo.go
│   │   ├── subscriber_repo.go
│   │   └── template_repo.go
│   ├── services/               # Business logic
│   │   ├── campaign_service.go
│   │   ├── mail_service.go
│   │   ├── subscriber_service.go
│   │   └── template_service.go
│   ├── queue/                  # Email queuing and throttling
│   │   ├── manager.go          # Queue operations
│   │   └── throttler.go        # Rate limiting
│   ├── scheduler/              # Job scheduling
│   │   └── scheduler.go        # Scheduled tasks
│   └── workers/                # Email sending workers
│       └── mail_worker.go      # Worker implementation
├── migrations/                 # Database migrations
│   ├── 000001_create_tables.up.sql
│   └── 000001_create_tables.down.sql
├── pkg/                        # Reusable packages
│   ├── email/                  # Email utilities
│   │   └── smtp.go             # SMTP client wrapper
│   └── utils/                  # Shared utilities
│       ├── logger.go           # Logging utilities
│       └── validator.go        # Input validation
├── scripts/                    # Build and deployment scripts
│   ├── run.sh
│   └── migrations.sh
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── Dockerfile                  # Docker configuration
└── README.md                   # Project documentation
```
