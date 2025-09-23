package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Version information (set during build)
var version = "1.0.2"

// Message represents a Gotify message
type GotifyMessage struct {
	Id       uint32 `json:"id"`
	Appid    uint32 `json:"appid"`
	Message  string `json:"message"`
	Title    string `json:"title"`
	Priority uint32 `json:"priority"`
	Date     string `json:"date"`
}

// Forwarder handles message forwarding
type Forwarder struct {
	ws          *websocket.Conn
	debugLogger *log.Logger
	targetURL   string
	gotifyHost  string
	gotifyToken string
	iconURL     string
}

func (f *Forwarder) sendMessage(msg *GotifyMessage) error {
	// Log the received message
	f.debugLogger.Printf("Received message - Title: %s, Content: %s\n", msg.Title, msg.Message)

	// Create the message payload for Bark server
	payload := map[string]interface{}{
		"title":      msg.Title,
		"body":       msg.Message,
		"badge":      1,
		"sound":      "minuet",
		"group":      "Gotify",
		"icon":       f.iconURL,
		"url":        f.gotifyHost,
		"device_key": strings.TrimPrefix(f.targetURL, "https://"),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Log the payload being sent
	f.debugLogger.Printf("Sending payload: %s\n", string(payloadBytes))

	// Create HTTP request
	req, err := http.NewRequest("POST", f.targetURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	f.debugLogger.Printf("Successfully forwarded message: %s\n", msg.Title)
	return nil
}

func (f *Forwarder) connectWebSocket() error {
	// Convert http(s):// to ws(s)://
	wsURL := f.gotifyHost
	if strings.HasPrefix(wsURL, "http://") {
		wsURL = "ws://" + wsURL[7:]
	} else if strings.HasPrefix(wsURL, "https://") {
		wsURL = "wss://" + wsURL[8:]
	}

	wsURL = fmt.Sprintf("%s/stream?token=%s", wsURL, f.gotifyToken)
	f.debugLogger.Printf("Connecting to WebSocket URL: %s\n", wsURL)

	// Create a custom dialer that skips TLS verification
	dialer := websocket.Dialer{
		EnableCompression: true,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	// Add custom headers
	headers := http.Header{}
	headers.Add("User-Agent", "Gotify-Forwarder/1.0")

	ws, resp, err := dialer.Dial(wsURL, headers)
	if err != nil {
		if resp != nil {
			f.debugLogger.Printf("WebSocket handshake failed with status: %d\n", resp.StatusCode)
			body, _ := io.ReadAll(resp.Body)
			f.debugLogger.Printf("Response body: %s\n", string(body))
		}
		return fmt.Errorf("failed to connect to websocket: %v", err)
	}

	f.ws = ws
	f.debugLogger.Println("Successfully connected to Gotify websocket")
	return nil
}

func (f *Forwarder) start() {
	for {
		if f.ws == nil {
			if err := f.connectWebSocket(); err != nil {
				f.debugLogger.Printf("Connection error: %v\n", err)
				time.Sleep(5 * time.Second)
				continue
			}
		}

		msg := &GotifyMessage{}
		err := f.ws.ReadJSON(msg)
		if err != nil {
			f.debugLogger.Printf("Error reading message: %v\n", err)
			f.ws.Close()
			f.ws = nil
			time.Sleep(5 * time.Second)
			continue
		}

		// Log the complete message details
		f.debugLogger.Printf("New message received:\n"+
			"  ID: %d\n"+
			"  App ID: %d\n"+
			"  Title: %s\n"+
			"  Message: %s\n"+
			"  Priority: %d\n"+
			"  Date: %s\n",
			msg.Id, msg.Appid, msg.Title, msg.Message, msg.Priority, msg.Date)

		if err := f.sendMessage(msg); err != nil {
			f.debugLogger.Printf("Error forwarding message: %v\n", err)
		}
	}
}

func main() {
	// Parse command line flags
	gotifyHost := flag.String("host", "", "Gotify server host (e.g., http://localhost:8080)")
	gotifyToken := flag.String("token", "", "Gotify client token")
	targetURL := flag.String("target", "", "Target URL to forward messages to")
	iconURL := flag.String("icon", "https://day.app/assets/images/avatar.jpg", "Icon URL for notifications")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version if requested
	if *showVersion {
		fmt.Printf("Gotify-to-Bark version %s\n", version)
		os.Exit(0)
	}

	// Validate required parameters
	if *gotifyHost == "" || *gotifyToken == "" || *targetURL == "" {
		fmt.Println("Usage: gotify-forwarder -host <gotify-host> -token <gotify-token> -target <target-url> [-icon <icon-url>]")
		fmt.Printf("Version: %s\n", version)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create forwarder instance
	forwarder := &Forwarder{
		debugLogger: log.New(os.Stdout, "Gotify Forwarder: ", log.Lshortfile),
		targetURL:   *targetURL,
		gotifyHost:  *gotifyHost,
		gotifyToken: *gotifyToken,
		iconURL:     *iconURL,
	}

	// Start forwarding
	forwarder.start()
}
