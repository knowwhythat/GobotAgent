//go:build !dev

package env

import (
	"os"
	"path/filepath"
)

func GetExecutablePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
