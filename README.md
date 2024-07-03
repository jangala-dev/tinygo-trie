# Trie Library README
This library provides a Trie implementation with support for single-level (+) and multi-level (#) wildcards. Wildcards are permitted in both Trie entries (with Insert()) and in Trie searches (with Match()). Tries are commonly used for looking up keys efficiently and can be extended to support wildcard queries. The library is designed to be flexible with the delimiter used to split the keys.

## Key Features:
* Insertion of key-value pairs into the Trie.
* Retrieval of values using exact key matches, treating wildcards as literals.
* Matching with wildcards to retrieve multiple values.
* Deletion of keys from the Trie, prunes unneeded nodes.

## Usage:
### 1. Initialization:
```
package main

func main() {
    // Use wildcard characters '+' and '#'
    tr := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"))

    // Use custom separator
    customSeparatorTrie := trie.New(trie.WithSingleWild("+"), trie.WithMultiWild("#"), trie.WithSeparator("/"))
}
```

### 2. Insertion:
Insert a key-value pair into the Trie.

```
tr.Insert("abcd", "value1")
customSeparatorTrie.Insert("a/b/c/d", "value3")
```

Wildcards can be used in entries:
* The multi-level wildcard (#) can only be used at the end of a key.
* The single-level wildcard (+) can be used at any level.
go

```
tr.Insert("ab+f", "value2")  // Acceptable usage of single-level wildcard
tr.Insert("ab+", "value3")   // Another acceptable usage
_, err := tr.Insert("a#b", "value5")  // Will error out due to incorrect wildcard usage
if err != nil {
    fmt.Println("Error inserting key:", err)
}
```

### 3. Retrieval:
Retrieve a value from the Trie using an exact key.

```
value, _ := tr.Retrieve("abcd")  // Returns "value1"
```

### 4. Matching:
Retrieve values matching a given pattern with potential wildcards.

```
matches := tr.Match("ab+f")  // Will match entries like "abcf", "abdf", etc.
matches = customSeparatorTrie.Match("a/b/+/#")  -- Returns matches for keys like "a/b/c/d" 
```
* Use the single-level wildcard (+) to match any single level of the key.
* Use the multi-level wildcard (#) to match zero or more levels at the end of the key.
  
### 5. Deletion:
Delete a key from the Trie.
```
tr.Delete("abcd")
```

### 6. Edge Cases:
The library can handle a variety of edge cases, such as empty keys, very long keys, and keys with special characters. Tests have been written to cover these scenarios.

## Wildcards:
* Single-Level Wildcard (+): Matches any single part of a key. For example, "ab+c" could match "abc" or "abcd", but not "abbc" or "abcc".
* Multi-Level Wildcard (#): Matches zero or more parts of a key at the end. For example, "ab#" could match "ab", "abc", "abcd", etc.

When using Match(), both wildcards can be part of the match key. This allows for complex queries to retrieve multiple entries that match the given pattern.

## Testing
To run the tests for this library, ensure that you have both trie.go and trie_test.go in the same directory with the correct package declaration. Then, run the tests using the go test command:

```
go test
```

This will automatically find and run all the tests in files named *_test.go in the current directory.