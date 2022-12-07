// +build !providerless

package aws

import "testing"

func TestDeviceAllocator(t *testing.T) {
	tests := []struct {
		name            string
		existingDevices ExistingDevices
		deviceMap       map[mountDevice]int
		expectedOutput  mountDevice
	}{
		{
			"empty device list with wrap",
			ExistingDevices{},
			generateUnsortedDeviceList(),
			"bd", // next to 'zz' is the first one, 'ba'
		},
	}

	for _, test := range tests {
		allocator := NewDeviceAllocator().(*deviceAllocator)
		for k, v := range test.deviceMap {
			allocator.possibleDevices[k] = v
		}

		got, err := allocator.GetNext(test.existingDevices)
		if err != nil {
			t.Errorf("text %q: unexpected error: %v", test.name, err)
		}
		if got != test.expectedOutput {
			t.Errorf("text %q: expected %q, got %q", test.name, test.expectedOutput, got)
		}
	}
}

func generateUnsortedDeviceList() map[mountDevice]int {
	possibleDevices := make(map[mountDevice]int)
	for _, firstChar := range []rune{'b', 'c'} {
		for i := 'a'; i <= 'z'; i++ {
			dev := mountDevice([]rune{firstChar, i})
			possibleDevices[dev] = 3
		}
	}
	possibleDevices["bd"] = 0
	return possibleDevices
}

func TestDeviceAllocatorError(t *testing.T) {
	allocator := NewDeviceAllocator().(*deviceAllocator)
	existingDevices := ExistingDevices{}

	// make all devices used
	var first, second byte
	for first = 'b'; first <= 'c'; first++ {
		for second = 'a'; second <= 'z'; second++ {
			device := [2]byte{first, second}
			existingDevices[mountDevice(device[:])] = "used"
		}
	}

	device, err := allocator.GetNext(existingDevices)
	if err == nil {
		t.Errorf("expected error, got device  %q", device)
	}
}
