package api

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pikachuim/deskcli/internal/config"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	CheckOrigin:      func(r *http.Request) bool { return true },
}

type wsMsg struct {
	Type string `json:"type"` // stdin | resize
	Data string `json:"data"` // base64 for stdin
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

type wsReply struct {
	Type    string `json:"type"` // stdout | exit | error
	Data    string `json:"data"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// GET /api/v1/ws/terminal/:name
func TerminalWS(c *gin.Context) {
	name, valid := containerName(c)
	if !valid {
		return
	}
	if !authorizeContainer(c, name) {
		return
	}
	cfg, err := config.Load()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Build exec command for the container
	binary := cfg.Engine
	if binary != "docker" && binary != "podman" {
		sendWSError(conn, "WebSocket terminal only supported for docker/podman engines")
		return
	}
	cmd := exec.Command(binary, "exec", "-it", name, "/bin/bash")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		sendWSError(conn, "failed to open stdin: "+err.Error())
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		sendWSError(conn, "failed to open stdout: "+err.Error())
		return
	}
	cmd.Stderr = cmd.Stdout // merge stderr into stdout

	if err := cmd.Start(); err != nil {
		sendWSError(conn, "exec failed: "+err.Error())
		return
	}

	done := make(chan struct{}, 1)

	// stdout → WebSocket
	go func() {
		defer func() { done <- struct{}{} }()
		buf := make([]byte, 4096)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				reply := wsReply{Type: "stdout", Data: base64.StdEncoding.EncodeToString(buf[:n])}
				b, _ := json.Marshal(reply)
				if err2 := conn.WriteMessage(websocket.TextMessage, b); err2 != nil {
					return
				}
			}
			if err != nil {
				return
			}
		}
	}()

	// WebSocket → stdin
	go func() {
		defer func() {
			stdin.Close()
			done <- struct{}{}
		}()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var m wsMsg
			if err := json.Unmarshal(msg, &m); err != nil {
				continue
			}
			switch m.Type {
			case "stdin":
				decoded, err := base64.StdEncoding.DecodeString(m.Data)
				if err != nil {
					continue
				}
				if _, err := stdin.Write(decoded); err != nil {
					return
				}
			case "resize":
				// PTY resize not available without creack/pty; acknowledged but no-op
			}
		}
	}()

	// Wait for process exit or disconnect
	<-done
	exitCode := 0
	if state := cmd.ProcessState; state != nil {
		exitCode = state.ExitCode()
	}
	reply := wsReply{Type: "exit", Code: exitCode}
	b, _ := json.Marshal(reply)
	_ = conn.WriteMessage(websocket.TextMessage, b)
}

func sendWSError(conn *websocket.Conn, msg string) {
	reply := wsReply{Type: "error", Message: msg}
	b, _ := json.Marshal(reply)
	_ = conn.WriteMessage(websocket.TextMessage, b)
}
