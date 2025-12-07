package scanner

import (
	"os"
	"path/filepath"
)

type Venv struct {
	Path string
	Size int64
}

func Scan(root string) ([]Venv, error) {
	var venvs []Venv

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If we can't access a path, just skip it
			return nil
		}

		if info.IsDir() {
			// Check for pyvenv.cfg
			if _, err := os.Stat(filepath.Join(path, "pyvenv.cfg")); err == nil {
				size, _ := getDirSize(path)
				venvs = append(venvs, Venv{Path: path, Size: size})
				return filepath.SkipDir // Don't look inside a venv
			}
            // Check for bin/activate
            if _, err := os.Stat(filepath.Join(path, "bin", "activate")); err == nil {
                 size, _ := getDirSize(path)
				venvs = append(venvs, Venv{Path: path, Size: size})
				return filepath.SkipDir
            }
            // Skip node_modules and .git to speed up
            if info.Name() == "node_modules" || info.Name() == ".git" {
                return filepath.SkipDir
            }
		}
		return nil
	})

	return venvs, err
}

func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}
