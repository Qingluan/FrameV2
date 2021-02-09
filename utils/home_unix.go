//+build linux darwin
package utils

import (
	"os"
	"path/filepath"
)

var (
	Home = filepath.Join(os.Getenv("HOME"), ".config", "ProxyAnyWhere")
)
