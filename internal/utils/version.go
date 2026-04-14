package utils

import "fmt"

// Convert version information to a string
func ParseVersion(v uint32) string {
	// 根据实际测试，版本格式为：
	// 0x06000001 = 6.0.1
	// 所以：major = v >> 24, minor = (v >> 16) & 0xFF, patch = v & 0xFFFF
	major := v >> 24
	minor := (v >> 16) & 0xFF
	patch := v & 0xFFFF
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
