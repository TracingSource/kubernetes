package windows

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// downloadFile saves a remote URL to a local temp file, and returns its path.
// It's the caller's responsibility to clean up the temp file when done.
func downloadFile(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "unable to download from %q", url)
	}
	defer response.Body.Close()

	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		return "", errors.Wrapf(err, "unable to create temp file")
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, response.Body)
	return tempFile.Name(), err
}
