package container

import (
	"encoding/json"
	"testing"

	"k8s.io/api/core/v1"
)

var (
	sampleContainer = `
{
  "name": "test_container",
  "image": "foo/image:v1",
  "command": [
    "/bin/testcmd"
  ],
  "args": [
    "/bin/sh",
    "-c",
    "echo abc"
  ],
  "ports": [
    {
      "containerPort": 8001
    }
  ],
  "env": [
    {
      "name": "ENV_FOO",
      "value": "bar"
    },
    {
      "name": "ENV_BAR",
      "valueFrom": {
        "secretKeyRef": {
          "name": "foo",
          "key": "bar",
          "optional": true
        }
      }
    }
  ],
  "resources": {
    "limits": {
      "foo": "1G"
    },
    "requests": {
      "foo": "500M"
    }
  }
}
`

	sampleV115HashValue = uint64(0x311670a)
	sampleV116HashValue = sampleV115HashValue
)

func TestConsistentHashContainer(t *testing.T) {
	container := &v1.Container{}
	if err := json.Unmarshal([]byte(sampleContainer), container); err != nil {
		t.Error(err)
	}

	currentHash := HashContainer(container)
	if currentHash != sampleV116HashValue {
		t.Errorf("mismatched hash value with v1.16")
	}

	if currentHash != sampleV115HashValue {
		t.Errorf("mismatched hash value with v1.15")
	}
}
