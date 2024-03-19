package utils

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var batchSize = 1000

// ReadStrings get []string from file
func ReadStrings(file *os.File) ([]string, error) {
	var result []string
	reader := bufio.NewReader(file)
	for {
		content, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				if content != "" {
					content = strings.TrimRight(content, "\n")
					result = append(result, content)
				}
				break
				// 注意此处还要处理最后一行
			} else {
				return result, errors.New(err.Error())
			}
		}
		if content == "" {
			continue
		}
		content = strings.TrimRight(content, "\n")
		content = strings.TrimRight(content, "\r")
		result = append(result, content)
	}
	return result, nil
}

func ReadString(filePath string) (string, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	var result = ""
	s, err := ReadStrings(file)
	if err != nil {
		return result, err
	}
	for _, v := range s {
		result += v
	}
	return result, nil
}
