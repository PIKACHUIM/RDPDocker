package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pikachuim/deskcli/internal/config"
	"github.com/pikachuim/deskcli/internal/engine"
	"github.com/pikachuim/deskcli/internal/forward"
	"github.com/pikachuim/deskcli/internal/store"
)

func Start(cfg *config.Config) error {
	// Initialize store (DB migrations + seed data)
	if err := store.Init(cfg.DataDir); err != nil {
		return fmt.Errorf("store init: %w", err)
	}

	go restoreIPTables(cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// CORS: allow all origins in dev; tighten in production via config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	v1 := r.Group("/api/v1")

	// ── Public endpoints ──────────────────────────────────────────────────
	v1.POST("/auth/login", Login)
	v1.POST("/auth/register", Register)
	v1.GET("/status", Status)

	// ── Authenticated endpoints ───────────────────────────────────────────
	auth := v1.Group("", jwtAuthMiddleware(cfg))
	{
		// Auth
		auth.GET("/auth/me", GetMe)
		auth.POST("/auth/refresh", RefreshToken)
		auth.POST("/auth/change-password", ChangePassword)

		// Users (admin only)
		adminOnly := auth.Group("", requireAdmin)
		{
			adminOnly.GET("/users", ListUsers)
			adminOnly.POST("/users", CreateUser)
		}

		// Container management (existing endpoints, paths unchanged)
		auth.GET("/containers", ListContainers)
		auth.POST("/containers", CreateContainer)
		auth.GET("/containers/:name", GetContainer)
		auth.DELETE("/containers/:name", RemoveContainer)
		auth.POST("/containers/:name/start", StartContainer)
		auth.POST("/containers/:name/stop", StopContainer)
		auth.POST("/containers/:name/restart", RestartContainer)
		if cfg.AllowExec {
			auth.POST("/containers/:name/exec", ExecContainer)
		}
		auth.POST("/containers/:name/passwd", SetPassword)

		// Image management
		auth.GET("/images", ListImages)
		auth.POST("/images/pull", PullImage)
		auth.DELETE("/images/:id", RemoveImage)
		auth.GET("/images/search", SearchImages)

		// Template management
		auth.GET("/templates", ListTemplates)
		auth.GET("/templates/os-types", GetOSTypes)
		auth.GET("/templates/de-names", GetDENames)
		auth.POST("/templates/:id/deploy", DeployTemplate)

		// Monitoring
		auth.GET("/monitor/stats", GetStats)
		auth.GET("/monitor/engines", GetEngineStatus)

		// Config (admin only)
		adminOnly.GET("/config", GetConfig)
		adminOnly.PUT("/config", UpdateConfig)
	}

	// WebSocket terminal (auth via query param)
	r.GET("/api/v1/ws/terminal/:name", wsAuthMiddleware(cfg), TerminalWS)

	// Static frontend files
	if cfg.StaticDir != "" {
		if _, err := os.Stat(cfg.StaticDir); err == nil {
			r.Static("/assets", cfg.StaticDir+"/assets")
			r.StaticFile("/favicon.ico", cfg.StaticDir+"/favicon.ico")
			// SPA fallback: all non-API routes serve index.html
			r.NoRoute(func(c *gin.Context) {
				path := c.Request.URL.Path
				if len(path) >= 4 && path[:4] == "/api" {
					c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "not found"})
					return
				}
				c.File(cfg.StaticDir + "/index.html")
			})
		}
	}

	listenAddr := cfg.ListenAddr
	if listenAddr == "" {
		listenAddr = "127.0.0.1"
	}
	return r.Run(fmt.Sprintf("%s:%d", listenAddr, cfg.Port))
}

func newEngine(cfg *config.Config) (engine.Engine, error) {
	return engine.New(cfg.Engine)
}

func newEngine2(engineType string) (engine.Engine, error) {
	return engine.New(engineType)
}

func restoreIPTables(cfg *config.Config) {
	containers, err := config.ListContainerConfigs()
	if err != nil || len(containers) == 0 {
		return
	}
	eng, err := engine.New(cfg.Engine)
	if err != nil {
		return
	}
	for _, cc := range containers {
		if len(cc.Ports) == 0 {
			continue
		}
		if ip, err := eng.GetIP(cc.Name); err == nil && ip != "" {
			_ = forward.ApplyContainerRules(ip, cc.Ports)
		}
	}
}
