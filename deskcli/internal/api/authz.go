package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikachuim/deskcli/internal/config"
	"github.com/pikachuim/deskcli/internal/store"
)

// currentUser returns the authenticated user id and role from the context.
func currentUser(c *gin.Context) (int64, string) {
	uid, _ := c.Get("user_id")
	role, _ := c.Get("role")
	var id int64
	if uid != nil {
		id, _ = uid.(int64)
	}
	var r string
	if role != nil {
		r, _ = role.(string)
	}
	return id, r
}

// authorizeContainer reports whether the current user may operate on the named
// container. Admins bypass ownership checks; regular users must own it.
// On denial it writes an HTTP error and returns false.
func authorizeContainer(c *gin.Context, name string) bool {
	_, role := currentUser(c)
	if role == "admin" {
		return true
	}
	uid, _ := currentUser(c)
	db := store.Get()
	if db == nil {
		fail(c, http.StatusInternalServerError, "store not initialized")
		return false
	}
	rec, err := db.Get(name)
	if err != nil {
		// Untracked (legacy) container: only admins may operate it.
		fail(c, http.StatusForbidden, "forbidden")
		return false
	}
	if rec.UserID != uid {
		fail(c, http.StatusForbidden, "forbidden")
		return false
	}
	return true
}

// saveContainerRecord persists ownership metadata for a container so later
// operations can enforce per-user isolation.
func saveContainerRecord(name, image, engineType string, ports []config.PortMap, uid, templateID int64) {
	db := store.Get()
	if db == nil {
		return
	}
	ssh, rdp, nx, vnc := portsToFields(ports)
	_ = db.Save(&store.ContainerRecord{
		Name:       name,
		Image:      image,
		Engine:     engineType,
		PortSSH:    ssh,
		PortRDP:    rdp,
		PortNX:     nx,
		PortVNC:    vnc,
		UserID:     uid,
		TemplateID: templateID,
	})
}

// portsToFields extracts the external ports for the well-known internal
// services (SSH=22, RDP=3389, NX=4000, VNC=5900).
func portsToFields(ports []config.PortMap) (ssh, rdp, nx, vnc int) {
	for _, p := range ports {
		switch p.Int {
		case 22:
			ssh = p.Ext
		case 3389:
			rdp = p.Ext
		case 4000:
			nx = p.Ext
		case 5900:
			vnc = p.Ext
		}
	}
	return
}
