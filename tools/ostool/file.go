package ostool

import (
	"errors"
	"log"
	"os"
)

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
