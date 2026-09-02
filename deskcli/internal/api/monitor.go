package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikachuim/deskcli/internal/config"
	"github.com/pikachuim/deskcli/internal/engine"
)

// GET /api/v1/monitor/stats
func GetStats(c *gin.Context) {
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	stats, err := eng.GetStats()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, stats)
}

// GET /api/v1/monitor/engines
func GetEngineStatus(c *gin.Context) {
	type engineStatus struct {
		Name      string `json:"name"`
		Available bool   `json:"available"`
		Version   string `json:"version"`
	}
	var results []engineStatus
	for _, name := range []string{"docker", "podman", "lxc", "lxd"} {
		eng, err := engine.New(name)
		status := engineStatus{Name: name}
		if err == nil {
			ver := eng.Version()
			status.Available = ver != "" && ver != "lxc (unavailable)" &&
				ver != "docker (unavailable)" && ver != "podman (unavailable)"
			status.Version = ver
		}
		results = append(results, status)
	}
	ok(c, results)
}

// GET /api/v1/status  (public)
func Status(c *gin.Context) {
	cfg, _ := config.Load()
	eng, err := engine.New(cfg.Engine)
	status := gin.H{
		"engine": cfg.Engine,
		"port":   cfg.Port,
	}
	if err == nil {
		stats, err2 := eng.GetStats()
		if err2 == nil {
			status["containers_running"] = stats.ContainersRunning
			status["containers_total"] = stats.ContainersTotal
			status["engine_version"] = stats.EngineVersion
		}
	}
	ok(c, status)
}

// GET /api/v1/config
func GetConfig(c *gin.Context) {
	cfg, err := config.Load()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, gin.H{
		"engine":     cfg.Engine,
		"port":       cfg.Port,
		"listen_addr": cfg.ListenAddr,
		"port_range": cfg.PortRange,
		"allow_exec": cfg.AllowExec,
		"static_dir": cfg.StaticDir,
	})
}

type updateConfigReq struct {
	Engine     string `json:"engine"`
	PortRange  string `json:"port_range"`
	AllowExec  *bool  `json:"allow_exec"`
	StaticDir  string `json:"static_dir"`
	JWTSecret  string `json:"jwt_secret"`
}

// PUT /api/v1/config  (admin only)
func UpdateConfig(c *gin.Context) {
	var req updateConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	cfg, err := config.Load()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if req.Engine != "" {
		cfg.Engine = req.Engine
	}
	if req.PortRange != "" {
		cfg.PortRange = req.PortRange
	}
	if req.AllowExec != nil {
		cfg.AllowExec = *req.AllowExec
	}
	if req.StaticDir != "" {
		cfg.StaticDir = req.StaticDir
	}
	if req.JWTSecret != "" {
		cfg.JWTSecret = req.JWTSecret
	}
	if err := config.Save(cfg); err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, gin.H{"message": "config updated"})
}
