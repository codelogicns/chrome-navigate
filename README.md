# 🎶 chrome-navigate

A simple Go utility to remotely control a **Chromium kiosk** (like Volumio’s NowPlaying screen)  
using Chrome’s **DevTools WebSocket protocol** via `--remote-debugging-port`.

This lets you change the NowPlaying plugin’s view — for example:  
`song`, `artist`, `album`, or `lyrics` — from the command line or another script.

---

## 🚀 Features

- Connects to Chrome’s remote debugging WebSocket (`--remote-debugging-port=9222`)
- Sends `Page.navigate` commands via the Chrome DevTools Protocol
- Automatically navigates your kiosk to a specific NowPlaying view:
  - `song`
  - `artist`
  - `album`
  - `lyrics`

---

## 🧱 Installation

### 1. Clone or copy the source

```bash
git clone https://github.com/codelogicns/chrome-navigate.git
cd chrome-navigate

go mod init chrome-navigate
go get github.com/gorilla/websocket
go mod tidy

go build chrome-navigate.go
```


## 🧠 Usage

1. Start Chromium / Chrome with remote debugging enabled

    On your kiosk or Volumio system:
```bash
chromium-browser \
  --kiosk http://localhost:4004 \
  --remote-debugging-port=9222 \
  --remote-debugging-address=0.0.0.0 &
```
2. Run the Go program
```bash
./chrome-navigate <view>
```

Where `view` can be one of:

| View     | Description                |
| -------- | -------------------------- |
| `song`   | Show the current song info |
| `artist` | Show artist details        |
| `album`  | Show album view            |
| `lyrics` | Show song lyrics           |


## 🧩 How It Works

The program calls http://localhost:9222/json to list all open Chrome targets.

It connects to the first one’s webSocketDebuggerUrl.

It sends this Chrome DevTools command:

```json
{
  "id": 1,
  "method": "Page.navigate",
  "params": {
    "url": "http://localhost:4004/?view=lyrics"
  }
}
```

## ⚙️ Requirements

- Go 1.19+
- Chrome or Chromium started with --remote-debugging-port
- The NowPlaying plugin’s web UI accessible at http://localhost:4004

## Tip

Tip: For development, you can see Chrome’s JSON targets:

```bash
curl http://localhost:9222/json | jq .
```