package main

import "fmt"

func ExampleToIntColumn() {
	names := []string{"-A", "A", "B", "Z", "AA", "AB", "AZ", "BA", "BZ", "ZZ", "AAC", "-AAC"}
	for _, name := range names {
		col, _ := toIntColumn(name)
		fmt.Printf("%s:%d\n", name, col)
	}

	// Output:
	// -A:-1
	// A:1
	// B:2
	// Z:26
	// AA:27
	// AB:28
	// AZ:52
	// BA:53
	// BZ:78
	// ZZ:702
	// AAC:705
	// -AAC:-705
}
