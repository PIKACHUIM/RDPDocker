package store

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrNotFound = errors.New("record not found")

// DefaultAdminUsername / DefaultAdminPassword are the credentials created on a
// fresh install so the web UI is reachable before any user exists.
const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin123"
)

// EnsureDefaultAdmin creates the bootstrap admin account when the users table is
// empty. It reports whether a new account was created so callers can warn the
// operator to change the password.
func (d *DB) EnsureDefaultAdmin() (bool, error) {
	if d.CountUsers() > 0 {
		return false, nil
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(DefaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	if _, err := d.CreateUser(DefaultAdminUsername, string(hash), "admin"); err != nil {
		return false, err
	}
	return true, nil
}

// UserRecord represents a row in the users table.
type UserRecord struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	Role         string `json:"role"`
	Enabled      bool   `json:"enabled"`
	CreatedAt    string `json:"created_at"`
}

func (d *DB) CreateUser(username, hash, role string) (*UserRecord, error) {
	res, err := d.db.Exec(
		`INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)`,
		username, hash, role,
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return d.GetUserByID(id)
}

func (d *DB) GetUserByID(id int64) (*UserRecord, error) {
	u := &UserRecord{}
	var enabled int
	err := d.db.QueryRow(
		`SELECT id, username, password_hash, role, enabled, created_at FROM users WHERE id=?`, id,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &enabled, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u.Enabled = enabled == 1
	return u, nil
}

func (d *DB) GetUserByUsername(username string) (*UserRecord, error) {
	u := &UserRecord{}
	var enabled int
	err := d.db.QueryRow(
		`SELECT id, username, password_hash, role, enabled, created_at FROM users WHERE username=?`, username,
	).Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &enabled, &u.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u.Enabled = enabled == 1
	return u, nil
}

func (d *DB) CountUsers() int {
	var count int
	_ = d.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count
}

func (d *DB) ListUsers() ([]UserRecord, error) {
	rows, err := d.db.Query(
		`SELECT id, username, password_hash, role, enabled, created_at FROM users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []UserRecord
	for rows.Next() {
		var u UserRecord
		var enabled int
		if err := rows.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &enabled, &u.CreatedAt); err == nil {
			u.Enabled = enabled == 1
			users = append(users, u)
		}
	}
	return users, nil
}

func (d *DB) UpdateUserPassword(id int64, hash string) error {
	_, err := d.db.Exec(`UPDATE users SET password_hash=? WHERE id=?`, hash, id)
	return err
}

func (d *DB) SetUserEnabled(id int64, enabled bool) error {
	v := 0
	if enabled {
		v = 1
	}
	_, err := d.db.Exec(`UPDATE users SET enabled=? WHERE id=?`, v, id)
	return err
}

func (d *DB) SetUserRole(id int64, role string) error {
	_, err := d.db.Exec(`UPDATE users SET role=? WHERE id=?`, role, id)
	return err
}
