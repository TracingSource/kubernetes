package flexvolume

import (
	"testing"
)

func TestHandleResponseDefaults(t *testing.T) {
	ds, err := handleCmdResponse("test", []byte(`{"status": "Success"}`))
	if err != nil {
		t.Error("error: ", err)
	}

	if *ds.Capabilities != *defaultCapabilities() {
		t.Error("wrong default capabilities: ", *ds.Capabilities)
	}
}
