package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pikachuim/deskcli/internal/config"
	"github.com/pikachuim/deskcli/internal/store"
)

// GET /api/v1/templates
func ListTemplates(c *gin.Context) {
	db := store.Get()
	if db == nil {
		fail(c, http.StatusInternalServerError, "store not initialized")
		return
	}
	osType := c.Query("os_type")
	deName := c.Query("de_name")
	category := c.Query("category")
	templates, err := db.ListTemplates(osType, deName, category)
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	if templates == nil {
		templates = []store.TemplateRecord{}
	}
	ok(c, templates)
}

// GET /api/v1/templates/os-types
func GetOSTypes(c *gin.Context) {
	db := store.Get()
	if db == nil {
		fail(c, http.StatusInternalServerError, "store not initialized")
		return
	}
	ok(c, db.DistinctOSTypes())
}

// GET /api/v1/templates/de-names
func GetDENames(c *gin.Context) {
	db := store.Get()
	if db == nil {
		fail(c, http.StatusInternalServerError, "store not initialized")
		return
	}
	ok(c, db.DistinctDENames())
}

type deployReq struct {
	Name      string   `json:"name"`
	Engine    string   `json:"engine"`
	Ports     []string `json:"ports"`
	Softwares []string `json:"softwares"`
}

// POST /api/v1/templates/:id/deploy
func DeployTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fail(c, http.StatusBadRequest, "invalid template id")
		return
	}
	db := store.Get()
	if db == nil {
		fail(c, http.StatusInternalServerError, "store not initialized")
		return
	}
	tmpl, err := db.GetTemplate(id)
	if err != nil {
		fail(c, http.StatusNotFound, "template not found")
		return
	}
	var req deployReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	// Determine container name
	name := req.Name
	if name == "" {
		name = fmt.Sprintf("%s-%s", tmpl.OSType, tmpl.DEName)
	}
	if !containerNamePattern.MatchString(name) {
		fail(c, http.StatusBadRequest, "invalid container name")
		return
	}
	// Validate softwares
	for _, sw := range req.Softwares {
		if !packagePattern.MatchString(sw) {
			fail(c, http.StatusBadRequest, fmt.Sprintf("invalid package name %q", sw))
			return
		}
	}

	// Load engine - prefer template-specified engine, fallback to config
	cfg, err := config.Load()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	engineType := cfg.Engine
	if req.Engine != "" {
		engineType = req.Engine
	}

	// If ports not provided, auto-allocate from used ports DB
	ports := req.Ports
	if len(ports) == 0 {
		usedPorts := db.UsedPorts()
		allocated, err := config.AllocatePorts(cfg.PortRange, usedPorts)
		if err != nil {
			fail(c, http.StatusInternalServerError, "no available ports: "+err.Error())
			return
		}
		ports = []string{
			fmt.Sprintf("%d:22", allocated[0]),
			fmt.Sprintf("%d:3389", allocated[1]),
			fmt.Sprintf("%d:4000", allocated[2]),
			fmt.Sprintf("%d:5900", allocated[3]),
		}
	}
	// Validate port mappings
	portMaps, err := parsePortMappings(ports)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}

	eng, err := newEngine2(engineType)
	if err != nil {
		fail(c, http.StatusBadRequest, "unsupported engine: "+engineType)
		return
	}
	if err := eng.Run(tmpl.ImageFull, name, ports, nil, nil); err != nil {
		fail(c, http.StatusInternalServerError, "container start failed: "+err.Error())
		return
	}
	for _, sw := range req.Softwares {
		_ = eng.Exec(name, []string{"apt-get", "install", "-y", sw}, true)
	}

	// Save to config
	_ = config.SaveContainer(&config.ContainerConfig{
		Name: name, Image: tmpl.ImageFull, Engine: engineType, Ports: portMaps,
	})

	// Persist ownership metadata for multi-tenant isolation.
	uid, _ := currentUser(c)
	saveContainerRecord(name, tmpl.ImageFull, engineType, portMaps, uid, tmpl.ID)

	ok(c, gin.H{
		"name":   name,
		"image":  tmpl.ImageFull,
		"ports":  ports,
		"engine": engineType,
	})
}
