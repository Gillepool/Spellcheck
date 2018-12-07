package spellcheck

import (
	"bufio"
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

// BuildWordLexicon builds a lexicon based on the words in the txt file.
func BuildWordLexicon() map[string]int {
	lines, err := readLines("big.txt")
	if err != nil {
		panic(err)
	}
	words := wordCount(lines)
	return words
}

func sum(input []int) int {
	sum := 0
	for i := range input {
		sum += input[i]
	}
	return sum
}

func wordCount(ss []string) map[string]int {
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

type pair struct {
	left, right string
}

func edits1(word string) []string {
	letters := "abcdefghijklmnopqrstuvwxyz"
	splits := make([]pair, 0)
	for i := 0; i < len(word)+1; i++ {
		splits = append(splits, pair{word[:i], word[i:]})
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

	finalSlice := make([]string, 0)
	finalSlice = append(finalSlice, deletes...)
	finalSlice = append(finalSlice, transposes...)
	finalSlice = append(finalSlice, replaces...)
	finalSlice = append(finalSlice, inserts...)

	finalSlice = removeDuplicates(finalSlice)
	return finalSlice
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

func edits2(word string) []string {
	finalSlice := make([]string, 0)
	for _, e1 := range edits1(word) {
		for _, e2 := range edits1(e1) {
			finalSlice = append(finalSlice, e2)
		}
	}
	return finalSlice
}

// Candidates returns possible candidates of the misspellt word
func Candidates(word string, wordMap map[string]int) []string {
	wordSlice := []string{word}
	a := Known(wordSlice, wordMap)
	b := Known(edits1(word), wordMap)
	c := Known(edits2(word), wordMap)
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

// Known returns a slice of known words
func Known(words []string, wordMap map[string]int) []string {
	finalSlice := make([]string, 0)
	for _, w := range words {
		if _, ok := wordMap[w]; ok && !contains(finalSlice, w) {
			finalSlice = append(finalSlice, w)
		}
	}
	return finalSlice

}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Correction returns the word with the highest probability of being the right word
func Correction(word string, wordMap map[string]int) (finalWord string) {
	var highestProb = 0.0
	for _, w := range Candidates(word, wordMap) {
		if p := probability(w, wordMap); p > highestProb {
			finalWord = w
		}
	}
	return
}
