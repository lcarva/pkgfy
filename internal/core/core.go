package core

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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

// Install downloads the given url and saves it in the config.InstallDir. It creates the
// config.InstallDir if needed. The file mode for the created dir and file is 0755.
func (p *Pkgfy) Install(url string) (err error) {
	// TODO: Make mode configurable
	err = os.MkdirAll(p.Config.InstallDir, 0755)
	if err != nil {
		return
	}

	filePath, err := p.download(url, p.Config.InstallDir)
	// TODO: Make mode configurable
	err = os.Chmod(filePath, 0755)
	if err != nil {
		return
	}
	return
}

// download fetches the url and saves the response as a new file in the given dir. The dir is
// expected to already exist.
func (p *Pkgfy) download(url, dir string) (filePath string, err error) {
	rsp, err := (*p.Client).Get(url)
	if err != nil {
		return
	}
	defer rsp.Body.Close()

	filename, err := extractFilename(rsp)
	if err != nil {
		return
	}

	filePath = path.Join(dir, filename)

	tempFilePath := filePath + ".part"
	file, err := os.Create(tempFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, rsp.Body)
	if err != nil {
		return
	}

	err = os.Rename(tempFilePath, filePath)
	return
}

// extractFilename retrieves a filename from the response by inspecting the Content-Disposition
// header and falling back to the last path element in the URL.
func extractFilename(response *http.Response) (filename string, err error) {
	contentDisposition := response.Header.Get("content-disposition")
	if contentDisposition != "" {
		filenameRgx := regexp.MustCompile(`filename=(.+)`)
		for _, candidates := range filenameRgx.FindAllStringSubmatch(contentDisposition, -1) {
			candidate := basePath(candidates[len(candidates)-1])
			candidate = strings.Trim(candidate, `"`)
			if candidate != "" {
				filename = candidate
			}
		}
	}
	if filename == "" {
		filename = basePath(response.Request.URL.Path)
	}
	if filename == "" {
		// We're really out of options here. Use the hashed value of the url to avoid collisions.
		filename = fmt.Sprintf("%x", md5.Sum([]byte(response.Request.URL.Redacted())))
	}
	return
}

// basePath returns the last path element of fullPath. If one is not found, "" is returned.
func basePath(fullPath string) (base string) {
	candidate := filepath.Base(fullPath)
	if candidate != "." && candidate != ".." && candidate != string(os.PathSeparator) {
		base = candidate
	}
	return
}
