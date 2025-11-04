# Notex

A modern note-taking application with AI capabilities.

## Features

- User authentication and authorization
- Create, edit, and manage notes
- Draft system for work-in-progress notes
- Categories and tags for organization
- Comments and notifications
- AI-powered writing assistance
- Image generation
- Markdown support

## Tech Stack

### Backend
- Go
- PostgreSQL/MySQL database
- RESTful API

### Frontend
- Vue.js
- Vite

## Getting Started

### Prerequisites
- Go 1.x or higher
- Node.js and npm
- Database (PostgreSQL or MySQL)

### Installation

1. Clone the repository
2. Set up the backend:
   ```bash
   cd backend
   cp config.example.yaml config/config.yaml
   # Edit config/config.yaml with your settings
   go mod download
   go run main.go
   ```

3. Set up the frontend:
   ```bash
   cd frontend/web
   npm install
   npm run dev
   ```

## Docker Support

You can also run the application using Docker:

```bash
docker-compose up
```

## Configuration

See `backend/config.example.yaml` for configuration options.

## License

See LICENSE file for details.

---

Powered by SWE-agent
