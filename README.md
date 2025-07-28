# Multi-Room Messenger App

A Go-based chat application with multiple rooms, using WebSocket for real-time communication and Redpanda as the message broker.

## Features
- Multi-room chat with unique room IDs.
- WebSocket connections via `/:chatName/username?username=<username>`.
- Messages stored and broadcast via Redpanda topics (`room-<chatName>`).
- Simple HTML/JavaScript frontend for testing.

## Prerequisites
- Go 1.21+: [Install Go](https://go.dev/dl/)
- Docker: [Install Docker](https://www.docker.com/)
- Redpanda: Configured via Docker Compose

## Setup
1. **Clone the Repository**:
   ```bash
   git clone <repository-url>
   cd Ar-messenger
   ```

2. **Install Dependencies**:
   ```bash
   go mod tidy
   ```

3. **Configure Environment**:
   Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```
   and then set the correct address for configurations

4. **Start chat**:
   ```bash
   docker-compose up -d
   ```

5. **Access the App**:
   - Open `http://localhost:8080` in a browser.
   - Enter a username and chat room (e.g., `room1`).
   - Send and receive messages in real-time.

## Usage
- Join a room: Enter a username and room ID (e.g., `room1`) in the web interface.
- WebSocket URL: `ws://localhost:8080/<roomID>/username?username=<username>`.
- Messages are broadcast to all clients in the same room via Redpanda.

## Project Structure
```
Ar-messenger/
├── cmd/messenger/main.go         # Entry point
├── internal/                   # Application logic
├── web/                        # Frontend static files
├── .env.example                # Sample environment variables
├── .gitignore                  # Git ignore rules
├── docker-compose.yml          # Redpanda and chat setup
└── go.mod                      # Go module definition
```