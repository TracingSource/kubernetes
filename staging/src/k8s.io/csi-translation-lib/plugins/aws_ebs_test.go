package plugins

import (
	"reflect"
	"testing"

	storage "k8s.io/api/storage/v1"
)

func TestKubernetesVolumeIDToEBSVolumeID(t *testing.T) {
	testCases := []struct {
		name         string
		kubernetesID string
		ebsVolumeID  string
		expErr       bool
	}{
		{
			name:         "Normal ID format",
			kubernetesID: "vol-02399794d890f9375",
			ebsVolumeID:  "vol-02399794d890f9375",
		},
		{
			name:         "aws:///{volumeId} format",
			kubernetesID: "aws:///vol-02399794d890f9375",
			ebsVolumeID:  "vol-02399794d890f9375",
		},
		{
			name:         "aws://{zone}/{volumeId} format",
			kubernetesID: "aws://us-west-2a/vol-02399794d890f9375",
			ebsVolumeID:  "vol-02399794d890f9375",
		},
		{
			name:         "fails on invalid volume ID",
			kubernetesID: "aws://us-west-2a/02399794d890f9375",
			expErr:       true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := KubernetesVolumeIDToEBSVolumeID(tc.kubernetesID)
			if err != nil {
				if !tc.expErr {
					t.Errorf("KubernetesVolumeIDToEBSVolumeID failed %v", err)
				}
			} else {
				if actual != tc.ebsVolumeID {
					t.Errorf("Wrong EBS Volume ID. actual: %s expected: %s", actual, tc.ebsVolumeID)
				}
			}
		})
	}
}

func TestTranslateEBSInTreeStorageClassToCSI(t *testing.T) {
	translator := NewAWSElasticBlockStoreCSITranslator()

	cases := []struct {
		name   string
		sc     *storage.StorageClass
		expSc  *storage.StorageClass
		expErr bool
	}{
		{
			name:  "translate normal",
			sc:    NewStorageClass(map[string]string{"foo": "bar"}, nil),
			expSc: NewStorageClass(map[string]string{"foo": "bar"}, nil),
		},
		{
			name:  "translate empty map",
			sc:    NewStorageClass(map[string]string{}, nil),
			expSc: NewStorageClass(map[string]string{}, nil),
		},

		{
			name:  "translate with fstype",
			sc:    NewStorageClass(map[string]string{"fstype": "ext3"}, nil),
			expSc: NewStorageClass(map[string]string{"csi.storage.k8s.io/fstype": "ext3"}, nil),
		},
	}

	for _, tc := range cases {
		t.Logf("Testing %v", tc.name)
		got, err := translator.TranslateInTreeStorageClassToCSI(tc.sc)
		if err != nil && !tc.expErr {
			t.Errorf("Did not expect error but got: %v", err)
		}

		if err == nil && tc.expErr {
			t.Errorf("Expected error, but did not get one.")
		}

		if !reflect.DeepEqual(got, tc.expSc) {
			t.Errorf("Got parameters: %v, expected :%v", got, tc.expSc)
		}

	}
}
