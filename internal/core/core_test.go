package core

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
)

type FakeHTTPClient struct {
	mockGet func(string) (*http.Response, error)
	// Get func(string) (*http.Response, error)
}

func (c *FakeHTTPClient) Get(url string) (*http.Response, error) {
	return c.mockGet(url)
}

func TestInstall(t *testing.T) {
	content := []byte("beep-bop-beep")
	body := ioutil.NopCloser(bytes.NewReader(content))
	var actualUrl string
	client := HTTPClient(&FakeHTTPClient{
		mockGet: func(u string) (*http.Response, error) {
			actualUrl = u
			return &http.Response{
				StatusCode: 200,
				Body:       body,
			}, nil
		},
	})

	installDir := path.Join(t.TempDir(), "bin")
	p := Pkgfy{Config: PkgfyConfig{InstallDir: installDir}, Client: &client}
	url := "https://example.com/gonzo"
	err := p.Install(url)
	if err != nil {
		t.Error(err)
	}
	if actualUrl != url {
		t.Errorf("Expected %s, got %s", url, actualUrl)
	}

	fi, err := os.Lstat(path.Join(installDir, "gonzo"))
	if err != nil {
		t.Error(err)
	}
	expectedMode := fs.FileMode(0755)
	if fi.Mode() != expectedMode {
		t.Errorf("Expected %v, got %v", expectedMode, fi.Mode())
	}

	actualContent, err := os.ReadFile(path.Join(installDir, "gonzo"))
	if err != nil {
		t.Error(err)
	}
	if string(actualContent) != string(content) {
		t.Errorf("Expected %s, got %s", content, actualContent)
	}
}
