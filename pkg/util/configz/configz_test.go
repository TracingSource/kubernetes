package configz

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConfigz(t *testing.T) {
	v, err := New("testing")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	v.Set("blah")

	s := httptest.NewServer(http.HandlerFunc(handle))
	defer s.Close()

	resp, err := http.Get(s.URL + "/configz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(body) != `{"testing":"blah"}` {
		t.Fatalf("unexpected output: %s", body)
	}

	v.Set("bing")
	resp, err = http.Get(s.URL + "/configz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(body) != `{"testing":"bing"}` {
		t.Fatalf("unexpected output: %s", body)
	}

	Delete("testing")
	resp, err = http.Get(s.URL + "/configz")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if string(body) != `{}` {
		t.Fatalf("unexpected output: %s", body)
	}
	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("unexpected Content-Type: %s", resp.Header.Get("Content-Type"))
	}
}
