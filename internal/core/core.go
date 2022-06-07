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

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

type Pkgfy struct {
	Config PkgfyConfig
	Client *HTTPClient
}

// Install downloads the given url and saves it in the config.InstallDir.
// It creates the config.InstallDir if needed. The file mode for the created
// dir and file is 0755.
func (p *Pkgfy) Install(url string) (err error) {
	// TODO: Make mode configurable
	err = os.MkdirAll(p.Config.InstallDir, 0755)
	if err != nil && !os.IsExist(err) {
		return
	}
	err = nil

	filepath, err := p.download(url, p.Config.InstallDir)
	// TODO: Make mode configurable
	err = os.Chmod(filepath, 0755)
	if err != nil {
		return
	}
	return
}

// download fetches the url and saves the response as a new file in the
// given dir. The filename is hard-coded to "gonzo". The dir is expected
// to already exist.
func (p *Pkgfy) download(url, dir string) (filepath string, err error) {
	// TODO: Extract this from the response header or the URL
	filename := "gonzo"

	rsp, err := (*p.Client).Get(url)
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
