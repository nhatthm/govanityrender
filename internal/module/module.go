package module

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	xerrors "go.nhat.io/vanityrender/internal/errors"
	"go.nhat.io/vanityrender/internal/must"
)

const (
	// ErrInvalidVersion indicates that the version is invalid.
	ErrInvalidVersion = xerrors.Error("invalid version")
)

var (
	// PathVersionRegExp matches a module version string.
	PathVersionRegExp = regexp.MustCompile(`^([a-zA-Z0-9]+([a-zA-Z0-9_/]+)?)?v\d+\.\d+\.\d+$`)
	// VersionRegExp matches a version string.
	VersionRegExp = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)$`)
)

const goMod = `go.mod`

// Path is the path to the module.
type Path string

// Version is the version of the module.
type Version struct {
	Major int
	Minor int
	Patch int
}

// String returns the string representation of the version.
func (v Version) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// PathVersion returns the path and version from a module version string.
func PathVersion(s string) (Path, Version) {
	if PathVersionRegExp.MatchString(s) {
		i := strings.LastIndex(s, "/")

		var (
			path    Path
			version Version
		)

		if i == -1 {
			path, version = ".", NewVersionFromString(s)
		} else {
			path = Path(s[:i])
			version = NewVersionFromString(s[i+1:])
		}

		if version.Major > 1 {
			path = Path(strings.TrimLeft(fmt.Sprintf("%s/v%d", strings.TrimLeft(string(path), "."), version.Major), "/"))
		}

		return path, version
	}

	return "", NewVersion(0, 0, 0)
}

// NewVersionFromString returns a new version from a string.
func NewVersionFromString(s string) Version {
	m := VersionRegExp.FindStringSubmatch(s)

	if len(m) == 0 {
		panic(fmt.Errorf("invalid version: %w", ErrInvalidVersion))
	}

	major, err := strconv.Atoi(m[1])
	must.NoError(err)

	minor, err := strconv.Atoi(m[2])
	must.NoError(err)

	patch, err := strconv.Atoi(m[3])
	must.NoError(err)

	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

// NewVersion returns a new version.
func NewVersion(major, minor, patch int) Version {
	return Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}
