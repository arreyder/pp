package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
)

// Detect and beautify JSON within a stream of text
func beautifyJSONInStream(text string) {
	// Regex to find JSON-like structures (objects and arrays)
	jsonRegex := regexp.MustCompile(`\{.*\}|\[.*\]`)
	matches := jsonRegex.FindAllStringIndex(text, -1)
	lastIndex := 0

	// Create color settings
	jsonColor := color.New(color.FgCyan).SprintFunc()
	nonJSONColor := color.New(color.FgWhite).SprintFunc()

	for _, match := range matches {
		// Print the non-JSON part with original color
		fmt.Print(nonJSONColor(text[lastIndex:match[0]]))

		// Extract the potential JSON part
		jsonPart := text[match[0]:match[1]]

		// Try to beautify the matched JSON
		var parsedJSON interface{}
		if err := json.Unmarshal([]byte(jsonPart), &parsedJSON); err == nil {
			// If it's valid JSON, print the beautified version with color
			beautifiedJSON, _ := json.MarshalIndent(parsedJSON, "", "    ")
			fmt.Println("\n" + jsonColor(string(beautifiedJSON)))
		} else {
			// If it's not valid JSON, print it as-is with original color
			fmt.Print(nonJSONColor(jsonPart))
		}

		// Update the lastIndex to the end of the current match
		lastIndex = match[1]
	}

	// Print any remaining text after the last JSON block
	fmt.Print(nonJSONColor(text[lastIndex:]) + "\n")
}

func main() {
	// Create a scanner to read input line by line from stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Buffer to accumulate text between JSON parts
	var buffer bytes.Buffer

	// Read each line from stdin and process it
	for scanner.Scan() {
		line := scanner.Text()
		buffer.WriteString(line)

		// Try to process JSON within the buffer if the line might contain JSON
		beautifyJSONInStream(buffer.String())

		// Clear the buffer after processing
		buffer.Reset()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
