package core

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeHTTPClient struct {
	mockGet func(string) (*http.Response, error)
	// Get func(string) (*http.Response, error)
}

func (c *FakeHTTPClient) Get(u string) (*http.Response, error) {
	return c.mockGet(u)
}

func TestInstall(t *testing.T) {
	var tests = []struct {
		name     string
		filename string
		fileurl  string
		headers  *http.Header
		alias    string
	}{
		{"filename from URL", "gonzo", "https://example.com/gonzo", &http.Header{}, ""},
		{"filename from URL, impartial Content-Disposition",
			"gonzo", "https://example.com/gonzo", cdHeader([]string{"attachement"}), ""},
		{"filename from Content-Disposition",
			"gonzo", "https://example.com/x", cdHeader([]string{`attachment; filename="gonzo"`}), ""},
		{"filename from Content-Disposition without quotes",
			"gonzo", "https://example.com/x", cdHeader([]string{`attachment; filename=gonzo`}), ""},
		{"filename from first Content-Disposition", "gonzo", "https://example.com/x",
			cdHeader([]string{`attachment; filename="gonzo"`, `attachment; filename="y"`}), ""},
		{"filename from URL hash", "182ccedb33a9e03fbf1079b209da1a31", "https://example.com/", &http.Header{}, ""},
		{"include alias", "gonzo", "https://example.com/gonzo", &http.Header{}, "gonzoo"},
	}

	for _, test := range tests {
		t.Logf("Testing %s", test.name)
		content := "hello, gonzo!"
		body := ioutil.NopCloser(bytes.NewReader([]byte(content)))
		var actualUrl string
		parsedUrl, err := url.Parse(test.fileurl)
		assert.Nil(t, err)
		client := HTTPClient(&FakeHTTPClient{
			mockGet: func(u string) (*http.Response, error) {
				actualUrl = u
				return &http.Response{
					StatusCode: 200,
					Header:     *test.headers,
					Body:       body,
					Request:    &http.Request{URL: parsedUrl},
				}, nil
			},
		})

		installDir := path.Join(t.TempDir(), "bin")
		p := Pkgfy{Config: PkgfyConfig{InstallDir: installDir}, Client: &client}

		err = p.Install(test.fileurl, &InstallOptions{Alias: test.alias})
		assert.Nil(t, err)
		assert.Equal(t, test.fileurl, actualUrl)

		files, err := ioutil.ReadDir(installDir)
		assert.Nil(t, err)
		actualFiles := map[string]bool{}
		for _, f := range files {
			actualFiles[f.Name()] = true
			switch f.Name() {
			case test.filename:
				expectedMode := fs.FileMode(0755)
				assert.Equal(t, expectedMode, f.Mode())
			case test.alias:
				// Symlinks do their own thing when it comes to file permissions.
			default:
				t.Fatalf("Found unexpected file in install dir: %q", f.Name())
			}
		}
		expectedFiles := map[string]bool{test.filename: true}
		if test.alias != "" {
			expectedFiles[test.alias] = true
		}
		assert.Equal(t, expectedFiles, actualFiles)

		expectedFilename := path.Join(installDir, test.filename)
		actualContent, err := os.ReadFile(expectedFilename)
		assert.Nil(t, err)
		assert.Equal(t, content, string(actualContent))
	}
}

// cdHeader returns Headers with given Content-Disposition values.
func cdHeader(values []string) *http.Header {
	return &http.Header{http.CanonicalHeaderKey("content-disposition"): values}
}
