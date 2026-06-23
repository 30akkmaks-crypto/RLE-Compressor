// rle.go
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ANSI colors
const (
	reset  = "\033[0m"
	green  = "\033[92m"
	red    = "\033[91m"
	yellow = "\033[93m"
)

func colorize(text, color string) string {
	return color + text + reset
}

const ESCAPE = '\\'

func compress(text string) string {
	if text == "" {
		return ""
	}
	var result strings.Builder
	n := len(text)
	i := 0
	for i < n {
		ch := text[i]
		j := i + 1
		for j < n && text[j] == ch {
			j++
		}
		count := j - i
		if count >= 3 {
			result.WriteByte(ch)
			result.WriteByte(ESCAPE)
			result.WriteString(strconv.Itoa(count))
		} else {
			for k := 0; k < count; k++ {
				if ch == ESCAPE {
					result.WriteByte(ESCAPE)
				}
				result.WriteByte(ch)
			}
		}
		i = j
	}
	return result.String()
}

func decompress(text string) (string, error) {
	if text == "" {
		return "", nil
	}
	var result strings.Builder
	n := len(text)
	i := 0
	for i < n {
		ch := text[i]
		if ch == ESCAPE {
			// Проверяем, экранированный ли это escape
			if i+1 < n && text[i+1] == ESCAPE {
				result.WriteByte(ESCAPE)
				i += 2
				continue
			}
			// Управляющая последовательность
			if i+1 >= n {
				return "", errors.New("неожиданный конец строки после escape")
			}
			char := text[i+1]
			i += 2
			// Читаем число
			numStr := ""
			for i < n && text[i] >= '0' && text[i] <= '9' {
				numStr += string(text[i])
				i++
			}
			if numStr == "" {
				return "", errors.New("отсутствует число после escape")
			}
			count, err := strconv.Atoi(numStr)
			if err != nil {
				return "", err
			}
			for k := 0; k < count; k++ {
				result.WriteByte(char)
			}
		} else {
			result.WriteByte(ch)
			i++
		}
	}
	return result.String(), nil
}

func readInput(filename string) (string, error) {
	if filename == "-" || filename == "" {
		// читаем stdin
		reader := bufio.NewReader(os.Stdin)
		var builder strings.Builder
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return "", err
			}
			builder.WriteString(line)
			if err == io.EOF {
				break
			}
		}
		return builder.String(), nil
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeOutput(filename string, content string) error {
	if filename == "-" || filename == "" {
		fmt.Print(content)
		return nil
	}
	return os.WriteFile(filename, []byte(content), 0644)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(colorize("Usage: rle compress|decompress [input] [output]", yellow))
		os.Exit(1)
	}
	mode := os.Args[1]
	inputFile := ""
	outputFile := ""
	if len(os.Args) >= 3 {
		inputFile = os.Args[2]
	}
	if len(os.Args) >= 4 {
		outputFile = os.Args[3]
	}
	if mode != "compress" && mode != "decompress" {
		fmt.Println(colorize("Invalid mode. Use compress or decompress.", red))
		os.Exit(1)
	}

	// Чтение
	data, err := readInput(inputFile)
	if err != nil {
		fmt.Println(colorize("Error reading input: "+err.Error(), red))
		os.Exit(1)
	}

	var result string
	var err2 error
	if mode == "compress" {
		result = compress(data)
	} else {
		result, err2 = decompress(data)
		if err2 != nil {
			fmt.Println(colorize("Decompression error: "+err2.Error(), red))
			os.Exit(1)
		}
	}

	if err := writeOutput(outputFile, result); err != nil {
		fmt.Println(colorize("Error writing output: "+err.Error(), red))
		os.Exit(1)
	}
	if outputFile != "" && outputFile != "-" {
		fmt.Println(colorize("Result written to "+outputFile, green))
	}
}
