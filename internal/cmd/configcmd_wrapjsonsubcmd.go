package cmd

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func (c *Config) newConfigWrapJsonSubcmd() *cobra.Command {
	subcmd := &cobra.Command{
		Use:   "wrapjson targets...",
		Short: "Wraps the config JSON in the Kinetic Platform v2 JSON config structure",
		Args:  cobra.ExactArgs(1),
		RunE:  c.runConfigWrapJsonE,
	}

	return subcmd
}

type wrappedConfig struct {
	Digest    string `json:"__digest__"`
	Timestamp string `json:"__timestamp__"`
	Version   int    `json:"__version__"`
	Config    string `json:"__config__"`
}

func (c *Config) runConfigWrapJsonE(cmd *cobra.Command, args []string) error {
	hasErrors := false

	t := time.Now()
	timestamp := t.Format("20060102150405")

	for _, file := range args {
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				c.logger.Error(fmt.Sprintf("%s: file does not exist", file))
				hasErrors = true
				continue
			}

			return err
		}

		data, err := os.ReadFile(file)
		if err != nil {
			c.logger.Error(fmt.Sprintf("%s: error reading file: %e", file, err))
			hasErrors = true
			continue
		}

		var testJson map[string]interface{}
		err = json.Unmarshal(data, &testJson)
		if err != nil {
			c.logger.Error(fmt.Sprintf("%s: error verifying JSON data: %e", file, err))
			hasErrors = true
			continue
		}

		var builder strings.Builder
		encoder := json.NewEncoder(&builder)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "")

		if err := encoder.Encode(testJson); err != nil {
			c.logger.Error(fmt.Sprintf("%s: error encoding JSON data: %e", file, err))
			hasErrors = true
			continue
		}

		configJson := strings.TrimSuffix(builder.String(), "\n")
		digest := fmt.Sprintf("sha256:%x", sha256.Sum256([]byte(configJson)))

		wrapped := &wrappedConfig{
			Config:    configJson,
			Digest:    digest,
			Timestamp: timestamp,
			Version:   2,
		}

		wrappedData, err := json.Marshal(wrapped)
		if err != nil {
			c.logger.Error(fmt.Sprintf("%s: error creating wrapped config: %e", file, err))
			hasErrors = true
			continue
		}

		dirname := filepath.Dir(file)
		extname := filepath.Ext(file)
		basename := strings.TrimSuffix(filepath.Base(file), extname)

		newFile := filepath.Join(dirname, fmt.Sprintf("%s-%s%s", basename, timestamp, extname))

		if err = os.WriteFile(newFile, wrappedData, 0600); err != nil {
			c.logger.Error(fmt.Sprintf("%s: error writing wrapped config (%s): %e", file, newFile, err))
			hasErrors = true
			continue
		}

		c.logger.Info(fmt.Sprintf("%s: wrapped as %s", file, newFile))
	}

	if hasErrors {
		os.Exit(1)
	}

	return nil
}
