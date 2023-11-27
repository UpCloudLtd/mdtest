//nolint:gochecknoglobals // These are defined as var instead of const to be able to override them on compile time.
package globals

import (
	"fmt"
	"runtime/debug"
)

var (
	BuildDate = "unknown"
	Version   = "dev"
)

func GetVersion() string {
	// Version was overridden during the build
	if Version != "dev" {
		return Version
	}

	// Try to read version from build info
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		version := buildInfo.Main.Version
		if version != "(devel)" && version != "" {
			return version
		}

		settingsMap := make(map[string]string)
		for _, setting := range buildInfo.Settings {
			settingsMap[setting.Key] = setting.Value
		}

		version = "dev"
		if rev, ok := settingsMap["vcs.revision"]; ok {
			version = fmt.Sprintf("%s-%s", version, rev[:8])
		}

		if dirty, ok := settingsMap["vcs.modified"]; ok && dirty == "true" {
			return version + "-dirty"
		}

		return version
	}

	// Fallback to the default value
	return Version
}
