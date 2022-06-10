package core

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

type Package struct {
	Name  string
	Alias string
	URL   string
}

// Install downloads the given url and saves it in the config.InstallDir with pkg.Name. It creates
// the config.InstallDir if needed. The file mode for the created dir and file is 0755.
func (p *Pkgfy) Install(pkg Package) (err error) {
	// TODO: Make mode configurable
	err = os.MkdirAll(p.Config.InstallDir, 0755)
	if err != nil {
		return
	}

	filePath := path.Join(p.Config.InstallDir, pkg.Name)
	err = p.download(pkg.URL, filePath)
	// TODO: Make mode configurable
	err = os.Chmod(filePath, 0755)
	if err != nil {
		return
	}

	if pkg.Alias != "" {
		err = symlinkAlias(pkg.Alias, filePath)
		if err != nil {
			return
		}
	}
	return
}

// download fetches the url and saves the response as a new file in the given dest.
func (p *Pkgfy) download(url, dest string) (err error) {
	rsp, err := (*p.Client).Get(url)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	tempDest := dest + ".part"
	file, err := os.Create(tempDest)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return
	}

	err = os.Rename(tempDest, dest)
	return
}

// symlinkAlias create a new symlink with the given alias for the target path.
func symlinkAlias(alias, target string) (err error) {
	aliasPath := path.Join(filepath.Dir(target), alias)
	return os.Symlink(target, aliasPath)
}
