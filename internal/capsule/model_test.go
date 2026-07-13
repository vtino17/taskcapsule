package capsule

import "testing"

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"payment-timeout", false},
		{"GH-482", false},
		{"checkout_redesign", false},
		{"a", false},
		{"", true},
		{".", true},
		{"..", true},
		{"../../etc", true},
		{"/task", true},
		{"task\\name", true},
		{"a-b-c-d-e-f-g-h-i-j-k-l-m-n-o-p-q-r-s-t-u-v-w-x-y-z-0123456789-1234567890-1234567890-1234567890-1234567890", true},
	}

	for _, tt := range tests {
		err := ValidateName(tt.name)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateName(%q) error = %v, wantErr = %v", tt.name, err, tt.wantErr)
		}
	}
}

func TestValidTransition(t *testing.T) {
	tests := []struct {
		from string
		to   string
		ok   bool
	}{
		{"", "preparing", true},
		{"preparing", "running", true},
		{"running", "pausing", true},
		{"pausing", "paused", true},
		{"paused", "resuming", true},
		{"resuming", "running", true},
		{"running", "error", true},
		{"deleting", "", true},
		{"", "running", false},
		{"running", "paused", false},
		{"paused", "running", false},
	}

	for _, tt := range tests {
		got := ValidTransition(tt.from, tt.to)
		if got != tt.ok {
			t.Errorf("ValidTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.ok)
		}
	}
}
