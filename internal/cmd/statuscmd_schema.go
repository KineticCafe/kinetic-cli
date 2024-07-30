package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/spf13/cobra"
)

// This represents a sqitch migration as returned by the `status/sqitch` endpoint.
type SqitchMigration struct {
	Change    string `json:"change"`
	ChangeID  string `json:"change_id"`
	PlannedAt string `json:"planned_at"`
}

var shortSchemaMap = map[string]string{
	"kinetic-cas-kiehls-schema": "kiehls",
	"kinetic-platform-schema":   "core",
}

func schemaName(env, name string) string {
	short := shortSchemaMap[name]

	if short == "" {
		short = name
	}

	return env + " - " + short
}

func (c *Config) runStatusSchemaE(cmd *cobra.Command, args []string) error {
	filter := len(args) > 0
	results := make([][]string, 0, len(statusUrls))
	errors := make([][]string, 0, len(statusUrls))

	for _, key := range urlOrder {
		if filter && slices.Index(args, key) == -1 {
			continue
		}

		if key == "local" && slices.Index(args, key) == -1 {
			continue
		}

		url := fmt.Sprintf(statusUrls[key], "sqitch")

		resp, err := http.Get(url)
		if err != nil {
			errors = append(errors, []string{key, err.Error()})
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errors = append(errors, []string{key, fmt.Sprintf("Expected 200 OK, got %d", resp.StatusCode)})
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errors = append(errors, []string{key, err.Error()})
			continue
		}

		var result map[string]SqitchMigration
		err = json.Unmarshal(body, &result)
		if err != nil {
			errors = append(errors, []string{key, err.Error()})
			continue
		}

		for name, migration := range result {
			if migration.ChangeID == "" {
				continue
			}

			results = append(results, []string{
				schemaName(key, name),
				migration.ChangeID,
				migration.PlannedAt,
			})
		}
	}

	styles := newReportStyles(c)

	reportResults([]string{"Env/Schema", "Change ID", "Planned At"}, results, styles)

	if len(errors) > 0 && len(results) > 0 {
		fmt.Println()
	}

	reportErrors(errors, styles)

	return nil
}
