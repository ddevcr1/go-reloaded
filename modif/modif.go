package modifier

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func Process(input string) string {
	tokens := tokenize(input)
	tokens = applyCommands(tokens)
	tokens = fixArticles(tokens)
	result := joinTokens(tokens)
	result = fixPunctuation(result)
	result = fixQuotes(result)
	result = fixDoubleQuotes(result)
	return result
}

// Разбивка на токены
func tokenize(input string) []string {
	reApos := regexp.MustCompile(`([A-Za-z0-9])'([A-Za-z0-9])`)
	input = reApos.ReplaceAllString(input, "$1<APOS>$2")

	input = strings.ReplaceAll(input, ")(", ") (")

	re := regexp.MustCompile(`([().,!?:;'"'])`)
	input = re.ReplaceAllString(input, " $1 ")

	input = strings.ReplaceAll(input, "<APOS>", "'")

	return strings.Fields(input)
}

func isCommand(tokens []string, i int) (string, int, int) {
	if i >= len(tokens) || tokens[i] != "(" {
		return "", 0, 0
	}
	if i+2 >= len(tokens) {
		return "", 0, 0
	}

	rawCmd := tokens[i+1]
	t := strings.ToLower(rawCmd)

	switch t {
	case "hex", "bin", "up", "low", "cap":
		if tokens[i+2] == ")" {
			return t, 3, 1
		}
		if i+4 < len(tokens) && tokens[i+2] == "," && tokens[i+4] == ")" {
			n, err := strconv.Atoi(tokens[i+3])
			if err == nil {
				return t, 5, n
			}
		}
	}
	return "", 0, 0
}

// Применение команд к тексту
func applyCommands(tokens []string) []string {
	for {
		changed := false
		var result []string
		i := 0
		for i < len(tokens) {
			cmd, consumed, count := isCommand(tokens, i)
			if cmd != "" {
				switch cmd {
				case "hex":
					if len(result) > 0 {
						result[len(result)-1] = hexToDecimal(result[len(result)-1])
					}
				case "bin":
					if len(result) > 0 {
						result[len(result)-1] = binToDecimal(result[len(result)-1])
					}
				case "up":
					applyToLastN(result, count, strings.ToUpper)
				case "low":
					applyToLastN(result, count, strings.ToLower)
				case "cap":
					applyToLastN(result, count, capitalize)
				}
				i += consumed
				changed = true
			} else {
				result = append(result, tokens[i])
				i++
			}
		}
		tokens = result
		if !changed {
			break
		}
	}
	return tokens
}

func applyToLastN(tokens []string, n int, fn func(string) string) {
	if n <= 0 {
		return
	}
	start := len(tokens) - n
	if start < 0 {
		start = 0
	}
	for j := start; j < len(tokens); j++ {
		tokens[j] = fn(tokens[j])
	}
}

func hexToDecimal(s string) string {
	val, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return s
	}
	return fmt.Sprintf("%d", val)
}

func binToDecimal(s string) string {
	val, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return s
	}
	return fmt.Sprintf("%d", val)
}

// Первая буква заглавная, остальные строчные
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	firstIdx := -1
	for i, r := range runes {
		if unicode.IsLetter(r) {
			firstIdx = i
			break
		}
	}

	if firstIdx >= 0 {
		runes[firstIdx] = unicode.ToUpper(runes[firstIdx])
		for i := firstIdx + 1; i < len(runes); i++ {
			if unicode.IsLetter(runes[i]) {
				runes[i] = unicode.ToLower(runes[i])
			}
		}
		return string(runes)
	}
	return s
}

// Исправление артиклей
func fixArticles(tokens []string) []string {
	for i := 0; i < len(tokens)-1; i++ {
		t := tokens[i]
		cleanT := strings.Trim(t, "'")
		lower := strings.ToLower(cleanT)

		isA := (lower == "a")
		isAn := (lower == "an")

		if isA || isAn {
			idx := i + 1
			for idx < len(tokens) {
				if isPunctuationToken(tokens[idx]) {
					idx++
					continue
				}

				cleanNext := strings.Trim(tokens[idx], "'")
				cleanNext = strings.TrimLeft(cleanNext, "(")

				if len(cleanNext) > 0 {
					firstChar := rune(cleanNext[0])
					if unicode.IsLetter(firstChar) || unicode.IsNumber(firstChar) {
						isVowel := isVowelOrH(firstChar)
						if isA && isVowel {
							tokens[i] = replaceArticle(t, "an")
						} else if isAn && !isVowel {
							tokens[i] = replaceArticle(t, "a")
						}
					}
				}
				break
			}
		}
	}
	return tokens
}

// Проверка на пунктуацию
func isPunctuationToken(s string) bool {
	return strings.ContainsAny(s, ".,!?:;()\"'")
}

// Замена артикля
func replaceArticle(original, newBase string) string {
	prefix := ""
	suffix := ""
	content := original
	if strings.HasPrefix(content, "'") {
		prefix = "'"
		content = content[1:]
	}
	if strings.HasSuffix(content, "'") {
		suffix = "'"
		content = content[:len(content)-1]
	}

	out := newBase
	if len(content) > 0 && unicode.IsUpper(rune(content[0])) {
		out = capitalize(newBase)
	}
	return prefix + out + suffix
}

// Проверка на гласную
func isVowelOrH(r rune) bool {
	switch unicode.ToLower(r) {
	case 'a', 'e', 'i', 'o', 'u', 'h':
		return true
	}
	return false
}

func joinTokens(tokens []string) string {
	return strings.Join(tokens, " ")
}

func fixPunctuation(text string) string {
	tokens := strings.Fields(text)
	var result []string
	for _, t := range tokens {
		if isSuffixPunctuation(t) {
			if len(result) > 0 {
				result[len(result)-1] += t
			} else {
				result = append(result, t)
			}
		} else {
			if len(result) > 0 && isPrefixPunctuation(result[len(result)-1]) {
				result[len(result)-1] += t
			} else {
				result = append(result, t)
			}
		}
	}
	return strings.Join(result, " ")
}

func isSuffixPunctuation(s string) bool {
	return strings.ContainsAny(s, ".,!?:;)")
}
func isPrefixPunctuation(s string) bool {
	if len(s) == 0 {
		return false
	}
	return s[len(s)-1] == '('
}

func fixQuotes(text string) string {
	return fixGenericQuotes(text, "'")
}
func fixDoubleQuotes(text string) string {
	return fixGenericQuotes(text, "\"")
}

// Логика кавычек
func fixGenericQuotes(text string, quoteChar string) string {
	tokens := strings.Fields(text)
	var result []string

	open := false

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		if t == quoteChar {
			if !open {
				if i+1 < len(tokens) {
					result = append(result, quoteChar+tokens[i+1])
					i++
					open = true
				} else {
					if len(result) > 0 {
						result[len(result)-1] += t
					} else {
						result = append(result, t)
					}
				}
			} else {
				if len(result) > 0 {
					result[len(result)-1] += quoteChar
				} else {
					result = append(result, t)
				}
				open = false
			}
		} else {
			result = append(result, t)
		}
	}

	return strings.Join(result, " ")
}
