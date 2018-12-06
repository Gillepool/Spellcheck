package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func buildWordLexicon() map[string]int {
	lines, err := readLines("big.txt")
	if err != nil {
		panic(err)
	}

	words := WordCount(lines)
	return words
}

func sum(input []int) int {
	sum := 0

	for i := range input {
		sum += input[i]
	}
	return sum
}

func WordCount(ss []string) map[string]int {
	isAlpha := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	wordmap := make(map[string]int)
	for _, s := range ss {
		words := strings.Split(s, " ")
		for i := 0; i < len(words); i++ {
			if isAlpha(strings.ToLower(words[i])) {
				wordmap[strings.ToLower(words[i])]++
			}
		}
	}
	return wordmap
}

func probability(word string, words map[string]int) float64 {
	m := make([]int, 0, len(words))
	for _, val := range words {
		m = append(m, val)
	}
	N := sum(m)
	return float64(words[word]) / float64(N)
}

type Pair struct {
	left, right string
}

func edits1(word string) []string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	splits := make([]Pair, 0)
	for i := 0; i < len(word)+1; i++ {
		splits = append(splits, Pair{word[:i], word[i:]})
	}

	deletes := make([]string, 0)
	for _, s := range splits {
		if len(s.right) > 0 {
			deletes = append(deletes, s.left+s.right[1:])
		}
	}

	transposes := make([]string, 0)
	for _, s := range splits {
		if len(s.right) > 1 {
			transposes = append(transposes, s.left+string(([]rune(s.right)[1]))+string(([]rune(s.right)[0]))+s.right[2:])
		}
	}

	replaces := make([]string, 0)
	for _, s := range splits {
		for _, c := range letters {
			if len(s.right) > 0 {
				replaces = append(replaces, s.left+string(c)+s.right[1:])

			}
		}

	}

	inserts := make([]string, 0)
	for _, s := range splits {
		for _, c := range letters {
			inserts = append(inserts, s.left+string(c)+s.right)
		}
	}

	final_slice := make([]string, 0)
	final_slice = append(final_slice, deletes...)
	final_slice = append(final_slice, transposes...)
	final_slice = append(final_slice, replaces...)
	final_slice = append(final_slice, inserts...)

	final_slice = removeDuplicates(final_slice)
	return final_slice
}

func edits2(word string) []string {
	final_slice := make([]string, 0)
	for _, e1 := range edits1(word) {
		for _, e2 := range edits1(e1) {
			final_slice = append(final_slice, e2)
		}
	}
	return final_slice
}

func candidates(word string, wordMap map[string]int) []string {
	wordSlice := []string{word}
	a := known(wordSlice, wordMap)
	b := known(edits1(word), wordMap)
	c := known(edits2(word), wordMap)
	d := wordSlice

	if len(a) > 0 {
		return a
	} else if len(b) > 0 {
		return b
	} else if len(c) > 0 {
		return c
	} else if len(d) > 0 {
		return d
	}
	return nil
}

func known(words []string, wordMap map[string]int) []string {
	final_slice := make([]string, 0)
	for _, w := range words {
		if _, ok := wordMap[w]; ok {
			final_slice = append(final_slice, w)
		}
	}

	final_slice = removeDuplicates(final_slice)
	return final_slice

}

func correction(word string, wordMap map[string]int) string {
	var finalWord = ""
	var highestProb = 0.0
	for _, w := range candidates(word, wordMap) {
		if p := probability(w, wordMap); p > highestProb {
			finalWord = w
		}
	}
	return finalWord
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func main() {
	words := buildWordLexicon()
	fmt.Println(correction("brd", words))
}
