package post_logic

import (
	"os"
	"path/filepath"
	"time"
)

func getDirectoryState(dir string) (map[string]time.Time, error) {
	state := make(map[string]time.Time)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// We only care about files, not directories.
		if !info.IsDir() {
			state[path] = info.ModTime()
		}
		return nil
	})
	return state, err
}
