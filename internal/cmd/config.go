package cmd

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/KineticCommerce/kinetic-cli/internal/kinetic"
	"github.com/KineticCommerce/kinetic-cli/internal/kineticerrors"
	"github.com/coreos/go-semver/semver"
	"github.com/spf13/cobra"
)

// This contains all data settable in the configuration file
type ConfigFile struct {
	Verbose bool `json:"verbose"         mapstructure:"verbose"         yaml:"verbose"`
}

// This contains all configuration.
type Config struct {
	ConfigFile

	configFilePath string
	debug          bool
	environment    environmentType
	version        semver.Version
	versionInfo    VersionInfo
	versionStr     string

	homeDir string
	wd      string
	project kinetic.Project

	logger *slog.Logger
	stderr io.Writer
	stdin  io.Reader
	stdout io.Writer
}

// A configOption sets and option on a Config.
type configOption func(*Config) error

func (c *Config) newRootCmd() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "kinetic-cli",
		Short: "Kinetic platform and sidecar management tools",
	}

	persistentFlags := rootCmd.PersistentFlags()

	persistentFlags.BoolVarP(&c.Verbose, "verbose", "v", c.Verbose, "Make output more verbose")
	persistentFlags.BoolVarP(&c.debug, "debug", "d", false, "Turn on debug output")
	persistentFlags.StringVarP(&c.configFilePath, "config", "c", "", "Set config file")
	persistentFlags.VarP(&c.environment, "environment", "e", "Set the environment (dit, stage, prod, prod-eu)")

	if err := kineticerrors.Combine(
		rootCmd.MarkPersistentFlagFilename("config"),
		persistentFlags.MarkHidden("debug"),
		rootCmd.RegisterFlagCompletionFunc("environment", environmentTypeCompletionFunc),
	); err != nil {
		return nil, err
	}

	for _, cmd := range []*cobra.Command{
		c.newStatusCmd(),
	} {
		if cmd != nil {
			rootCmd.AddCommand(cmd)
		}
	}

	return rootCmd, nil
}

func newConfig(options ...configOption) (*Config, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homeDir, err := filepath.Abs(userHomeDir)
	if err != nil {
		return nil, err
	}

	logger := slog.Default()

	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	project, err := kinetic.ResolveProject(wd)
	if err != nil {
		return nil, err
	}

	c := &Config{
		ConfigFile: ConfigFile{
			Verbose: false,
		},
		homeDir: homeDir,
		logger:  logger,
		wd:      wd,
		project: project,
		stdin:   os.Stdin,
		stdout:  os.Stdout,
		stderr:  os.Stderr,
	}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Config) execute(args []string) error {
	rootCmd, err := c.newRootCmd()
	if err != nil {
		return err
	}

	rootCmd.SetArgs(args)
	return rootCmd.Execute()
}

func withVersionInfo(versionInfo VersionInfo) configOption {
	return func(c *Config) error {
		var version *semver.Version
		var versionElems []string

		if versionInfo.Version != "" {
			var err error

			version, err = semver.NewVersion(strings.TrimPrefix(versionInfo.Version, "v"))
			if err != nil {
				return err
			}

			versionElems = append(versionElems, "v"+version.String())
		} else {
			versionElems = append(versionElems, "dev")
		}

		if versionInfo.Commit != "" {
			versionElems = append(versionElems, "commit "+versionInfo.Commit)
		}

		if versionInfo.Date != "" {
			date := versionInfo.Date
			if sec, err := strconv.ParseInt(date, 10, 64); err == nil {
				date = time.Unix(sec, 0).UTC().Format(time.RFC3339)
			}
			versionElems = append(versionElems, "built at "+date)
		}

		if versionInfo.BuiltBy != "" {
			versionElems = append(versionElems, "built by "+versionInfo.BuiltBy)
		}

		if version != nil {
			c.version = *version
		}
		c.versionInfo = versionInfo
		c.versionStr = strings.Join(versionElems, ", ")
		return nil
	}
}
