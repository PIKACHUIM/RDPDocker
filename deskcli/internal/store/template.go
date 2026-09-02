package store

// TemplateRecord represents a row in the templates table.
type TemplateRecord struct {
	ID          int64  `json:"id"`
	OSType      string `json:"os_type"`
	OSVersion   string `json:"os_version"`
	DEName      string `json:"de_name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	ImageFull   string `json:"image_full"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
}

// SeedTemplates inserts preset templates if the table is empty.
func (d *DB) SeedTemplates() error {
	var count int
	_ = d.db.QueryRow(`SELECT COUNT(*) FROM templates`).Scan(&count)
	if count > 0 {
		return nil
	}

	templates := []TemplateRecord{
		// ── Debian 12 ─────────────────────────────────────────────────────────
		{OSType: "debian", OSVersion: "12", DEName: "server", DisplayName: "Debian 12 · Server",
			Description: "Minimal Debian 12 headless server, ideal for CLI-only workloads.",
			ImageFull:   "pikachuim/debian:12-server", Icon: "server", Category: "server"},
		{OSType: "debian", OSVersion: "12", DEName: "x11gui", DisplayName: "Debian 12 · X11 Core",
			Description: "Lightweight X11 base layer, no full desktop environment.",
			ImageFull:   "pikachuim/debian:12-x11gui", Icon: "monitor", Category: "minimal"},
		{OSType: "debian", OSVersion: "12", DEName: "gnome3", DisplayName: "Debian 12 · GNOME 3",
			Description: "Debian 12 with the classic GNOME 3 desktop environment.",
			ImageFull:   "pikachuim/debian:12-gnome3", Icon: "layout-grid", Category: "desktop"},
		{OSType: "debian", OSVersion: "12", DEName: "xfce4l", DisplayName: "Debian 12 · Xfce 4",
			Description: "Debian 12 with fast, lightweight Xfce 4 desktop.",
			ImageFull:   "pikachuim/debian:12-xfce4l", Icon: "zap", Category: "desktop"},
		{OSType: "debian", OSVersion: "12", DEName: "plasma", DisplayName: "Debian 12 · KDE Plasma",
			Description: "Debian 12 with feature-rich KDE Plasma desktop.",
			ImageFull:   "pikachuim/debian:12-plasma", Icon: "globe", Category: "desktop"},
		{OSType: "debian", OSVersion: "12", DEName: "deepin", DisplayName: "Debian 12 · Deepin DE",
			Description: "Debian 12 with the elegant Deepin Desktop Environment.",
			ImageFull:   "pikachuim/debian:12-deepin", Icon: "sparkles", Category: "desktop"},
		{OSType: "debian", OSVersion: "12", DEName: "lingmo", DisplayName: "Debian 12 · Lingmo",
			Description: "Debian 12 with the Lingmo desktop environment.",
			ImageFull:   "pikachuim/debian:12-lingmo", Icon: "wind", Category: "desktop"},
		// ── Ubuntu 22.04 ──────────────────────────────────────────────────────
		{OSType: "ubuntu", OSVersion: "22.04", DEName: "server", DisplayName: "Ubuntu 22.04 · Server",
			Description: "Ubuntu 22.04 LTS minimal headless server.",
			ImageFull:   "pikachuim/ubuntu:22.04-server", Icon: "server", Category: "server"},
		{OSType: "ubuntu", OSVersion: "22.04", DEName: "x11gui", DisplayName: "Ubuntu 22.04 · X11 Core",
			Description: "Ubuntu 22.04 lightweight X11 base layer.",
			ImageFull:   "pikachuim/ubuntu:22.04-x11gui", Icon: "monitor", Category: "minimal"},
		{OSType: "ubuntu", OSVersion: "22.04", DEName: "gnome3", DisplayName: "Ubuntu 22.04 · GNOME 3",
			Description: "Ubuntu 22.04 LTS with GNOME 3 desktop.",
			ImageFull:   "pikachuim/ubuntu:22.04-gnome3", Icon: "layout-grid", Category: "desktop"},
		{OSType: "ubuntu", OSVersion: "22.04", DEName: "xfce4l", DisplayName: "Ubuntu 22.04 · Xfce 4",
			Description: "Ubuntu 22.04 with lightweight Xfce 4 desktop.",
			ImageFull:   "pikachuim/ubuntu:22.04-xfce4l", Icon: "zap", Category: "desktop"},
		{OSType: "ubuntu", OSVersion: "22.04", DEName: "plasma", DisplayName: "Ubuntu 22.04 · KDE Plasma",
			Description: "Ubuntu 22.04 with KDE Plasma desktop.",
			ImageFull:   "pikachuim/ubuntu:22.04-plasma", Icon: "globe", Category: "desktop"},
		{OSType: "ubuntu", OSVersion: "22.04", DEName: "deepin", DisplayName: "Ubuntu 22.04 · Deepin DE",
			Description: "Ubuntu 22.04 with Deepin Desktop Environment.",
			ImageFull:   "pikachuim/ubuntu:22.04-deepin", Icon: "sparkles", Category: "desktop"},
		// ── Ubuntu 24.04 ──────────────────────────────────────────────────────
		{OSType: "ubuntu", OSVersion: "24.04", DEName: "server", DisplayName: "Ubuntu 24.04 · Server",
			Description: "Ubuntu 24.04 LTS minimal headless server.",
			ImageFull:   "pikachuim/ubuntu:24.04-server", Icon: "server", Category: "server"},
		{OSType: "ubuntu", OSVersion: "24.04", DEName: "x11gui", DisplayName: "Ubuntu 24.04 · X11 Core",
			Description: "Ubuntu 24.04 lightweight X11 base layer.",
			ImageFull:   "pikachuim/ubuntu:24.04-x11gui", Icon: "monitor", Category: "minimal"},
		{OSType: "ubuntu", OSVersion: "24.04", DEName: "gnome3", DisplayName: "Ubuntu 24.04 · GNOME 3",
			Description: "Ubuntu 24.04 LTS with GNOME 3 desktop.",
			ImageFull:   "pikachuim/ubuntu:24.04-gnome3", Icon: "layout-grid", Category: "desktop"},
		{OSType: "ubuntu", OSVersion: "24.04", DEName: "plasma", DisplayName: "Ubuntu 24.04 · KDE Plasma",
			Description: "Ubuntu 24.04 with KDE Plasma 6 desktop.",
			ImageFull:   "pikachuim/ubuntu:24.04-plasma", Icon: "globe", Category: "desktop"},
		// ── Alpine 3.19 ───────────────────────────────────────────────────────
		{OSType: "alpine", OSVersion: "3.19", DEName: "server", DisplayName: "Alpine 3.19 · Server",
			Description: "Ultra-minimal Alpine 3.19 server (< 10 MB base).",
			ImageFull:   "pikachuim/alpine:3.19-server", Icon: "server", Category: "server"},
		{OSType: "alpine", OSVersion: "3.19", DEName: "x11gui", DisplayName: "Alpine 3.19 · X11 Core",
			Description: "Alpine 3.19 with minimal X11 layer.",
			ImageFull:   "pikachuim/alpine:3.19-x11gui", Icon: "monitor", Category: "minimal"},
		{OSType: "alpine", OSVersion: "3.19", DEName: "xfce4l", DisplayName: "Alpine 3.19 · Xfce 4",
			Description: "Alpine 3.19 with Xfce 4 for a tiny footprint GUI.",
			ImageFull:   "pikachuim/alpine:3.19-xfce4l", Icon: "zap", Category: "desktop"},
		// ── Fedora 40 ─────────────────────────────────────────────────────────
		{OSType: "fedora", OSVersion: "40", DEName: "server", DisplayName: "Fedora 40 · Server",
			Description: "Fedora 40 headless server with latest packages.",
			ImageFull:   "pikachuim/fedora:40-server", Icon: "server", Category: "server"},
		{OSType: "fedora", OSVersion: "40", DEName: "x11gui", DisplayName: "Fedora 40 · X11 Core",
			Description: "Fedora 40 lightweight X11 base layer.",
			ImageFull:   "pikachuim/fedora:40-x11gui", Icon: "monitor", Category: "minimal"},
		{OSType: "fedora", OSVersion: "40", DEName: "gnome3", DisplayName: "Fedora 40 · GNOME 3",
			Description: "Fedora 40 with GNOME 3 (Fedora's flagship desktop).",
			ImageFull:   "pikachuim/fedora:40-gnome3", Icon: "layout-grid", Category: "desktop"},
		{OSType: "fedora", OSVersion: "40", DEName: "plasma", DisplayName: "Fedora 40 · KDE Plasma",
			Description: "Fedora 40 with KDE Plasma desktop.",
			ImageFull:   "pikachuim/fedora:40-plasma", Icon: "globe", Category: "desktop"},
		// ── Arch Linux ────────────────────────────────────────────────────────
		{OSType: "arch", OSVersion: "latest", DEName: "server", DisplayName: "Arch Linux · Server",
			Description: "Rolling-release Arch Linux minimal server.",
			ImageFull:   "pikachuim/arch:latest-server", Icon: "server", Category: "server"},
		{OSType: "arch", OSVersion: "latest", DEName: "x11gui", DisplayName: "Arch Linux · X11 Core",
			Description: "Arch Linux with minimal X11 layer.",
			ImageFull:   "pikachuim/arch:latest-x11gui", Icon: "monitor", Category: "minimal"},
		{OSType: "arch", OSVersion: "latest", DEName: "gnome3", DisplayName: "Arch Linux · GNOME 3",
			Description: "Arch Linux with GNOME 3 — bleeding edge.",
			ImageFull:   "pikachuim/arch:latest-gnome3", Icon: "layout-grid", Category: "desktop"},
		{OSType: "arch", OSVersion: "latest", DEName: "plasma", DisplayName: "Arch Linux · KDE Plasma",
			Description: "Arch Linux with KDE Plasma 6.",
			ImageFull:   "pikachuim/arch:latest-plasma", Icon: "globe", Category: "desktop"},
		{OSType: "arch", OSVersion: "latest", DEName: "hyprland", DisplayName: "Arch Linux · Hyprland",
			Description: "Arch Linux with the modern dynamic tiling Wayland compositor Hyprland.",
			ImageFull:   "pikachuim/arch:latest-hyprland", Icon: "layout-panel-left", Category: "wayland"},
		{OSType: "arch", OSVersion: "latest", DEName: "niri", DisplayName: "Arch Linux · Niri",
			Description: "Arch Linux with Niri — scrollable-tiling Wayland compositor.",
			ImageFull:   "pikachuim/arch:latest-niri", Icon: "columns", Category: "wayland"},
	}

	stmt, err := d.db.Prepare(
		`INSERT INTO templates
		 (os_type,os_version,de_name,display_name,description,image_full,icon,category,is_active)
		 VALUES (?,?,?,?,?,?,?,?,1)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, t := range templates {
		if _, err := stmt.Exec(t.OSType, t.OSVersion, t.DEName, t.DisplayName,
			t.Description, t.ImageFull, t.Icon, t.Category); err != nil {
			return err
		}
	}
	return nil
}

// ─── Template CRUD ────────────────────────────────────────────────────────────

func (d *DB) ListTemplates(osType, deName, category string) ([]TemplateRecord, error) {
	query := `SELECT id,os_type,os_version,de_name,display_name,description,
	                 image_full,icon,category,is_active,created_at
	          FROM templates WHERE is_active=1`
	args := []interface{}{}
	if osType != "" {
		query += " AND os_type=?"
		args = append(args, osType)
	}
	if deName != "" {
		query += " AND de_name=?"
		args = append(args, deName)
	}
	if category != "" {
		query += " AND category=?"
		args = append(args, category)
	}
	query += " ORDER BY os_type, os_version, de_name"

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []TemplateRecord
	for rows.Next() {
		var t TemplateRecord
		var active int
		if err := rows.Scan(&t.ID, &t.OSType, &t.OSVersion, &t.DEName,
			&t.DisplayName, &t.Description, &t.ImageFull,
			&t.Icon, &t.Category, &active, &t.CreatedAt); err == nil {
			t.IsActive = active == 1
			out = append(out, t)
		}
	}
	return out, nil
}

func (d *DB) GetTemplate(id int64) (*TemplateRecord, error) {
	t := &TemplateRecord{}
	var active int
	err := d.db.QueryRow(
		`SELECT id,os_type,os_version,de_name,display_name,description,
		        image_full,icon,category,is_active,created_at
		 FROM templates WHERE id=?`, id,
	).Scan(&t.ID, &t.OSType, &t.OSVersion, &t.DEName,
		&t.DisplayName, &t.Description, &t.ImageFull,
		&t.Icon, &t.Category, &active, &t.CreatedAt)
	if err != nil {
		return nil, ErrNotFound
	}
	t.IsActive = active == 1
	return t, nil
}

func (d *DB) DistinctOSTypes() []string {
	rows, _ := d.db.Query(`SELECT DISTINCT os_type FROM templates WHERE is_active=1 ORDER BY os_type`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if rows.Scan(&v) == nil {
			out = append(out, v)
		}
	}
	return out
}

func (d *DB) DistinctDENames() []string {
	rows, _ := d.db.Query(`SELECT DISTINCT de_name FROM templates WHERE is_active=1 ORDER BY de_name`)
	if rows == nil {
		return nil
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var v string
		if rows.Scan(&v) == nil {
			out = append(out, v)
		}
	}
	return out
}
