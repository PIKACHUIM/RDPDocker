package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	ConfigDir  = "/opt/deskcli/conf"
	DataDir    = "/opt/deskcli/data"
	MainConfig = "/opt/deskcli/conf/config.yaml"
)

type Config struct {
	Token      string `yaml:"token"`
	JWTSecret  string `yaml:"jwt_secret"`
	Engine     string `yaml:"engine"`
	ListenAddr string `yaml:"listen_addr"`
	Port       int    `yaml:"port"`
	DataDir    string `yaml:"data_dir"`
	PortRange  string `yaml:"port_range"`
	AllowExec  bool   `yaml:"allow_exec"`
	StaticDir  string `yaml:"static_dir"`
}

type PortMap struct {
	Ext int `yaml:"ext"`
	Int int `yaml:"int"`
}

type ContainerConfig struct {
	Name   string    `yaml:"name"`
	Image  string    `yaml:"image"`
	Engine string    `yaml:"engine"`
	Ports  []PortMap `yaml:"ports"`
}

func Load() (*Config, error) {
	data, err := os.ReadFile(MainConfig)
	if err != nil {
		return &Config{
			Engine: "docker", ListenAddr: "127.0.0.1",
			Port: 8080, DataDir: DataDir, PortRange: "50000-60000",
		}, nil
	}
	var c Config
	if err := yaml.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.ListenAddr == "" {
		c.ListenAddr = "127.0.0.1"
	}
	if c.Engine == "" {
		c.Engine = "docker"
	}
	if c.DataDir == "" {
		c.DataDir = DataDir
	}
	if c.PortRange == "" {
		c.PortRange = "50000-60000"
	}
	return &c, nil
}

func Save(c *Config) error {
	if err := os.MkdirAll(ConfigDir, 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(MainConfig, data, 0600)
}

func containerConfigPath(name string) string {
	return filepath.Join(ConfigDir, "containers", name+".yaml")
}

func LoadContainer(name string) (*ContainerConfig, error) {
	data, err := os.ReadFile(containerConfigPath(name))
	if err != nil {
		return nil, err
	}
	var c ContainerConfig
	return &c, yaml.Unmarshal(data, &c)
}

func SaveContainer(c *ContainerConfig) error {
	dir := filepath.Join(ConfigDir, "containers")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(containerConfigPath(c.Name), data, 0600)
}

func ListContainerConfigs() ([]*ContainerConfig, error) {
	dir := filepath.Join(ConfigDir, "containers")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	var out []*ContainerConfig
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".yaml" {
			continue
		}
		name := e.Name()[:len(e.Name())-5]
		if cc, err := LoadContainer(name); err == nil {
			out = append(out, cc)
		}
	}
	return out, nil
}

func DeleteContainer(name string) error {
	if err := os.Remove(containerConfigPath(name)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete container config: %w", err)
	}
	return nil
}

// AllocatePorts returns 4 external ports (for SSH=22, RDP=3389, NX=4000, VNC=5900)
// from portRange, skipping any already used.
func AllocatePorts(portRange string, usedPorts map[int]bool) ([4]int, error) {
	parts := strings.SplitN(portRange, "-", 2)
	if len(parts) != 2 {
		return [4]int{}, fmt.Errorf("invalid port_range: %s", portRange)
	}
	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || start < 1 || start > 65535 {
		return [4]int{}, fmt.Errorf("invalid port_range start: %s", parts[0])
	}
	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil || end < start || end > 65535 {
		return [4]int{}, fmt.Errorf("invalid port_range end: %s", parts[1])
	}
	var result [4]int
	taken := make(map[int]bool)
	for k, v := range usedPorts {
		taken[k] = v
	}
	for i := 0; i < 4; i++ {
		for p := start; p <= end; p++ {
			if !taken[p] {
				result[i] = p
				taken[p] = true
				break
			}
		}
		if result[i] == 0 {
			return result, fmt.Errorf("no available port in range %s", portRange)
		}
	}
	return result, nil
}
