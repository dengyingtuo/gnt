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

func ExampleToLetterColumn() {
	cols := []int{-1, 1, 2, 26, 27, 28, 52, 53, 78, 702, 705, -705}
	for _, col := range cols {
		name := toLetterColumn(col)
		fmt.Printf("%d:%s\n", col, name)
	}

	// Output:
	// -1:-A
	// 1:A
	// 2:B
	// 26:Z
	// 27:AA
	// 28:AB
	// 52:AZ
	// 53:BA
	// 78:BZ
	// 702:ZZ
	// 705:AAC
	// -705:-AAC
}
func ExampleConvFunc() {
	s := "2001:220,2002:320"
	item := &Item{
		ColsConv: map[string]string{"A": "split2 , : id chance"},
	}
	conv := item.GetConvFunc("A")
	if conv != nil {
		vals, keys := conv(s)
		fmt.Println(vals)
		fmt.Println(keys)
	}

	// Output:
	// [[2001 220] [2002 320]]
	// [id chance]
}
