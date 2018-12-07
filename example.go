package main

import (
	"fmt"

	"github.com/gillepool/Spellcheck/spellcheck"
)

func main() {
	words := spellcheck.BuildWordLexicon()
	fmt.Println(spellcheck.Correction("helly", words))
}
