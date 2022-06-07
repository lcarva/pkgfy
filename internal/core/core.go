package core

import (
	"io"
	"net/http"
	"os"
	"path"
)

type PkgfyConfig struct {
	InstallDir string
}

// Install downloads the given url and saves it in the config.InstallDir.
// It also sets the file mode to 0755.
func Install(url string, config *PkgfyConfig) (err error) {
	// TODO: Create dir if needed
	filepath, err := download(url, config.InstallDir)
	// TODO: Make mode configurable
	err = os.Chmod(filepath, 0755)
	if err != nil {
		return
	}
	return
}

// download fetches the url and saves the response as a new file in the
// given dir. The filename is hard-coded to "gonzo".
func download(url string, dir string) (filepath string, err error) {
	// TODO: Extract this from the response header or the URL
	filename := "gonzo"

	rsp, err := http.Get(url)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	filepath = path.Join(dir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return
	}

	return
}
