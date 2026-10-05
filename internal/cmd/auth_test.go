package cmd

import "testing"

func TestParsePastedCode(t *testing.T) {
	cases := []struct {
		name, input, want string
		wantErr           bool
	}{
		{"redirect URL", "http://localhost:18080/callback?code=abc123&state=s1", "abc123", false},
		{"bare code", "abc123", "abc123", false},
		{"no state", "http://localhost:18080/callback?code=abc123", "abc123", false},
		{"state mismatch", "http://localhost:18080/callback?code=abc123&state=other", "", true},
		{"denied", "http://localhost:18080/callback?error=access_denied&state=s1", "", true},
		{"no code", "http://localhost:18080/callback?state=s1", "", true},
	}
	for _, c := range cases {
		got, err := parsePastedCode(c.input, "s1")
		if (err != nil) != c.wantErr || got != c.want {
			t.Errorf("%s: got (%q, %v), want (%q, err=%v)", c.name, got, err, c.want, c.wantErr)
		}
	}
}
