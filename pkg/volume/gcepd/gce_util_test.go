// +build !providerless

package gcepd

import "testing"

func TestParseScsiSerial(t *testing.T) {
	cases := []struct {
		name      string
		output    string
		diskName  string
		expectErr bool
	}{
		{
			name:     "valid",
			output:   "0Google  PersistentDisk  test-disk",
			diskName: "test-disk",
		},
		{
			name:     "valid with newline",
			output:   "0Google  PersistentDisk  test-disk\n",
			diskName: "test-disk",
		},
		{
			name:      "invalid prefix",
			output:    "00Google  PersistentDisk  test-disk",
			expectErr: true,
		},
		{
			name:      "invalid suffix",
			output:    "0Google  PersistentDisk  test-disk  more",
			expectErr: true,
		},
	}

	for _, test := range cases {
		serial, err := parseScsiSerial(test.output)
		if err != nil && !test.expectErr {
			t.Errorf("test %v failed: %v", test.name, err)
		}
		if err == nil && test.expectErr {
			t.Errorf("test %q failed: got success", test.name)
		}
		if serial != test.diskName {
			t.Errorf("test %v failed: expected serial %q, got %q", test.name, test.diskName, serial)
		}
	}
}
