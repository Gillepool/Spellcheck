package spellcheck

import (
	"testing"
)

func TestWordLexicon(t *testing.T) {
	words := BuildWordLexicon()
	if len(words) != 25310 {
		t.Error("Failed")
	}

}

func TestCandidates(t *testing.T) {
	words := BuildWordLexicon()
	corrections := Correction("foreig", words)
	if corrections != "foreign" {
		t.Error("Failed")
	}
}
