package docker

import "testing"

func TestImageTagMatches(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want bool
	}{
		{"autograder", []string{"autograder:latest"}, true},
		{"autograder", []string{"other:latest"}, false},
		{"autograder", []string{"autograder:v2", "other:latest"}, true},
		{"autograder", []string{}, false},
		{"autograder", nil, false},
	}
	for _, tt := range tests {
		if got := imageTagMatches(tt.name, tt.tags); got != tt.want {
			t.Errorf("imageTagMatches(%q, %v) = %v, want %v", tt.name, tt.tags, got, tt.want)
		}
	}
}
