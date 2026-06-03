package service

import "testing"

func TestValidateStorageNameAcceptsCommonNames(t *testing.T) {
	valid := []string{
		"core-sw-01",
		"核心交换机 01",
		"edge.switch_1",
	}
	for _, name := range valid {
		if err := validateStorageName("设备名称", name); err != nil {
			t.Fatalf("validateStorageName(%q) returned error: %v", name, err)
		}
	}
}

func TestValidateStorageNameRejectsPathLikeNames(t *testing.T) {
	invalid := []string{
		"",
		"   ",
		".",
		"..",
		"../outside",
		"group/device",
		"group\\device",
		"group:device",
		"device\nname",
	}
	for _, name := range invalid {
		if err := validateStorageName("设备名称", name); err == nil {
			t.Fatalf("validateStorageName(%q) expected error", name)
		}
	}
}
