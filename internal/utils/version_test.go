package utils

import "testing"

func TestParseVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		src    uint32
		expect string
	}{
		{"1.0.0", 0x01000000, "1.0.0"},
		{"2.3.4", 0x02030004, "2.3.4"},
		{"5.0.0", 0x05000000, "5.0.0"},
		{"6.0.1", 0x06000001, "6.0.1"},
		{"0.0.0", 0, "0.0.0"},
		{"255.255.65535", 0xFFFF_FFFF, "255.255.65535"},
		{"3.1.0", 0x03010000, "3.1.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ParseVersion(tt.src)
			if got != tt.expect {
				t.Errorf("ParseVersion(0x%08X) = %q, want %q", tt.src, got, tt.expect)
			}
		})
	}
}
