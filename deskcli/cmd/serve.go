package cmd

import (
	"github.com/pikachuim/deskcli/internal/api"
	"github.com/pikachuim/deskcli/internal/config"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		// Set default static dir to embedded web/dist
		if cfg.StaticDir == "" {
			cfg.StaticDir = "/opt/deskcli/web/dist"
		}
		return api.Start(cfg)
	},
}
