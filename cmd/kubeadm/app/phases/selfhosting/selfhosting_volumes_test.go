package selfhosting

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func createTemporaryFile(name string) *os.File {
	content := []byte("foo")
	tmpfile, err := ioutil.TempFile("", name)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	return tmpfile
}

func TestCreateTLSSecretFromFile(t *testing.T) {
	tmpCert := createTemporaryFile("foo.crt")
	defer os.Remove(tmpCert.Name())
	tmpKey := createTemporaryFile("foo.key")
	defer os.Remove(tmpKey.Name())

	_, err := createTLSSecretFromFiles("foo", tmpCert.Name(), tmpKey.Name())
	if err != nil {
		log.Fatal(err)
	}

	if err := tmpCert.Close(); err != nil {
		log.Fatal(err)
	}

	if err := tmpKey.Close(); err != nil {
		log.Fatal(err)
	}
}

func TestCreateOpaqueSecretFromFile(t *testing.T) {
	tmpFile := createTemporaryFile("foo")
	defer os.Remove(tmpFile.Name())

	_, err := createOpaqueSecretFromFile("foo", tmpFile.Name())
	if err != nil {
		log.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
}
