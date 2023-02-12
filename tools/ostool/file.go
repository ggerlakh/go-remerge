package ostool

import (
	"errors"
	"log"
	"os"
)

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		} else {
			log.Fatalf("Error stat file %s: %v\n", path, err)
		}
	}
	return true
}
