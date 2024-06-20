package simutils

import (
	"os"
	"path/filepath"
)

func CurrentDirectory() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(ex)
}
