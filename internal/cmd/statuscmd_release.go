package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/spf13/cobra"
)

// This represents a platform release as returned by the `status/release` endpoint.
type ReleasePackage struct {
	Elixir    interface{} `json:"elixir"`
	Hashref   string      `json:"hashref"`
	Name      string      `json:"name"`
	Repo      interface{} `json:"repo"`
	Timestamp string      `json:"timestamp"`
}

func (c *Config) runStatusReleaseE(cmd *cobra.Command, args []string) error {
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

		url := fmt.Sprintf(statusUrls[key], "release")

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

		var result map[string]ReleasePackage
		err = json.Unmarshal(body, &result)
		if err != nil {
			errors = append(errors, []string{key, err.Error()})
			continue
		}

		pkg := result["package"]

		results = append(results, []string{key, "v" + pkg.Timestamp + "-" + pkg.Hashref})
	}

	styles := newReportStyles(c, withWidth(80))

	reportResults([]string{"Env", "Tag"}, results, styles)

	if len(errors) > 0 && len(results) > 0 {
		fmt.Println()
	}

	reportErrors(errors, styles)

	return nil
}
