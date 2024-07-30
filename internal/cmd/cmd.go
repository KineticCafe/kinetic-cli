package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/KineticCommerce/kinetic-cli/internal/kinetic"
	"github.com/KineticCommerce/kinetic-cli/internal/kineticset"
)

var deDuplicateErrorRx = regexp.MustCompile(`:\s+`)

type VersionInfo struct {
	Version string
	Commit  string
	Date    string
	BuiltBy string
}

func (v VersionInfo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("version", v.Version),
		slog.String("commit", v.Commit),
		slog.String("date", v.Date),
		slog.String("builtBy", v.BuiltBy),
	)
}

func Main(versionInfo VersionInfo, args []string) int {
	if err := runMain(versionInfo, args); err != nil {
		if errExitCode := kinetic.ExitCodeError(0); errors.As(err, &errExitCode) {
			return int(errExitCode)
		}

		fmt.Fprintf(os.Stderr, "kinetic-cli: %s\n", deDuplicateError(err))
		return 1
	}

	return 0
}

func runMain(versionInfo VersionInfo, args []string) (err error) {
	if versionInfo.Commit == "" || versionInfo.Date == "" {
		if buildInfo, ok := debug.ReadBuildInfo(); ok {
			var vcs, vcsRevision, vcsTime, vcsModified string
			for _, setting := range buildInfo.Settings {
				switch setting.Key {
				case "vcs":
					vcs = setting.Value
				case "vcs.revision":
					vcsRevision = setting.Value
				case "vcs.time":
					vcsTime = setting.Value
				case "vcs.modified":
					vcsModified = setting.Value
				}
			}
			if versionInfo.Commit == "" && vcs == "git" {
				versionInfo.Commit = vcsRevision
				if modified, err := strconv.ParseBool(vcsModified); err == nil && modified {
					versionInfo.Commit += "-dirty"
				}
			}
			if versionInfo.Date == "" {
				versionInfo.Date = vcsTime
			}
		}
	}

	var config *Config
	if config, err = newConfig(
		withVersionInfo(versionInfo)); err != nil {
		return err
	}

	return config.execute(args)
}

func deDuplicateError(err error) string {
	components := deDuplicateErrorRx.Split(err.Error(), -1)
	seenComponents := kineticset.NewWithCapacity[string](len(components))
	uniqueComponents := make([]string, 0, len(components))
	for _, component := range components {
		if seenComponents.Contains(component) {
			continue
		}
		uniqueComponents = append(uniqueComponents, component)
		seenComponents.Add(component)
	}
	return strings.Join(uniqueComponents, ": ")
}
