package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . <input> <output>")
		return
	}

	in := os.Args[1]
	out := os.Args[2]

	data, err := os.ReadFile(in)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	text := string(data)

	tokens := tokenize(text)
	tokens = processOps(tokens)
	tokens = fixAtoAn(tokens)
	tokens = fixPunctuation(tokens)
	tokens = fixQuotes(tokens)

	result := rebuild(tokens)

	os.WriteFile(out, []byte(result), 0644)
}

///////////////////////////////////////////////////////////
// 1) TOKENIZER
///////////////////////////////////////////////////////////

func tokenize(s string) []string {
	var tokens []string
	current := ""

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current += string(r)
		} else {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			if !unicode.IsSpace(r) {
				tokens = append(tokens, string(r))
			}
		}
	}

	if current != "" {
		tokens = append(tokens, current)
	}

	return tokens
}

///////////////////////////////////////////////////////////
// 2) PROCESS OPERATIONS
///////////////////////////////////////////////////////////

func processOps(tokens []string) []string {
	var out []string

	for i := 0; i < len(tokens); i++ {

		// (hex)
		if match(tokens, i, "(", "hex", ")") {
			n := out[len(out)-1]
			val, _ := strconv.ParseInt(n, 16, 64)
			out[len(out)-1] = strconv.Itoa(int(val))
			i += 2
			continue
		}

		// (bin)
		if match(tokens, i, "(", "bin", ")") {
			n := out[len(out)-1]
			val := binToDec(n)
			out[len(out)-1] = strconv.Itoa(val)
			i += 2
			continue
		}

		// (up)
		if match(tokens, i, "(", "up", ")") {
			out[len(out)-1] = strings.ToUpper(out[len(out)-1])
			i += 2
			continue
		}

		// (low)
		if match(tokens, i, "(", "low", ")") {
			out[len(out)-1] = strings.ToLower(out[len(out)-1])
			i += 2
			continue
		}

		// (cap)
		if match(tokens, i, "(", "cap", ")") {
			out[len(out)-1] = capitalize(out[len(out)-1])
			i += 2
			continue
		}

		// (up, n)
		if isRange(tokens, i, "up") {
			n, _ := strconv.Atoi(tokens[i+3])
			applyRange(out, n, strings.ToUpper)
			i += 4
			continue
		}

		// (low, n)
		if isRange(tokens, i, "low") {
			n, _ := strconv.Atoi(tokens[i+3])
			applyRange(out, n, strings.ToLower)
			i += 4
			continue
		}

		// (cap, n)
		if isRange(tokens, i, "cap") {
			n, _ := strconv.Atoi(tokens[i+3])
			applyRange(out, n, capitalize)
			i += 4
			continue
		}

		out = append(out, tokens[i])
	}

	return out
}

func match(t []string, i int, a, b, c string) bool {
	return i+2 < len(t) && t[i] == a && t[i+1] == b && t[i+2] == c
}

func isRange(t []string, i int, op string) bool {
	return i+4 < len(t) && t[i] == "(" && t[i+1] == op && t[i+2] == "," && t[i+4] == ")"
}

func binToDec(s string) int {
	res := 0
	for _, r := range s {
		res = res*2 + int(r-'0')
	}
	return res
}

func applyRange(words []string, n int, fn func(string) string) {
	for i := len(words) - 1; i >= 0 && n > 0; i-- {
		if isWord(words[i]) {
			words[i] = fn(words[i])
			n--
		}
	}
}

///////////////////////////////////////////////////////////
// 3) FIX QUOTES
///////////////////////////////////////////////////////////

func fixQuotes(tokens []string) []string {
	var out []string
	var buf []string
	in := false

	for _, t := range tokens {
		if t == "'" {
			if !in {
				in = true
				buf = []string{}
			} else {
				in = false
				out = append(out, "'"+strings.Join(buf, " ")+"'")
			}
			continue
		}

		if in {
			buf = append(buf, t)
		} else {
			out = append(out, t)
		}
	}

	return out
}

///////////////////////////////////////////////////////////
// 4) FIX PUNCTUATION
///////////////////////////////////////////////////////////

func fixPunctuation(tokens []string) []string {
	var out []string
	p := ".,!?:;"

	for _, t := range tokens {
		// multi punctuation (... !! !?)
		if len(t) > 1 && isMultiPunct(t) {
			out[len(out)-1] += t
			continue
		}

		// single punctuation
		if strings.ContainsRune(p, rune(t[0])) {
			out[len(out)-1] += t
			continue
		}

		out = append(out, t)
	}

	return out
}

func isMultiPunct(s string) bool {
	return s == "..." || s == "!!" || s == "!?" || s == "??"
}

///////////////////////////////////////////////////////////
// 5) FIX "a" → "an"
///////////////////////////////////////////////////////////

func fixAtoAn(tokens []string) []string {
	for i := 0; i < len(tokens)-1; i++ {

		if strings.ToLower(tokens[i]) == "a" &&
			startsWithVowel(tokens[i+1]) {

			tokens[i] = "an"
		}
	}
	return tokens
}

func startsWithVowel(s string) bool {
	if s == "" {
		return false
	}
	c := unicode.ToLower(rune(s[0]))
	return c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' || c == 'h'
}

///////////////////////////////////////////////////////////
// 6) UTILITIES
///////////////////////////////////////////////////////////

func rebuild(tokens []string) string {
	return strings.Join(tokens, " ")
}

func isWord(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
