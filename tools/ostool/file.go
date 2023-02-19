package ostool

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
)

func IsDir(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	if fi.IsDir() {
		return true, nil
	} else {
		return false, nil
	}
}

func FilterComments(filename string) []string {
	var linesWithoutComments []string
	// Open file for reading
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Read file line by line and filter out comments
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !isComment(line) {
			linesWithoutComments = append(linesWithoutComments, line)
			//fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return linesWithoutComments
}

func isComment(line string) bool {
	line = stripWhitespace(line)
	if len(line) == 0 {
		return false
	}

	if line[0] == '/' && len(line) > 1 {
		if line[1] == '/' || line[1] == '*' {
			return true
		}
	}
	return false
}

func stripWhitespace(line string) string {
	var result []rune
	for _, ch := range line {
		if ch != ' ' && ch != '\t' {
			result = append(result, ch)
		}
	}
	return string(result)
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		} else {
			log.Fatalf("Error stat file %s: %v\n", path, err)
		}
	}
	return true
}
