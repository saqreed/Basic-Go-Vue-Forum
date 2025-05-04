# Forum Application

A full-stack forum application built with Go (backend) and Vue.js (frontend).

## Features

- User authentication and authorization
- Post creation, editing, and deletion
- Comment system
- Real-time chat using WebSocket
- User profiles and statistics
- Admin panel for user management

## Tech Stack

### Backend
- Go
- PostgreSQL
- JWT for authentication
- WebSocket for real-time communication
- Gorilla Mux for routing

### Frontend
- Vue.js 3
- Vue Router
- Pinia for state management
- Axios for HTTP requests
- WebSocket for real-time chat

## Project Structure

```
forum-go-vue/
├── backend/           # Go backend
│   ├── cmd/          # Main application entry point
│   ├── internal/     # Internal packages
│   │   ├── auth/     # Authentication logic
│   │   ├── database/ # Database operations
│   │   ├── handlers/ # HTTP handlers
│   │   ├── middleware/# Middleware functions
│   │   └── models/   # Data models
│   └── go.mod        # Go module file
│
└── frontend/         # Vue.js frontend
    ├── src/          # Source files
    │   ├── components/# Vue components
    │   ├── stores/   # Pinia stores
    │   ├── views/    # Page components
    │   └── router/   # Vue Router configuration
    └── package.json  # NPM dependencies
```

## Setup and Installation

### Prerequisites
- Go 1.21 or later
- Node.js 18 or later
- PostgreSQL 15 or later

### Backend Setup
1. Navigate to the backend directory
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up the database:
   ```bash
   psql -U postgres -f internal/database/init.sql
   ```
4. Start the server:
   ```bash
   go run cmd/main.go
   ```

### Frontend Setup
1. Navigate to the frontend directory
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```

## API Endpoints

### Authentication
- POST /api/register - User registration
- POST /api/login - User login
- POST /api/change-password - Password change

### Posts
- GET /api/posts - Get all posts
- GET /api/posts/{id} - Get post by ID
- POST /api/posts - Create post
- PUT /api/posts/{id} - Update post
- DELETE /api/posts/{id} - Delete post

### Comments
- GET /api/posts/{post_id}/comments - Get post comments
- POST /api/posts/{post_id}/comments - Add comment
- PUT /api/comments/{id} - Update comment
- DELETE /api/comments/{id} - Delete comment

### Chat
- GET /ws/chat - WebSocket endpoint for real-time chat

### Profile
- GET /api/profile - Get user profile
- GET /api/profile/stats - Get user statistics
- PUT /api/profile/password - Change password

### Admin
- GET /api/admin/users - Get all users
- GET /api/admin/users/{id} - Get user by ID
- PUT /api/admin/users/{id} - Update user
- DELETE /api/admin/users/{id} - Delete user

## Security Features
- JWT-based authentication
- Password hashing
- Role-based access control
- CORS protection
- WebSocket authentication

## License
MIT 