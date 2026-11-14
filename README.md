# goTalk

A small Go application that serves a simple web chat interface and exposes a basic JSON API for sending messages.

## Overview

This project is a minimal prototype for message exchange in a browser. It serves static files from the `public/` folder and accepts HTTP requests at `/api/send`, returning the message back as JSON. The UI is intentionally lightweight and can be used as a starting point for a chat or messaging feature.

## Project structure

```text
.
├── LICENSE
├── README.md
├── go.mod
├── go.sum
├── main.go
├── public/
│   ├── chat.html
│   └── index.html
└──
```

## Features

- Go HTTP server
- Static frontend served from `/web/`
- JSON endpoint: `POST /api/send`
- Simple browser UI for sending and displaying messages
- Basic validation for empty messages

## How it works

The server does the following:

- serves files under `/web/` from the `public` directory
- receives JSON payloads such as:

```json
{
  "id": 1,
  "message": "Olá, mundo!"
}
```

- validates the incoming message
- returns a JSON response like:

```json
{
  "status": "success",
  "message": "Olá, mundo!"
}
```

## Run locally

1. Make sure Go is installed.
2. Install dependencies:

```bash
go mod download
```

3. Start the server:

```bash
go run .
```

4. Open the app in your browser:

- `http://localhost:8080/web/index.html`
- `http://localhost:8080/web/chat.html`

## API example

You can test the endpoint with curl:

```bash
curl -X POST http://localhost:8080/api/send \
  -H "Content-Type: application/json" \
  -d '{"id":1,"message":"hello from goTalk"}'
```

Example response:

```json
{"status":"success","message":"hello from goTalk"}
```

## Important note

The current implementation attempts to connect to a WebSocket at `ws://localhost:3000/ws` when the application starts. If no WebSocket server is running on that port, the app may terminate during startup. To keep the app running as a simple HTTP demo, remove or comment out the `initSocket()` call in `main.go`.

## Status

This app is a lightweight prototype and does not yet include:

- real-time multi-user chat
- persistent message storage
- authentication
- database integration
- websocket server implementation

It is best suited as a starting point for learning Go web apps and building a chat feature step by step.
