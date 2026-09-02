package engine

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type LxcEngine struct{ isLxd bool }

func (l *LxcEngine) List() ([]ContainerInfo, error) {
	var out []byte
	var err error
	if l.isLxd {
		out, err = exec.Command("lxc", "list", "--format", "csv", "-c", "ns").Output()
	} else {
		out, err = exec.Command("lxc-ls", "-f", "--fancy-format", "NAME,STATE").Output()
	}
	if err != nil {
		return nil, err
	}
	var list []ContainerInfo
	for i, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" || (!l.isLxd && i == 0) {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 1 {
			continue
		}
		info := ContainerInfo{Name: strings.Split(parts[0], ",")[0]}
		if l.isLxd {
			sp := strings.SplitN(line, ",", 2)
			info.Name = sp[0]
			if len(sp) > 1 {
				info.Status = sp[1]
			}
		} else if len(parts) > 1 {
			info.Status = parts[1]
		}
		list = append(list, info)
	}
	return list, nil
}

func (l *LxcEngine) Pull(image string) error {
	if l.isLxd {
		c := exec.Command("lxc", "image", "copy", image, "local:")
		c.Stdout, c.Stderr = os.Stdout, os.Stderr
		return c.Run()
	}
	return fmt.Errorf("pull not supported for native lxc")
}

func (l *LxcEngine) Start(name string) error {
	if l.isLxd {
		return exec.Command("lxc", "start", name).Run()
	}
	return exec.Command("lxc-start", "-n", name).Run()
}

func (l *LxcEngine) Stop(name string) error {
	if l.isLxd {
		return exec.Command("lxc", "stop", name).Run()
	}
	return exec.Command("lxc-stop", "-n", name).Run()
}

func (l *LxcEngine) Restart(name string) error {
	_ = l.Stop(name)
	return l.Start(name)
}

func (l *LxcEngine) Remove(name string) error {
	if l.isLxd {
		return exec.Command("lxc", "delete", "--force", name).Run()
	}
	return exec.Command("lxc-destroy", "-n", name).Run()
}

func (l *LxcEngine) Exec(name string, cmd []string, detach bool) error {
	var c *exec.Cmd
	if l.isLxd {
		c = exec.Command("lxc", append([]string{"exec", name, "--"}, cmd...)...)
	} else {
		c = exec.Command("lxc-attach", append([]string{"-n", name, "--"}, cmd...)...)
	}
	if !detach {
		c.Stdin, c.Stdout, c.Stderr = os.Stdin, os.Stdout, os.Stderr
	}
	return c.Run()
}

func (l *LxcEngine) Run(image, name string, ports []string, env []string, devices []string) error {
	if l.isLxd {
		if err := exec.Command("lxc", "init", image, name).Run(); err != nil {
			return err
		}
		return exec.Command("lxc", "start", name).Run()
	}
	return exec.Command("lxc-create", "-n", name, "-t", image).Run()
}

func (l *LxcEngine) SetPassword(name, password string) error {
	return l.Exec(name, []string{"bash", "-c", fmt.Sprintf("echo 'root:%s' | chpasswd", password)}, false)
}

func (l *LxcEngine) Info(name string) (*ContainerInfo, error) {
	var out []byte
	var err error
	if l.isLxd {
		out, err = exec.Command("lxc", "info", name).Output()
	} else {
		out, err = exec.Command("lxc-info", "-n", name).Output()
	}
	if err != nil {
		return nil, err
	}
	info := &ContainerInfo{Name: name}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(strings.ToLower(line), "status:") || strings.Contains(strings.ToLower(line), "state:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				info.Status = strings.TrimSpace(parts[1])
			}
		}
	}
	return info, nil
}

func (l *LxcEngine) GetIP(name string) (string, error) {
	if l.isLxd {
		out, err := exec.Command("lxc", "list", name, "--format", "csv", "-c", "4").Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(out)), nil
	}
	out, err := exec.Command("lxc-info", "-n", name, "-iH").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// ─── New image / stats methods ────────────────────────────────────────────────

func (l *LxcEngine) ListImages() ([]ImageInfo, error) {
	if !l.isLxd {
		return nil, fmt.Errorf("list images not supported for native LXC")
	}
	out, err := exec.Command("lxc", "image", "list", "--format", "csv", "-c", "flsd").Output()
	if err != nil {
		return nil, err
	}
	var list []ImageInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 4)
		img := ImageInfo{}
		if len(parts) > 0 {
			img.Tags = []string{parts[0]}
		}
		if len(parts) > 1 {
			img.ID = parts[1]
		}
		if len(parts) > 2 {
			img.Size = parts[2]
		}
		if len(parts) > 3 {
			img.Created = parts[3]
		}
		list = append(list, img)
	}
	return list, nil
}

func (l *LxcEngine) RemoveImage(id string) error {
	if !l.isLxd {
		return fmt.Errorf("remove image not supported for native LXC")
	}
	return exec.Command("lxc", "image", "delete", id).Run()
}

func (l *LxcEngine) GetStats() (*EngineStats, error) {
	stats := &EngineStats{EngineVersion: l.Version()}
	containers, err := l.List()
	if err != nil {
		return stats, nil
	}
	stats.ContainersTotal = len(containers)
	for _, c := range containers {
		lower := strings.ToLower(c.Status)
		if lower == "running" || lower == "started" {
			stats.ContainersRunning++
		}
	}
	if l.isLxd {
		imgs, err := l.ListImages()
		if err == nil {
			stats.ImagesTotal = len(imgs)
		}
	}
	return stats, nil
}

func (l *LxcEngine) Version() string {
	binary := "lxc-ls"
	args := []string{"--version"}
	if l.isLxd {
		binary = "lxc"
		args = []string{"version"}
	}
	out, err := exec.Command(binary, args...).Output()
	if err != nil {
		return "lxc (unavailable)"
	}
	return strings.TrimSpace(string(out))
}
