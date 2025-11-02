package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

type Target struct {
	Description          string `json:"description"`
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

type Command struct {
	ID     int                    `json:"id"`
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: chrome-navigate <song|artist|album|lyrics>")
		os.Exit(1)
	}
	view := os.Args[1]
	validViews := map[string]bool{"song": true, "artist": true, "album": true, "lyrics": true}
	if !validViews[view] {
		log.Fatalf("Invalid view: %s (valid: song, artist, album, lyrics)", view)
	}

	// Get the list of Chrome targets
	resp, err := http.Get("http://localhost:9222/json")
	if err != nil {
		log.Fatalf("Cannot connect to Chrome remote debugger: %v", err)
	}
	defer resp.Body.Close()

	var targets []Target
	if err := json.NewDecoder(resp.Body).Decode(&targets); err != nil {
		log.Fatalf("Cannot parse Chrome /json response: %v", err)
	}

	if len(targets) == 0 {
		log.Fatalf("No open Chrome tabs found. Is Chrome running with --remote-debugging-port=9222?")
	}

	// Pick the first tab (or filter for Now Playing)
	target := targets[0]
	fmt.Printf("🔗 Connecting to tab: %s (%s)\n", target.Title, target.WebSocketDebuggerURL)

	// Connect via WebSocket
	ws, _, err := websocket.DefaultDialer.Dial(target.WebSocketDebuggerURL, nil)
	if err != nil {
		log.Fatalf("WebSocket connection failed: %v", err)
	}
	defer ws.Close()

	// Build the navigate command
	cmd := Command{
		ID:     1,
		Method: "Page.navigate",
		Params: map[string]interface{}{
			"url": fmt.Sprintf("http://localhost:4004/?view=%s", view),
		},
	}

	data, _ := json.Marshal(cmd)
	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Fatalf("Failed to send navigation command: %v", err)
	}

	fmt.Printf("✅ Navigated to NowPlaying view: %s\n", view)
}
