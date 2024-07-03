package trie_test

import (
	"strings"
	"testing"

	trie "github.com/jangala-dev/tinygo-trie"
)

func TestInitialization(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))
	if tr == nil {
		t.Fatal("Trie initialization failed with default parameters")
	}

	customSeparatorTrie := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"), trie.WithSeparator("/"))
	if customSeparatorTrie == nil {
		t.Fatal("Trie initialization failed with custom parameters")
	}
}

func TestInsertion(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))

	if ok, _ := tr.Insert("abcd", "value1"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}
	if ok, _ := tr.Insert("abcf", "value2"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}

	customSeparatorTrie := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"), trie.WithSeparator("/"))

	if ok, _ := customSeparatorTrie.Insert("a/b/c/d", "value3"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}
	if ok, _ := customSeparatorTrie.Insert("a/b/c/f", "value4"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}

	// Test wildcard insertion
	if ok, err := tr.Insert("a#b", "value5"); ok || err == nil {
		t.Fatal("Trie accepted wildcard insertion wrongly")
	}

	if ok, _ := customSeparatorTrie.Insert("a/b/#+", "value6"); !ok {
		t.Fatal("Trie failed to insert a valid multi-level wildcard")
	}
}

func TestRetrieval(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))
	tr.Insert("abcd", "value1")

	val, _ := tr.Retrieve("abcd")
	if val != "value1" {
		t.Fatal("Failed to retrieve value for key")
	}

	val, _ = tr.Retrieve("abc")
	if val != nil {
		t.Fatal("Retrieved value for non-existent key")
	}

	// Test wildcards in retrieval
	tr.Insert("ab+f", "value2")
	val, _ = tr.Retrieve("ab+f")
	if val != "value2" {
		t.Fatal("Failed to retrieve value for key with wildcard")
	}
}

func TestMatching(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"), trie.WithSeparator("/"))

	tr.Insert("a/b/c/d", "value1")
	tr.Insert("a/b/c/f", "value2")
	tr.Insert("a/b/d/#", "value3")

	matches := tr.Match("a/b/c/d")
	if len(matches) != 1 {
		t.Fatal("Incorrect number of matches returned")
	}

	matches = tr.Match("a/b/c/+")
	if len(matches) != 2 {
		t.Fatal("Incorrect number of matches for single wildcard")
	}

	matches = tr.Match("a/b/+/d")
	if len(matches) != 2 {
		t.Fatal("Incorrect number of matches for mix of wildcards and keys")
	}

	matches = tr.Match("a/b/+/#")
	if len(matches) != 3 {
		t.Fatal("Incorrect number of matches for overlapping key patterns")
	}
}

func TestDeletion(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))
	tr.Insert("abcd", "value1")
	tr.Insert("abc", "value2")

	if !tr.Delete("abcd") {
		t.Fatal("Failed to delete key-value pair - part I")
	}
	if val, _ := tr.Retrieve("abcd"); val != nil {
		t.Fatal("Failed to delete key-value pair")
	}

	if tr.Delete("abcd") {
		t.Fatal("Incorrectly deleted non-existent key")
	}

	if val, _ := tr.Retrieve("abc"); val != "value2" {
		t.Fatal("Deleted key that's a prefix of another key")
	}
}

func TestDeleteLeafNode(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))
	tr.Insert("abcd", "value1")
	tr.Insert("abcde", "value2")

	if !tr.Delete("abcde") {
		t.Fatal("Failed to delete leaf node")
	}
	if val, _ := tr.Retrieve("abcde"); val != nil {
		t.Fatal("Failed to delete value of leaf node")
	}
	if val, _ := tr.Retrieve("abcd"); val != "value1" {
		t.Fatal("Affected other keys while deleting leaf node")
	}
}

func TestDeleteInternalNode(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))
	tr.Insert("abcd", "value1")
	tr.Insert("abcde", "value2")

	if !tr.Delete("abcd") {
		t.Fatal("Failed to delete internal node")
	}
	if val, _ := tr.Retrieve("abcd"); val != nil {
		t.Fatal("Failed to delete value of internal node")
	}
	if val, _ := tr.Retrieve("abcde"); val != "value2" {
		t.Fatal("Affected other keys while deleting internal node")
	}
}

func TestDeleteWithWildcards(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))
	tr.Insert("ab+cd", "value1")
	tr.Insert("ab+ce", "value2")

	if !tr.Delete("ab+cd") {
		t.Fatal("Failed to delete key with wildcard")
	}
	if val, _ := tr.Retrieve("ab+cd"); val != nil {
		t.Fatal("Failed to delete value of key with wildcard")
	}
	if val, _ := tr.Retrieve("ab+ce"); val != "value2" {
		t.Fatal("Affected other keys with wildcards while deleting")
	}
}

func TestChainDeletion(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))
	tr.Insert("a", "value1")
	tr.Insert("ab", "value2")
	tr.Insert("abc", "value3")

	if !tr.Delete("abc") {
		t.Fatal("Failed to delete leaf node in chain")
	}
	if !tr.Delete("ab") {
		t.Fatal("Failed to delete internal node in chain")
	}
	if val, _ := tr.Retrieve("a"); val != "value1" {
		t.Fatal("Affected root of the chain while deleting chained nodes")
	}
	if val, _ := tr.Retrieve("ab"); val != nil {
		t.Fatal("Failed to chain delete correctly")
	}
	if val, _ := tr.Retrieve("abc"); val != nil {
		t.Fatal("Failed to chain delete correctly")
	}
}

func TestDeleteNonExistentKey(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))
	tr.Insert("abcd", "value1")

	if tr.Delete("abcde") {
		t.Fatal("Incorrectly deleted non-existent key")
	}
}

func TestEdgeCases(t *testing.T) {
	tr := trie.New()

	if ok, _ := tr.Insert("", "empty"); !ok {
		t.Fatal("Failed to insert empty key")
	}
	if val, _ := tr.Retrieve(""); val != "empty" {
		t.Fatal("Failed to retrieve empty key")
	}

	longKey := strings.Repeat("a", 1000)
	if ok, _ := tr.Insert(longKey, "value"); !ok {
		t.Fatal("Failed to insert long key")
	}
	if val, _ := tr.Retrieve(longKey); val != "value" {
		t.Fatal("Failed to retrieve long key")
	}

	specialKey := "a!@#$%^&*()-_=[]{}|;:',.<>?/~"
	if ok, _ := tr.Insert(specialKey, "value"); !ok {
		t.Fatal("Failed to insert key with special characters")
	}
	if val, _ := tr.Retrieve(specialKey); val != "value" {
		t.Fatal("Failed to retrieve key with special characters")
	}

	newTrie := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))

	if ok, err := newTrie.Insert("a#b", "value"); ok || err == nil {
		t.Fatal("Accepted multi-level wildcard as single-level wildcard")
	}
}

func TestOverlappingKeysWithWildcards(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"), trie.WithSeparator("/"))

	if ok, _ := tr.Insert("a/b/c/d", "value1"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}
	if ok, _ := tr.Insert("a/+/c/d", "value2"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}
	if ok, _ := tr.Insert("a/b/c/#", "value3"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}
	if ok, _ := tr.Insert("a/+/+/#", "value4"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}

	matches := tr.Match("a/b/c/d")
	if len(matches) != 4 {
		t.Fatal("Incorrect number of matches for overlapping keys with wildcards")
	}
}

func TestCustomSeparatorEdgeCases(t *testing.T) {
	tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"), trie.WithSeparator("*"))

	if ok, _ := tr.Insert("a*b*c*d", "value1"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}
	matches := tr.Match("a*+*c*d")
	if len(matches) != 1 {
		t.Fatal("Incorrect number of matches with custom separator")
	}

	if ok, _ := tr.Insert("a*b*+*#", "value2"); !ok {
		t.Fatal("Failed to insert key-value pair")
	}
	matches = tr.Match("a*b*+*d")
	if len(matches) != 2 {
		t.Fatal("Incorrect number of matches with custom separator and wildcards")
	}
}
