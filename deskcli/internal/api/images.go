package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/pikachuim/deskcli/internal/config"
	"github.com/pikachuim/deskcli/internal/engine"
)

// GET /api/v1/images
func ListImages(c *gin.Context) {
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	images, err := eng.ListImages()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if images == nil {
		images = []engine.ImageInfo{}
	}
	ok(c, images)
}

type pullImageReq struct {
	Image string `json:"image" binding:"required"`
}

// POST /api/v1/images/pull  — streams pull progress via SSE
func PullImage(c *gin.Context) {
	var req pullImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "image name required")
		return
	}
	if !imagePattern.MatchString(req.Image) {
		fail(c, http.StatusBadRequest, "invalid image name")
		return
	}
	cfg, err := config.Load()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	binary := cfg.Engine
	if binary == "podman" || binary == "docker" {
		// already correct
	} else {
		binary = "docker"
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	cmd := exec.Command(binary, "pull", req.Image)
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(c.Writer, "data: {\"error\":\"%v\"}\n\n", err)
		c.Writer.Flush()
		return
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Fprintf(c.Writer, "data: %s\n\n", line)
		c.Writer.Flush()
	}
	if err := cmd.Wait(); err != nil {
		fmt.Fprintf(c.Writer, "data: {\"error\":\"%v\"}\n\n", err)
	} else {
		fmt.Fprintf(c.Writer, "data: {\"status\":\"done\"}\n\n")
	}
	c.Writer.Flush()
}

// DELETE /api/v1/images/:id
func RemoveImage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		fail(c, http.StatusBadRequest, "image id required")
		return
	}
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	if err := eng.RemoveImage(id); err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, nil)
}

// GET /api/v1/images/search?q=xxx
func SearchImages(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		fail(c, http.StatusBadRequest, "q parameter required")
		return
	}
	cfg, _ := config.Load()
	binary := cfg.Engine
	if binary != "docker" && binary != "podman" {
		fail(c, http.StatusBadRequest, "search only supported for docker/podman engine")
		return
	}
	out, err := exec.Command(binary, "search", "--format", "{{.Name}}\t{{.StarCount}}\t{{.Description}}", q).Output()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	type searchResult struct {
		Name        string `json:"name"`
		Stars       string `json:"stars"`
		Description string `json:"description"`
	}
	var results []searchResult
	for _, line := range splitLines(string(out)) {
		parts := splitTab(line, 3)
		if len(parts) < 1 {
			continue
		}
		r := searchResult{Name: parts[0]}
		if len(parts) > 1 {
			r.Stars = parts[1]
		}
		if len(parts) > 2 {
			r.Description = parts[2]
		}
		results = append(results, r)
	}
	ok(c, results)
}
