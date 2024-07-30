package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/spf13/cobra"
)

type ConfigStatusResponse struct {
	Timestamp string `json:"timestamp"`
	Hashref   string `json:"hashref"`
	Version   int    `json:"version"`
}

func (c *Config) runStatusConfigE(cmd *cobra.Command, args []string) error {
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

		url := fmt.Sprintf(statusUrls[key], "config")

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

		var result ConfigStatusResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			errors = append(errors, []string{key, err.Error()})
			continue
		}

		if result.Version == 0 {
			result.Version = 1
		}

		results = append(results, []string{key, result.Timestamp, result.Hashref, fmt.Sprint(result.Version)})
	}

	styles := newReportStyles(
		c,
		withWidth(80),
	)

	reportResults([]string{"Env", "Timestamp", "Hashref", "Version"}, results, styles)

	if len(errors) > 0 && len(results) > 0 {
		fmt.Println()
	}

	reportErrors(errors, styles)

	return nil
}
