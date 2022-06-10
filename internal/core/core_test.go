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
)

type FakeHTTPClient struct {
	mockGet func(string) (*http.Response, error)
	// Get func(string) (*http.Response, error)
}

func (c *FakeHTTPClient) Get(url string) (*http.Response, error) {
	return c.mockGet(url)
}

func TestInstall(t *testing.T) {
	var tests = []struct {
		filename string
		fileurl  string
		headers  *http.Header
	}{
		{"gonzo", "https://example.com/gonzo", &http.Header{}},
		// Content-Disposition header does not include a filename
		{"gonzo", "https://example.com/gonzo", cdHeader([]string{"attachement"})},
		// Content-Disposition header includes a filename
		{"gonzo", "https://example.com/x", cdHeader([]string{`attachment; filename="gonzo"`})},
		// Content-Disposition header includes a filename without quotes
		{"gonzo", "https://example.com/x", cdHeader([]string{`attachment; filename=gonzo`})},
		// Multiple Content-Disposition headers include different filenames
		{"gonzo", "https://example.com/x",
			cdHeader([]string{`attachment; filename="gonzo"`, `attachment; filename="y"`})},
		// Default to hash of URL
		{"182ccedb33a9e03fbf1079b209da1a31", "https://example.com/", &http.Header{}},
	}

	for _, test := range tests {
		content := "hello, gonzo!"
		body := ioutil.NopCloser(bytes.NewReader([]byte(content)))
		var actualUrl string
		parsedUrl, err := url.Parse(test.fileurl)
		if err != nil {
			t.Fatalf("Unable to parse test url, %q: %s", test.fileurl, err)
		}
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
		url := test.fileurl
		err = p.Install(test.fileurl)
		if err != nil {
			t.Fatal(err)
		}
		if actualUrl != url {
			t.Fatalf("Expected %s, got %s", url, actualUrl)
		}

		files, err := ioutil.ReadDir(installDir)
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range files {
			if f.Name() != test.filename {
				t.Fatalf("Found unexpected file in install dir: %q", f.Name())
			}
			expectedMode := fs.FileMode(0755)
			if f.Mode() != expectedMode {
				t.Fatalf("Expected %v, got %v", expectedMode, f.Mode())
			}
		}

		expectedFilename := path.Join(installDir, test.filename)

		actualContent, err := os.ReadFile(expectedFilename)
		if err != nil {
			t.Fatal(err)
		}
		if string(actualContent) != string(content) {
			t.Fatalf("Expected %s, got %s", content, actualContent)
		}
	}
}

// cdHeader returns Headers with given Content-Disposition values.
func cdHeader(values []string) *http.Header {
	// values := make([]string, 0, len(filenames))
	// for _, filename := range filenames {
	// 	// value := "attachment"
	// 	// if filename != "" {
	// 	// 	value = fmt.Sprintf("%s; filename=%q", value, filename)
	// 	// }
	// 	values = append(values, value)
	// }
	return &http.Header{http.CanonicalHeaderKey("content-disposition"): values}
}
