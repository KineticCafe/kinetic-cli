package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// This is still imperfect, but provides a useful starting point.
//
// 1.  Add a `--json` output flag to these commands to dump the raw JSON value.
//     If color is not disabled, make it syntax highlighted (possibly with
//     github.com/alecthomas/chroma, since the coloring in github.com/itchyny/gojq appears
//     to only be available in the cli).
//
// 2.  Add a `full` command that uses `status/full` instead of calling three times the way
//     the default function does now.
//
// 3.  Migrate status collectors to `internal/kinetic` or `internal/platform`.
//
// 4.  Add status collectors for sidecars (see `kinetic.sh`) and make it configurable. The
//     default configuration should probably be in an S3 bucket so we can update the
//     behaviour as required.
//
//     Adjust the commands as necessary.

var (
	urlOrder = []string{"local", "dit", "staging", "prod", "prod-eu"}

	statusUrls = map[string]string{
		"local":   "https://localhost:4000/status/%s",
		"dit":     "https://kcs-dev.kineticcommercetech.io/status/%s",
		"staging": "https://kcs-staging.kineticcommercetech.io/status/%s",
		"prod":    "https://kcs.kineticcommerce.io/status/%s",
		"prod-eu": "https://kcs-prod-eu-platform.kineticcommerce.io/status/%s",
	}
)

func (c *Config) newStatusCmd() *cobra.Command {
	statusCmd := &cobra.Command{
		Use:   "status [command]",
		Short: "Query platform system status",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(c.stdout, "Release:")
			if err := c.runStatusReleaseE(cmd, args); err != nil {
				return err
			}

			fmt.Fprintln(c.stdout, "\nConfig:")
			if err := c.runStatusConfigE(cmd, args); err != nil {
				return err
			}

			fmt.Fprintln(c.stdout, "\nSchema:")
			if err := c.runStatusSchemaE(cmd, args); err != nil {
				return err
			}

			return nil
		},
	}

	statusCmd.AddCommand(&cobra.Command{
		Use:   "config [environments...]",
		Short: "Display the current config package",
		RunE:  c.runStatusConfigE,
	})

	statusCmd.AddCommand(&cobra.Command{
		Use:   "schema [environments...]",
		Short: "Display the current Sqitch migrations",
		RunE:  c.runStatusSchemaE,
	})

	statusCmd.AddCommand(&cobra.Command{
		Use:   "release [environments...]",
		Short: "Display the current release package",
		RunE:  c.runStatusReleaseE,
	})

	return statusCmd
}
