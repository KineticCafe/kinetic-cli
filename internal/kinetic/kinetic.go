package kinetic

import (
	"strings"

	"github.com/spf13/cobra"
)

func FlagCompletionFunc(allCompletions []string) func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		for _, completion := range allCompletions {
			if strings.HasPrefix(completion, toComplete) {
				completions = append(completions, completion)
			}
		}

		return completions, cobra.ShellCompDirectiveNoFileComp
	}
}
