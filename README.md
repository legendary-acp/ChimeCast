# ChimeCast
ChimeCast is a robust backend solution for a Zoom-like video calling application, crafted in Go. Utilizing the `mux` router, Gorilla WebSocket, and `pion/webrtc`, this backend offers seamless real-time communication through peer-to-peer (P2P) video calls. With features like simple session-based authentication and user management, ChimeCast is designed to provide a reliable foundation for high-quality video conferencing experiences.

## Features

- **Session-based Authentication**: Custom session tokens for secure access without OAuth.
- **User Management**: Register, login, and manage user data.
- **Real-time Communication**: WebSocket-based signaling for WebRTC connections.
- **Peer-to-Peer Video/Audio**: WebRTC support for high-quality media streaming.
- **SQLite Database**: Persistent storage for user data and session logs.

## Tech Stack

- **Language**: Go
- **Framework**: `mux` for routing and API handling
- **Database**: SQLite via `GORM` (ORM for Go)
- **Real-Time Communication**:
  - **WebSocket**: Gorilla WebSocket for signaling
  - **WebRTC**: `pion/webrtc` for P2P media streaming

---
## API Endpoints
### 1. Authentication

- `POST /register`<br>
Registers a new user.

```json

{
  "username": "example",
  "password": "password123"
}
```
- `POST /login`<br>
Logs in a user and returns a session token.

```json

    {
      "username": "example",
      "password": "password123"
    }
```
- `POST /logout`<br>
Ends a user session (requires session token).

### 2. Video Call Management

- `POST /call/start`<br>
    Initiates a call session with a target user. Requires authentication.

- `POST /call/end`<br>
    Ends an ongoing call.

### 3. WebSocket Connection (for signaling)

-`/ws` <br>
    Establishes a WebSocket connection for WebRTC signaling between peers.

---

## Libraries and Packages
- **mux**: HTTP request router and dispatcher.
- **Gorilla WebSocket**: For real-time signaling using WebSocket.
- **pion/webrtc**: WebRTC implementation for Go to support P2P media streaming.
- **GORM**: ORM for Go, simplifies database interactions with SQLite.
