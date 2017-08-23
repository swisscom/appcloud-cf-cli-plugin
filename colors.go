package main

import "github.com/fatih/color"

// Text formatting and colors for print text
var (
	bold      = color.New(color.Bold).SprintFunc()
	greenBold = color.New(color.FgGreen).Add(color.Bold).SprintFunc()
	cyanBold  = color.New(color.FgCyan).Add(color.Bold).SprintFunc()
)
