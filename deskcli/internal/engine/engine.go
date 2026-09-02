package engine

import "fmt"

type ContainerInfo struct {
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Image  string   `json:"image"`
	Ports  []string `json:"ports"`
	ID     string   `json:"id"`
}

// ImageInfo describes a locally available image.
type ImageInfo struct {
	ID      string   `json:"id"`
	Tags    []string `json:"tags"`
	Size    string   `json:"size"`
	Created string   `json:"created"`
}

// EngineStats holds aggregate statistics from an engine.
type EngineStats struct {
	ContainersRunning int    `json:"containers_running"`
	ContainersTotal   int    `json:"containers_total"`
	ImagesTotal       int    `json:"images_total"`
	EngineVersion     string `json:"engine_version"`
}

type Engine interface {
	// ── existing methods ──────────────────────────────────────────────────
	List() ([]ContainerInfo, error)
	Start(name string) error
	Stop(name string) error
	Restart(name string) error
	Remove(name string) error
	Exec(name string, cmd []string, detach bool) error
	Run(image, name string, ports []string, env []string, devices []string) error
	Pull(image string) error
	SetPassword(name, password string) error
	Info(name string) (*ContainerInfo, error)
	GetIP(name string) (string, error)

	// ── new methods ───────────────────────────────────────────────────────
	ListImages() ([]ImageInfo, error)
	RemoveImage(id string) error
	GetStats() (*EngineStats, error)
	Version() string
}

func New(engineType string) (Engine, error) {
	switch engineType {
	case "docker":
		return &DockerEngine{binary: "docker"}, nil
	case "podman":
		return &DockerEngine{binary: "podman"}, nil
	case "lxc":
		return &LxcEngine{isLxd: false}, nil
	case "lxd":
		return &LxcEngine{isLxd: true}, nil
	default:
		return nil, fmt.Errorf("unsupported engine: %s", engineType)
	}
}
