package strutil

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func IntToBase62(n int) string {
	if n == 0 {
		return string(base62[0])
	}

	var result []byte
	for n > 0 {
		result = append(result, base62[n%62])
		n /= 62
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		return s
	}
	return output
}

func RemoveEmojis(input string) string {
	var output []rune
	for _, r := range input {
		if r <= unicode.MaxASCII {
			output = append(output, r)
		}
	}
	return string(output)
}

func Clean(input string) string {
	output := RemoveAccents(input)
	output = RemoveEmojis(output)
	return output
}

func ReplaceSlashes(input string, replace string) string {
	rule := []string{}
	rule = append(rule, "/", replace, "_", replace, "-", replace, "%", replace, "(", replace, ")", replace)
	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	return input
}

func RemoveSlashes(input string) string {
	return ReplaceSlashes(input, "")
}

func KebabCase(input string) string {
	rule := []string{}
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "-")
}

func WordCount(s string) int {
	var r rune
	var size, count int

	isWord := false

	for len(s) > 0 {
		r, size = utf8.DecodeRuneInString(s)

		switch {
		case isLetter(r):
			if !isWord {
				isWord = true
				count++
			}

		case isWord && (r == '\'' || r == '-'):
			// is word

		default:
			isWord = false
		}

		s = s[size:]
	}

	return count
}

func CamelCase(input string) string {
	caser := cases.Title(language.BrazilianPortuguese)
	rule := []string{}
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	return strings.Join(words, "")
}

func PascalCase(input string) string {
	caser := cases.Title(language.BrazilianPortuguese)
	rule := []string{}
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	return strings.Join(words, "")
}

func SnakeCase(input string) string {
	rule := []string{}
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.ToLower(word)
	}
	return strings.Join(words, "_")
}

func UpperSnakeCase(input string) string {
	rule := []string{}
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.ToUpper(word)
	}
	return strings.Join(words, "_")
}

func UpperFirst(s string) string {
	if len(s) == 0 {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	r = unicode.ToUpper(r)

	return string(r) + s[size:]
}

func LowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	r = unicode.ToLower(r)

	return string(r) + s[size:]
}

func UpperCamelCase(input string) string {
	caser := cases.Title(language.BrazilianPortuguese)
	rule := []string{}
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	return strings.Join(words, "")
}

func LowerCamelCase(input string) string {
	caser := cases.Title(language.English)
	rule := []string{}
	rule = append(rule, ".", " ", "_", " ", "-", " ")
	replacer := strings.NewReplacer(rule...)
	input = replacer.Replace(input)
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = caser.String(word)
	}
	return strings.Join(words, "")
}

// isLetter checks r is a letter but not CJK character.
func isLetter(r rune) bool {
	if !unicode.IsLetter(r) {
		return false
	}

	switch {
	// cjk char: /[\u3040-\u30ff\u3400-\u4dbf\u4e00-\u9fff\uf900-\ufaff\uff66-\uff9f]/

	// hiragana and katakana (Japanese only)
	case r >= '\u3034' && r < '\u30ff':
		return false

	// CJK unified ideographs extension A (Chinese, Japanese, and Korean)
	case r >= '\u3400' && r < '\u4dbf':
		return false

	// CJK unified ideographs (Chinese, Japanese, and Korean)
	case r >= '\u4e00' && r < '\u9fff':
		return false

	// CJK compatibility ideographs (Chinese, Japanese, and Korean)
	case r >= '\uf900' && r < '\ufaff':
		return false

	// half-width katakana (Japanese only)
	case r >= '\uff66' && r < '\uff9f':
		return false
	}

	return true
}
