package store

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type ContainerRecord struct {
	Name       string
	Image      string
	Engine     string
	Password   string
	PortSSH    int
	PortRDP    int
	PortNX     int
	PortVNC    int
	UserID     int64
	TemplateID int64
}

type DB struct{ db *sql.DB }

func Open(dataDir string) (*DB, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", filepath.Join(dataDir, "deskapi.db")+"?_journal=WAL&_timeout=5000")
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		// containers table (original)
		`CREATE TABLE IF NOT EXISTS containers (
			name TEXT PRIMARY KEY, image TEXT, engine TEXT, password TEXT,
			port_ssh INTEGER, port_rdp INTEGER, port_nx INTEGER, port_vnc INTEGER
		)`,
		// users table
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT DEFAULT 'user',
			enabled INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		// templates table
		`CREATE TABLE IF NOT EXISTS templates (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			os_type TEXT NOT NULL,
			os_version TEXT NOT NULL,
			de_name TEXT NOT NULL,
			display_name TEXT NOT NULL,
			description TEXT DEFAULT '',
			image_full TEXT NOT NULL,
			icon TEXT DEFAULT 'monitor',
			category TEXT DEFAULT '',
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		// sessions table
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	// Add new columns to containers table if they don't exist yet (idempotent)
	for _, col := range []struct{ name, def string }{
		{"user_id", "INTEGER DEFAULT 0"},
		{"template_id", "INTEGER DEFAULT 0"},
	} {
		var cnt int
		_ = db.QueryRow(
			`SELECT COUNT(*) FROM pragma_table_info('containers') WHERE name=?`, col.name,
		).Scan(&cnt)
		if cnt == 0 {
			_, _ = db.Exec("ALTER TABLE containers ADD COLUMN " + col.name + " " + col.def)
		}
	}
	return nil
}

// ─── Container methods ────────────────────────────────────────────────────────

func (d *DB) Save(r *ContainerRecord) error {
	_, err := d.db.Exec(
		`INSERT OR REPLACE INTO containers
		 (name,image,engine,password,port_ssh,port_rdp,port_nx,port_vnc,user_id,template_id)
		 VALUES (?,?,?,?,?,?,?,?,?,?)`,
		r.Name, r.Image, r.Engine, r.Password,
		r.PortSSH, r.PortRDP, r.PortNX, r.PortVNC,
		r.UserID, r.TemplateID)
	return err
}

func (d *DB) Get(name string) (*ContainerRecord, error) {
	r := &ContainerRecord{}
	err := d.db.QueryRow(
		`SELECT name,image,engine,password,port_ssh,port_rdp,port_nx,port_vnc,
		        COALESCE(user_id,0),COALESCE(template_id,0)
		 FROM containers WHERE name=?`, name,
	).Scan(&r.Name, &r.Image, &r.Engine, &r.Password,
		&r.PortSSH, &r.PortRDP, &r.PortNX, &r.PortVNC,
		&r.UserID, &r.TemplateID)
	return r, err
}

func (d *DB) UpdatePassword(name, password string) error {
	_, err := d.db.Exec(`UPDATE containers SET password=? WHERE name=?`, password, name)
	return err
}

func (d *DB) Delete(name string) error {
	_, err := d.db.Exec(`DELETE FROM containers WHERE name=?`, name)
	return err
}

func (d *DB) ListContainers() ([]ContainerRecord, error) {
	rows, err := d.db.Query(
		`SELECT name,image,engine,password,port_ssh,port_rdp,port_nx,port_vnc,
		        COALESCE(user_id,0),COALESCE(template_id,0)
		 FROM containers ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ContainerRecord
	for rows.Next() {
		var r ContainerRecord
		if err := rows.Scan(&r.Name, &r.Image, &r.Engine, &r.Password,
			&r.PortSSH, &r.PortRDP, &r.PortNX, &r.PortVNC,
			&r.UserID, &r.TemplateID); err == nil {
			out = append(out, r)
		}
	}
	return out, nil
}

func (d *DB) UsedPorts() map[int]bool {
	used := make(map[int]bool)
	rows, err := d.db.Query(`SELECT port_ssh,port_rdp,port_nx,port_vnc FROM containers`)
	if err != nil {
		return used
	}
	defer rows.Close()
	for rows.Next() {
		var a, b, c, e int
		if rows.Scan(&a, &b, &c, &e) == nil {
			used[a], used[b], used[c], used[e] = true, true, true, true
		}
	}
	return used
}

func (d *DB) Close() error { return d.db.Close() }
