//+build windows
package utils

import "os"

var (
	Home, _ = os.UserHomeDir()
)
