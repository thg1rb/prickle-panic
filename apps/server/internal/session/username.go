package session

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
	"golang.org/x/text/unicode/norm"
)

var ErrUsernameInvalid = errors.New("username is invalid")

const MaxUsernameBytes = 256

// NormalizeUsername applies the approved Thai-first username allowlist and
// returns NFC display text. It never interprets the value as markup.
func NormalizeUsername(raw string) (string, error) {
	if !utf8.ValidString(raw) || len(raw) > MaxUsernameBytes {
		return "", ErrUsernameInvalid
	}
	for _, r := range raw {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) || isBidiControl(r) {
			return "", ErrUsernameInvalid
		}
	}
	value := norm.NFC.String(strings.TrimFunc(raw, unicode.IsSpace))
	if value == "" {
		return "", ErrUsernameInvalid
	}
	lastThaiBase := false
	thaiMarks := 0
	thaiVowelMark := false
	thaiFinalMark := false
	lastSpace := false
	for _, r := range value {
		if r == ' ' {
			if lastSpace {
				return "", ErrUsernameInvalid
			}
			lastSpace = true
			lastThaiBase = false
			continue
		}
		lastSpace = false
		switch {
		case r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z':
			lastThaiBase = false
		case r >= '0' && r <= '9' || r >= '\u0e50' && r <= '\u0e59':
			lastThaiBase = false
		case unicode.IsLetter(r) && unicode.In(r, unicode.Thai):
			lastThaiBase = true
			thaiMarks = 0
			thaiVowelMark = false
			thaiFinalMark = false
		case unicode.IsMark(r) && unicode.In(r, unicode.Thai) && lastThaiBase:
			// Thai vowel and tone marks may extend a Thai base or its mark sequence.
			if thaiMarks >= 2 {
				return "", ErrUsernameInvalid
			}
			if r >= '\u0e48' && r <= '\u0e4e' {
				if thaiFinalMark {
					return "", ErrUsernameInvalid
				}
				thaiFinalMark = true
			} else {
				if thaiVowelMark || thaiFinalMark {
					return "", ErrUsernameInvalid
				}
				thaiVowelMark = true
			}
			thaiMarks++
		case isBidiControl(r) || unicode.IsControl(r) || unicode.Is(unicode.Cf, r):
			return "", ErrUsernameInvalid
		default:
			return "", ErrUsernameInvalid
		}
	}
	if lastSpace {
		return "", ErrUsernameInvalid
	}
	count := 0
	for g := graphemes.FromString(value); g.Next(); {
		count++
	}
	if count < 2 || count > 16 {
		return "", ErrUsernameInvalid
	}
	return value, nil
}

func isBidiControl(r rune) bool {
	return r == '\u061c' || r == '\u200e' || r == '\u200f' ||
		r >= '\u202a' && r <= '\u202e' || r >= '\u2066' && r <= '\u2069'
}
