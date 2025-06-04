package utils

import (
	"fmt"
	"os"
	"regexp"
)

func GetRegexInFile(path string, pattern string, replacement string, printError bool) (string, error) {
	// Get regex from file.
	read, err := os.ReadFile(path)
	if err != nil {
		if printError {
			fmt.Println(err)
		}
		return "", err
	}

	// Find regex and only output based on the pattern specified.
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(regex.FindString(string(read)), replacement), nil
}

func Find(pattern, replacement, input string) string {
	regex := regexp.MustCompile(pattern)
	return regex.ReplaceAllString(regex.FindString(input), replacement)
}

func Replace(pattern, replacement string, input string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}
	return re.ReplaceAllString(input, replacement), nil
}

func ReplaceAll(replacements map[string]string, input string) (string, error) {
	result := input
	for pattern, replacement := range replacements {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return "", err
		}
		result = re.ReplaceAllString(result, replacement)
	}
	return result, nil
}
