package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pikachuim/deskcli/internal/config"
	"github.com/pikachuim/deskcli/internal/engine"
	"github.com/pikachuim/deskcli/internal/forward"
	"github.com/pikachuim/deskcli/internal/store"
)

func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data, "message": "ok"})
}

func fail(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"success": false, "data": nil, "message": msg})
}

var (
	containerNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,62}$`)
	imagePattern         = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_./:@-]{0,254}$`)
	packagePattern       = regexp.MustCompile(`^[a-z0-9][a-z0-9+._-]{0,127}$`)
)

func containerName(c *gin.Context) (string, bool) {
	name := c.Param("name")
	if !containerNamePattern.MatchString(name) {
		fail(c, http.StatusBadRequest, "invalid container name")
		return "", false
	}
	return name, true
}

func parsePortMappings(values []string) ([]config.PortMap, error) {
	mappings := make([]config.PortMap, 0, len(values))
	for _, value := range values {
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid port mapping %q", value)
		}
		ext, err := strconv.Atoi(parts[0])
		if err != nil || ext < 1 || ext > 65535 {
			return nil, fmt.Errorf("invalid external port in %q", value)
		}
		internal, err := strconv.Atoi(parts[1])
		if err != nil || internal < 1 || internal > 65535 {
			return nil, fmt.Errorf("invalid internal port in %q", value)
		}
		mappings = append(mappings, config.PortMap{Ext: ext, Int: internal})
	}
	return mappings, nil
}

func loadEngine(c *gin.Context) (engine.Engine, *config.Config, bool) {
	cfg, err := config.Load()
	if err != nil {
		fail(c, 500, err.Error())
		return nil, nil, false
	}
	eng, err := newEngine(cfg)
	if err != nil {
		fail(c, 500, err.Error())
		return nil, nil, false
	}
	return eng, cfg, true
}

func ListContainers(c *gin.Context) {
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	list, err := eng.List()
	if err != nil {
		fail(c, 500, err.Error())
		return
	}
	uid, role := currentUser(c)
	if role == "admin" {
		ok(c, list)
		return
	}
	// Non-admin users only see containers they own.
	db := store.Get()
	owned := map[string]bool{}
	if db != nil {
		if recs, err := db.ListContainers(); err == nil {
			for _, r := range recs {
				if r.UserID == uid {
					owned[r.Name] = true
				}
			}
		}
	}
	filtered := make([]engine.ContainerInfo, 0, len(list))
	for _, ci := range list {
		if owned[ci.Name] {
			filtered = append(filtered, ci)
		}
	}
	ok(c, filtered)
}

type createReq struct {
	Image     string   `json:"image"`
	Name      string   `json:"name"`
	Ports     []string `json:"ports"`
	Softwares []string `json:"softwares"`
}

func CreateContainer(c *gin.Context) {
	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if !imagePattern.MatchString(req.Image) {
		fail(c, http.StatusBadRequest, "invalid image name")
		return
	}
	name := req.Name
	if name == "" {
		name = strings.ReplaceAll(strings.Split(req.Image, ":")[0], "/", "-")
	}
	if !containerNamePattern.MatchString(name) {
		fail(c, http.StatusBadRequest, "invalid container name")
		return
	}
	portMaps, err := parsePortMappings(req.Ports)
	if err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	for _, sw := range req.Softwares {
		if !packagePattern.MatchString(sw) {
			fail(c, http.StatusBadRequest, fmt.Sprintf("invalid package name %q", sw))
			return
		}
	}
	eng, cfg, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	if err := eng.Run(req.Image, name, req.Ports, nil, nil); err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	for _, sw := range req.Softwares {
		if err := eng.Exec(name, []string{"apt-get", "install", "-y", sw}, true); err != nil {
			fail(c, http.StatusInternalServerError, fmt.Sprintf("container created but package installation failed: %v", err))
			return
		}
	}
	if err := config.SaveContainer(&config.ContainerConfig{Name: name, Image: req.Image, Engine: cfg.Engine, Ports: portMaps}); err != nil {
		fail(c, http.StatusInternalServerError, fmt.Sprintf("container created but configuration could not be saved: %v", err))
		return
	}
	uid, _ := currentUser(c)
	saveContainerRecord(name, req.Image, cfg.Engine, portMaps, uid, 0)
	if ip, err := eng.GetIP(name); err == nil && ip != "" {
		if err := forward.ApplyContainerRules(ip, portMaps); err != nil {
			fail(c, http.StatusInternalServerError, fmt.Sprintf("container created but port forwarding failed: %v", err))
			return
		}
	}
	ok(c, gin.H{"name": name})
}

func GetContainer(c *gin.Context) {
	name, valid := containerName(c)
	if !valid {
		return
	}
	if !authorizeContainer(c, name) {
		return
	}
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	info, err := eng.Info(name)
	if err != nil {
		fail(c, 404, fmt.Sprintf("container %s not found: %v", name, err))
		return
	}
	ok(c, info)
}

func RemoveContainer(c *gin.Context) {
	name, valid := containerName(c)
	if !valid {
		return
	}
	if !authorizeContainer(c, name) {
		return
	}
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	if err := eng.Remove(name); err != nil {
		fail(c, 500, err.Error())
		return
	}
	_ = config.DeleteContainer(name)
	if db := store.Get(); db != nil {
		_ = db.Delete(name)
	}
	ok(c, nil)
}

func StartContainer(c *gin.Context) {
	name, valid := containerName(c)
	if !valid {
		return
	}
	if !authorizeContainer(c, name) {
		return
	}
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	if err := eng.Start(name); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

func StopContainer(c *gin.Context) {
	name, valid := containerName(c)
	if !valid {
		return
	}
	if !authorizeContainer(c, name) {
		return
	}
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	if err := eng.Stop(name); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

func RestartContainer(c *gin.Context) {
	name, valid := containerName(c)
	if !valid {
		return
	}
	if !authorizeContainer(c, name) {
		return
	}
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	if err := eng.Restart(name); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

type execReq struct {
	Cmd    []string `json:"cmd"`
	Detach bool     `json:"detach"`
}

func ExecContainer(c *gin.Context) {
	name, valid := containerName(c)
	if !valid {
		return
	}
	if !authorizeContainer(c, name) {
		return
	}
	var req execReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Cmd) == 0 || strings.TrimSpace(req.Cmd[0]) == "" {
		fail(c, http.StatusBadRequest, "cmd must contain an executable")
		return
	}
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	if err := eng.Exec(name, req.Cmd, req.Detach); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}

type passwdReq struct {
	Password string `json:"password"`
}

func SetPassword(c *gin.Context) {
	name, valid := containerName(c)
	if !valid {
		return
	}
	if !authorizeContainer(c, name) {
		return
	}
	var req passwdReq
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Password) < 12 || len(req.Password) > 256 {
		fail(c, http.StatusBadRequest, "password must be 12-256 characters")
		return
	}
	eng, _, ok2 := loadEngine(c)
	if !ok2 {
		return
	}
	if err := eng.SetPassword(name, req.Password); err != nil {
		fail(c, 500, err.Error())
		return
	}
	ok(c, nil)
}
