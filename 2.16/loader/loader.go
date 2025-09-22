package loader

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Downloads(resources string, dir string) (string, error) {
	resp, err := http.Get(resources)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	fname := filepath.Base(resources)
	if fname == "" {
		fname = "index"
	}

	path := filepath.Join(dir, fname)
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return path, err
}
