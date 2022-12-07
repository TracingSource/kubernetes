package glusterfs

import (
	"bufio"
	"fmt"
	"os"

	"k8s.io/klog"
)

// readGlusterLog will take the last 2 lines of the log file
// on failure of gluster SetUp and return those so kubelet can
// properly expose them
// return error on any failure
func readGlusterLog(path string, podName string) error {

	var line1 string
	var line2 string
	linecount := 0

	klog.Infof("failure, now attempting to read the gluster log for pod %s", podName)

	// Check and make sure path exists
	if len(path) == 0 {
		return fmt.Errorf("log file does not exist for pod %s", podName)
	}

	// open the log file
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open log file for pod %s", podName)
	}
	defer file.Close()

	// read in and scan the file using scanner
	// from stdlib
	fscan := bufio.NewScanner(file)

	// rather than guessing on bytes or using Seek
	// going to scan entire file and take the last two lines
	// generally the file should be small since it is pod specific
	for fscan.Scan() {
		if linecount > 0 {
			line1 = line2
		}
		line2 = "\n" + fscan.Text()

		linecount++
	}

	if linecount > 0 {
		return fmt.Errorf("%v", line1+line2+"\n")
	}
	return nil
}
