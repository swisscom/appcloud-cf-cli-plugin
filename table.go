package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Table is a table in the console.
type Table interface {
	Add(row ...string)
	Print()
}

// PrintableTable is a printable table in the console.
type PrintableTable struct {
	headers       []string
	headerPrinted bool
	maxSizes      []int
	rows          [][]string
}

// NewTable creates a new table.
func NewTable(headers []string) Table {
	return &PrintableTable{
		headers:  headers,
		maxSizes: make([]int, len(headers)),
	}
}

// Add adds a new row to a table.
func (t *PrintableTable) Add(row ...string) {
	t.rows = append(t.rows, row)
}

// Print prints a table to stdout.
func (t *PrintableTable) Print() {
	for _, row := range append(t.rows, t.headers) {
		t.calculateMaxSize(row)
	}

	if t.headerPrinted == false {
		t.printHeader()
		t.headerPrinted = true
	}

	for _, line := range t.rows {
		t.printRow(line)
	}

	t.rows = [][]string{}
}

func (t *PrintableTable) calculateMaxSize(row []string) {
	for index, value := range row {
		cellLength := utf8.RuneCountInString(value)
		if t.maxSizes[index] < cellLength {
			t.maxSizes[index] = cellLength
		}
	}
}

func (t *PrintableTable) printHeader() {
	output := ""
	for col, value := range t.headers {
		output = output + t.cellValue(col, value)
	}
	fmt.Println(bold(output))
}

func (t *PrintableTable) printRow(row []string) {
	output := ""
	for columnIndex, value := range row {
		output = output + t.cellValue(columnIndex, value)
	}
	fmt.Println(output)
}

func (t *PrintableTable) cellValue(col int, value string) string {
	padding := ""
	if col < len(t.headers)-1 {
		padding = strings.Repeat(" ", t.maxSizes[col]-utf8.RuneCountInString(value))
	}
	return fmt.Sprintf("%s%s   ", value, padding)
}
