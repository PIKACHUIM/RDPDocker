package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pikachuim/deskcli/internal/config"
	"github.com/pikachuim/deskcli/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// POST /api/v1/auth/login
func Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "username and password required")
		return
	}
	db := store.Get()
	if db == nil {
		fail(c, http.StatusInternalServerError, "store not initialized")
		return
	}
	user, err := db.GetUserByUsername(req.Username)
	if err != nil {
		fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if !user.Enabled {
		fail(c, http.StatusForbidden, "account disabled")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		fail(c, http.StatusUnauthorized, "invalid credentials")
		return
	}
	cfg, _ := config.Load()
	token, err := signJWT(cfg, user.ID, user.Username, user.Role)
	if err != nil {
		fail(c, http.StatusInternalServerError, "token generation failed")
		return
	}
	ok(c, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username, "role": user.Role},
	})
}

type registerReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8"`
}

// POST /api/v1/auth/register  (only allowed when no users exist, i.e. first admin)
func Register(c *gin.Context) {
	db := store.Get()
	if db == nil {
		fail(c, http.StatusInternalServerError, "store not initialized")
		return
	}
	if db.CountUsers() > 0 {
		fail(c, http.StatusForbidden, "registration is disabled; contact your administrator")
		return
	}
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		fail(c, http.StatusInternalServerError, "password hash failed")
		return
	}
	user, err := db.CreateUser(req.Username, string(hash), "admin")
	if err != nil {
		fail(c, http.StatusConflict, "username already taken")
		return
	}
	cfg, _ := config.Load()
	token, _ := signJWT(cfg, user.ID, user.Username, user.Role)
	ok(c, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username, "role": user.Role},
	})
}

// GET /api/v1/auth/me
func GetMe(c *gin.Context) {
	uid, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	ok(c, gin.H{"id": uid, "username": username, "role": role})
}

// POST /api/v1/auth/refresh
func RefreshToken(c *gin.Context) {
	uid, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")
	cfg, _ := config.Load()
	token, err := signJWT(cfg, uid.(int64), username.(string), role.(string))
	if err != nil {
		fail(c, http.StatusInternalServerError, "token generation failed")
		return
	}
	ok(c, gin.H{"token": token})
}

type changePassReq struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// POST /api/v1/auth/change-password
func ChangePassword(c *gin.Context) {
	var req changePassReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	uid, _ := c.Get("user_id")
	db := store.Get()
	user, err := db.GetUserByID(uid.(int64))
	if err != nil {
		fail(c, http.StatusNotFound, "user not found")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		fail(c, http.StatusUnauthorized, "wrong current password")
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err := db.UpdateUserPassword(uid.(int64), string(hash)); err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, gin.H{"message": "password updated"})
}

// ─── User management (admin) ──────────────────────────────────────────────────

// GET /api/v1/users
func ListUsers(c *gin.Context) {
	db := store.Get()
	users, err := db.ListUsers()
	if err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, users)
}

type createUserReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role"`
}

// POST /api/v1/users
func CreateUser(c *gin.Context) {
	var req createUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	role := req.Role
	if role != "admin" && role != "user" {
		role = "user"
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	db := store.Get()
	user, err := db.CreateUser(req.Username, string(hash), role)
	if err != nil {
		fail(c, http.StatusConflict, "username already taken")
		return
	}
	ok(c, user)
}

// PUT /api/v1/users/:id/enabled
func SetUserEnabled(c *gin.Context) {
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		fail(c, http.StatusBadRequest, err.Error())
		return
	}
	idStr := c.Param("id")
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		fail(c, http.StatusBadRequest, "invalid user id")
		return
	}
	db := store.Get()
	if err := db.SetUserEnabled(id, body.Enabled); err != nil {
		fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	ok(c, nil)
}
