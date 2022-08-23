package version

import (
	"runtime"
	"runtime/debug"
)

// Build information. Populated at build-time.
//
// nolint:gochecknoglobals
var (
	version      = "dev"
	revision     string
	branch       string
	buildUser    string
	buildDate    string
	dependencies []*debug.Module
)

// Information holds app version info.
type Information struct {
	Version      string
	Revision     string
	Branch       string
	BuildUser    string
	BuildDate    string
	GoVersion    string
	GoOS         string
	GoArch       string
	Dependencies []*debug.Module
}

// Info returns app version info.
func Info() Information {
	return Information{
		Version:      version,
		Revision:     revision,
		Branch:       branch,
		BuildUser:    buildUser,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		GoOS:         runtime.GOOS,
		GoArch:       runtime.GOARCH,
		Dependencies: dependencies,
	}
}

//nolint:gochecknoinits
func init() {
	if info, available := debug.ReadBuildInfo(); available {
		dependencies = info.Deps
	}
}
