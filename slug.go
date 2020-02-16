package slug

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks.
}

var invalidCharacters = regexp.MustCompile(`[^a-z0-9\s-_]`)
var separators = regexp.MustCompile(`[\s_-]+`)

func Slugify(text string) string {
	text = strings.ToLower(text)
	text = zhCharToPinyin(text)
	text = removeDiacritics(text)
	text = invalidCharacters.ReplaceAllString(text, "")
	text = separators.ReplaceAllString(text, " ") // Remove duplicate separators.
	text = strings.TrimSpace(text)
	text = separators.ReplaceAllString(text, "-") // Standardize separators.
	return text
}

var a = pinyin.NewArgs()

func zhCharToPinyin(p string) (s string) {
	for _, r := range p {
		if unicode.Is(unicode.Han, r) {
			s += string(pinyin.Pinyin(string(r), a)[0][0])
		} else {
			s += string(r)
		}
	}
	return
}

var t = transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

func removeDiacritics(text string) string {
	result, _, _ := transform.String(t, text)
	return result
}
