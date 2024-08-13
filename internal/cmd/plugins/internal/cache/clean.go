package cache

import (
	"os"

	"buffalo-cli/internal/plugins"
	"github.com/spf13/cobra"
)

// CleanCmd cleans the plugins cache
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "cleans the plugins cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		os.RemoveAll(plugins.CachePath)
		return nil
	},
}
