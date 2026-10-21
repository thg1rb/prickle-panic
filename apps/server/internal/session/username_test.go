package session

import "testing"

func TestNormalizeUsername(t *testing.T) {
	cases := []struct {
		name, in, want string
		valid          bool
	}{
		{"english", "  Coco  ", "Coco", true},
		{"thai marks", "กิข้", "กิข้", true},
		{"thai vowel and tone order", "กิ้ข", "กิ้ข", true},
		{"thai digits", "๑๒", "๑๒", true},
		{"ascii digits", "A1", "A1", true},
		{"one internal space", "Co Co", "Co Co", true},
		{"NFC equivalent character", "KB", "KB", true},
		{"trimmed whitespace", "\u00a0AB\u00a0", "AB", true},
		{"repeated spaces", "A  B", "", false},
		{"tab is rejected before trim", "\tAB", "", false},
		{"bidi", "AB\u202e", "", false},
		{"zero width", "A\u200bB", "", false},
		{"emoji", "A😀", "", false},
		{"unsupported script", "AЖ", "", false},
		{"markup", "<b>A</b>", "", false},
		{"standalone Thai mark", "\u0e34A", "", false},
		{"repeated malformed Thai mark", "กิิข", "", false},
		{"malformed Thai mark order", "ก้ิข", "", false},
		{"single grapheme", "กิ", "", false},
		{"two Thai graphemes", "กิข้", "กิข้", true},
		{"16 graphemes", "abcdefghijklmnop", "abcdefghijklmnop", true},
		{"17 graphemes", "abcdefghijklmnopq", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NormalizeUsername(tc.in)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v err=%v", tc.valid, err)
			}
			if err == nil && got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestNormalizeUsernameRawByteLimitAndUTF8(t *testing.T) {
	if _, err := NormalizeUsername("\xffAB"); err == nil {
		t.Fatal("expected invalid UTF-8 rejection")
	}
	if _, err := NormalizeUsername("AB" + repeat("x", MaxUsernameBytes-1)); err == nil {
		t.Fatal("expected raw byte limit rejection")
	}
	if _, err := NormalizeUsername("AB" + repeat("x", MaxUsernameBytes-2)); err == nil {
		t.Fatal("256-byte input with excessive grapheme count must be rejected")
	}
	if got, err := NormalizeUsername("abcdefghijklmno"); err != nil || got != "abcdefghijklmno" {
		t.Fatalf("valid grapheme boundary rejected: %v", err)
	}
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}
