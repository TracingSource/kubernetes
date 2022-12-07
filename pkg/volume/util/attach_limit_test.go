package util

import (
	"crypto/sha1"
	"encoding/hex"
	"testing"

	"k8s.io/api/core/v1"
	v1helper "k8s.io/kubernetes/pkg/apis/core/v1/helper"
)

func TestGetCSIAttachLimitKey(t *testing.T) {
	// When driverName is less than 39 characters
	csiLimitKey := GetCSIAttachLimitKey("com.amazon.ebs")
	if csiLimitKey != "attachable-volumes-csi-com.amazon.ebs" {
		t.Errorf("Expected com.amazon.ebs got %s", csiLimitKey)
	}

	// When driver is longer than 39 chars
	longDriverName := "com.amazon.kubernetes.eks.ec2.ebs/csi-driver"
	csiLimitKeyLonger := GetCSIAttachLimitKey(longDriverName)
	if !v1helper.IsAttachableVolumeResourceName(v1.ResourceName(csiLimitKeyLonger)) {
		t.Errorf("Expected %s to have attachable prefix", csiLimitKeyLonger)
	}

	expectedCSIKey := getDriverHash(longDriverName)
	if csiLimitKeyLonger != expectedCSIKey {
		t.Errorf("Expected limit to be %s got %s", expectedCSIKey, csiLimitKeyLonger)
	}
}

func getDriverHash(driverName string) string {
	charsFromDriverName := driverName[:23]
	hash := sha1.New()
	hash.Write([]byte(driverName))
	hashed := hex.EncodeToString(hash.Sum(nil))
	hashed = hashed[:16]
	return CSIAttachLimitPrefix + charsFromDriverName + hashed
}
