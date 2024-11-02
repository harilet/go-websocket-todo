# go-websocket-todo

This is a simple Todo App built with Go and Svelte that uses WebSocket for real-time communication between the client and the server.

## Features
- Add, remove, and mark as done Todo items
- Real-time updates of Todo items using WebSocket
- Simple and easy-to-use UI

## Getting Started

1. Clone the repository
    ```
    git clone https://github.com/harilet/go-websocket-todo.git
    ```
2. Installing the dependencies and building the ui

    ```
        cd go-websocket-todo
        cd web
        npm install
        npm run build
    ```

3. Run the server

    ```
    cd go-websocket-todo
    go run .
    ```

4. Open the client in your browser at http://localhost:3333