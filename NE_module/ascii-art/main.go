package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . \"text\"")
		return
	}

	input := os.Args[1]

	data, err := os.ReadFile("standard.txt")
	if err != nil {
		panic("Error: banner file (standard.txt) is missing")
	}

	lines := splitdata(string(data))
	asciimap := createmap(lines)
	input = strings.ReplaceAll(input, "\\n", "\n")

	output := buildASCII(input, asciimap)

	for _, line := range output {
		fmt.Println(line)
	}
}

func splitdata(data string) []string {
	return strings.Split(strings.ReplaceAll(data, "\r\n", "\n"), "\n")
}

func createmap(lines []string) map[rune][]string {
	chrnum := 32
	asciimap := make(map[rune][]string)

	for i := 0; i < len(lines); i += 9 {
		if i+8 > len(lines) {
			break
		}
		asciimap[rune(chrnum)] = lines[i+1 : i+8+1]
		chrnum++
	}

	return asciimap
}

func buildASCII(input string, asciimap map[rune][]string) []string {
	output := make([]string, 8)
	flag := true
	for _, ch := range input {
		if flag {
			if ch == '\n' {
				for i := 0; i < 8; i++ {
					fmt.Println(output[i])
					output[i] = ""
					flag = false
				}
				continue
			}
		}

		if val, ok := asciimap[ch]; ok {
			for i := 0; i < 8; i++ {
				output[i] += val[i]
			}
		}
	}

	return output
}
