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
	var testCases = []struct {
		name string
		pkg  Package
	}{
		{"simple installation",
			Package{Name: "gonzo", URL: "https://example.com/x"}},
		{"installation with an alias",
			Package{Name: "gonzo", URL: "https://example.com/gonzo", Alias: "gonzoo"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing %s", tc.name)
			content := "hello, gonzo!"
			body := ioutil.NopCloser(bytes.NewReader([]byte(content)))
			var actualUrl string
			parsedUrl, err := url.Parse(tc.pkg.URL)
			assert.Nil(t, err)
			client := HTTPClient(&FakeHTTPClient{
				mockGet: func(u string) (*http.Response, error) {
					actualUrl = u
					return &http.Response{
						StatusCode: 200,
						Body:       body,
						Request:    &http.Request{URL: parsedUrl},
					}, nil
				},
			})

			installDir := path.Join(t.TempDir(), "bin")
			p := Pkgfy{Config: PkgfyConfig{InstallDir: installDir}, Client: &client}

			err = p.Install(tc.pkg)
			assert.Nil(t, err)
			assert.Equal(t, tc.pkg.URL, actualUrl)

			files, err := ioutil.ReadDir(installDir)
			assert.Nil(t, err)
			actualFiles := []string{}
			for _, f := range files {
				actualFiles = append(actualFiles, f.Name())
				switch f.Name() {
				case tc.pkg.Name:
					expectedMode := fs.FileMode(0755)
					assert.Equal(t, expectedMode, f.Mode())
				case tc.pkg.Alias:
					// Symlinks do their own thing when it comes to file permissions.
				default:
					t.Fatalf("Found unexpected file in install dir: %q", f.Name())
				}
			}
			expectedFiles := []string{tc.pkg.Name}
			if tc.pkg.Alias != "" {
				expectedFiles = append(expectedFiles, tc.pkg.Alias)
			}
			assert.Equal(t, expectedFiles, actualFiles)

			expectedFilename := path.Join(installDir, tc.pkg.Name)
			actualContent, err := os.ReadFile(expectedFilename)
			assert.Nil(t, err)
			assert.Equal(t, content, string(actualContent))
		})
	}
}
