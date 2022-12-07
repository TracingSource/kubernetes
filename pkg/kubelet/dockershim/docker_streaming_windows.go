// +build windows

package dockershim

import (
	"bytes"
	"fmt"
	"io"

	"k8s.io/kubernetes/pkg/kubelet/util/ioutils"
)

func (r *streamingRuntime) portForward(podSandboxID string, port int32, stream io.ReadWriteCloser) error {
	stderr := new(bytes.Buffer)
	err := r.exec(podSandboxID, []string{"wincat.exe", "127.0.0.1", fmt.Sprint(port)}, stream, stream, ioutils.WriteCloserWrapper(stderr), false, nil, 0)
	if err != nil {
		return fmt.Errorf("%v: %s", err, stderr.String())
	}

	return nil
}
