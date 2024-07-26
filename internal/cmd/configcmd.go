package cmd

import "github.com/spf13/cobra"

// 1. Add `op` support to pull the configuration files from 1Password. The `op://`
//    document URLs will be hardcoded but overridable via configuration.
// 2. Add `SSM` and `S3` support to publish the configurations like the `manage-secrets`
//    (livechat associate) / `manage-config` (employee web) / `manage` (import triggers)
//    scripts do.
// 3. All of this should be driven by configuration (default configuration can be in S3)
//    because not all configuration files are written the same (see the above scripts) and
//    the code that is currently written is only written to *wrap* a single configuration
//    file already on disk. Publication is currently manual.

func (c *Config) newConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config [command]",
		Short: "Manage platform configuration",
	}

	configCmd.AddCommand(c.newConfigWrapJsonSubcmd())

	return configCmd
}
