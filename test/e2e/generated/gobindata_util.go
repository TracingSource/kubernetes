package generated

import "k8s.io/klog"

/*
ReadOrDie reads a file from gobindata.
Relies heavily on the successful generation of build artifacts as per the go:generate directives above.
*/
func ReadOrDie(filePath string) []byte {

	fileBytes, err := Asset(filePath)
	if err != nil {
		gobindataMsg := "An error occurred, possibly gobindata doesn't know about the file you're opening. For questions on maintaining gobindata, contact the sig-testing group."
		klog.Infof("Available gobindata files: %v ", AssetNames())
		klog.Fatalf("Failed opening %v , with error %v.  %v.", filePath, err, gobindataMsg)
	}
	return fileBytes
}
