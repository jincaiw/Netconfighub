package service

import (
	"fmt"
	"strings"
)

func validateStorageName(label, value string) error {
	name := strings.TrimSpace(value)
	if name == "" {
		return fmt.Errorf("%s不能为空", label)
	}
	if name == "." || name == ".." {
		return fmt.Errorf("%s不能为 . 或 ..", label)
	}
	for _, r := range name {
		if r < 0x20 || r == 0x7f {
			return fmt.Errorf("%s不能包含控制字符", label)
		}
		switch r {
		case '/', '\\', ':':
			return fmt.Errorf("%s不能包含路径分隔符或冒号", label)
		}
	}
	return nil
}
