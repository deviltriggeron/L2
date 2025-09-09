package main

import (
	"fmt"
	"sort"
	"strings"
)

func sortString(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

func anagrams(words []string) map[string][]string {
	valueMap := make(map[string][]string)
	keysMap := make(map[string]string)

	for i := range words {
		wordLower := strings.ToLower(words[i])
		sorted := sortString(wordLower)

		valueMap[sorted] = append(valueMap[sorted], wordLower)

		if _, exists := keysMap[sorted]; !exists {
			keysMap[sorted] = wordLower
		}

	}

	result := make(map[string][]string)
	for key, group := range valueMap {
		if len(group) > 1 {
			sort.Strings(group)
			result[keysMap[key]] = group
		}
	}

	return result
}

func main() {
	words := []string{"пятак", "листок", "пятка", "слиток", "тяпка", "столик", "стол"}
	anagrams := anagrams(words)

	for k, v := range anagrams {
		fmt.Printf("%s: %s\n", k, v)
	}
}
