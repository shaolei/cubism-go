package core

import "testing"

func TestParseMajorVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version string
		want    int
		wantErr bool
	}{
		{"standard 5.0.0", "5.0.0", 5, false},
		{"standard 6.0.1", "6.0.1", 6, false},
		{"single digit major", "4.1", 4, false},
		{"large major", "12.0.0", 12, false},
		{"single number", "5", 5, false},
		{"non-numeric major", "abc.0.0", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseMajorVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseMajorVersion(%q) error = %v, wantErr %v", tt.version, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseMajorVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}
