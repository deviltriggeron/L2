package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"unicode"
)

func containsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func countNumberToStr(str string, count int) string {
	res := ""
	if str != "\\" {
		for i := 1; i < count; i++ {
			res += str
		}
	} else {
		res += strconv.Itoa(count)
	}
	return res
}

func transcript(str string) (string, error) {
	if !containsLetter(str) {
		return "", errors.New("there are only numbers in the line")
	}

	if !unicode.IsLetter(rune(str[0])) {
		return "", errors.New("the first letter should not be a number")
	}

	res := ""
	for i := range str {
		j, err := strconv.Atoi(string(str[i]))
		if err == nil {
			tmp := countNumberToStr(string(str[i-1]), j)
			res += tmp
		} else if str[i] == '\\' {
			continue
		} else {
			res += string(str[i])
		}
	}
	return res, nil
}

func main() {
	str := `qwe4asd\3`
	res, err := transcript(str)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("res: %v\n", res)
}
