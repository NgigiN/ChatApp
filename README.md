# ChatApp

A real-time chat application built with Go and WebSockets, designed for educational environments. Features room-based messaging, user authentication,  a modern responsive web interface.

## Features

- **Real-time Messaging**: WebSocket-based instant messaging
- **Room Management**: Create and join chat rooms for different subjects
- **User Authentication**: Secure login and registration system
- **Responsive Design**: Mobile-first UI with modern styling
- **Database Persistence**: MySQL for storing messages and user data
- **Docker Support**: Easy deployment with Docker Compose
- **RESTful API**: Complete API for room and message management

## Tech Stack

### Backend
- **Go 1.23** - Main application language
- **Gin** - HTTP web framework
- **WebSockets** - Real-time communication
- **MySQL** - Database for persistence
- **Redis** - Caching and session management
- **Prometheus** - Metrics and monitoring

### Frontend
- **Vanilla JavaScript** - Client-side functionality
- **HTML5/CSS3** - Modern responsive design
- **Font Awesome** - Icons
- **Google Fonts** - Typography

## Quick Start

### Prerequisites
- Go 1.23+
- Docker and Docker Compose
- MySQL 8.4+
- Redis 7+

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone <repository-url>
cd ChatApp
```

2. Start the application:
```bash
docker-compose up -d
```

3. Access the application:
- Web interface: http://localhost:8000
- API: http://localhost:8000/api

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Start MySQL and Redis services

3. Update database connection in `main.go` if needed

4. Run the application:
```bash
go run main.go
```

## API Endpoints

### Authentication
- `POST /login` - User login
- `POST /register` - User registration

### Rooms
- `GET /api/rooms` - List all rooms
- `POST /api/rooms` - Create a new room
- `GET /api/rooms/:id/messages` - Get room messages
- `POST /api/rooms/:id/messages` - Send message to room

### WebSocket
- `GET /ws` - WebSocket connection for real-time chat

## Project Structure

```
ChatApp/
├── cmd/server/          # Application entry point
├── internal/
│   ├── chat/           # Chat server and WebSocket handling
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # Authentication, logging, rate limiting
│   ├── models/         # Data models
│   ├── repositories/   # Database access layer
│   ├── services/       # Business logic
│   └── user/           # User management
├── static/             # Frontend assets
│   ├── css/           # Stylesheets
│   ├── js/            # JavaScript files
│   └── *.html         # HTML pages
├── pkg/               # Shared packages
├── docker-compose.yml # Docker configuration
└── Dockerfile         # Container build instructions
```

## Default Rooms

The application comes with pre-configured rooms:
- Math 101 - Calculus
- Physics Lab
- Chemistry Study Group
- General Discussion

## Configuration

Environment variables:
- `SERVER_PORT` - Server port (default: 8000)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `REDIS_HOST` - Redis host
- `REDIS_PORT` - Redis port
- `JWT_SECRET` - JWT signing secret

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o chatapp ./cmd/server
```

### Database Migrations
The application automatically creates required tables on startup.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
