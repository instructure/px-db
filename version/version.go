package version

import "fmt"

const version = "v0.1.0"

// Print the version of the current release
func Print() string {
	return fmt.Sprintf("%s", version)
}
