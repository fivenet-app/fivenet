// Package version provides the version information for the application.
package version

const (
	// Owner is the owner of the source repository.
	Owner = "fivenet-app"
	// Repo is the name of the source repository.
	Repo = "fivenet"
)

// Version is the current version of the application (set during build).
var Version = "UNKNOWN"

// UnknownVersion is the default version string used when the version is not set.
const UnknownVersion = "UNKNOWN"
