package markdowntableformatter

import (
	"fmt"
	"regexp"
	"strings"
)

func splitLine(text string) []string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	for idx, line := range lines {
		lines[idx] = strings.TrimSpace(line)
	}

	return lines
}

func splitColumn(text string) []string {
	columns := strings.Split(strings.Trim(text, "|"), "|")
	for idx, line := range columns {
		columns[idx] = strings.TrimSpace(line)
	}

	return columns
}

func checkTable(text string) bool {
	return regexp.MustCompile(`^\|.+\|`).MatchString(text)
}

func getSeparateTable(text string) [][]string {
	var lines [][]string

	for _, line := range splitLine(text) {
		lines = append(lines, splitColumn(line))
	}

	return lines
}

func getHasSeparatorLine(lines [][]string) bool {
	has := false
	if len(lines) > 2 && len(lines[0]) > 0 {
		has = regexp.MustCompile(`^\-+$`).MatchString(lines[1][0])
	}

	return has
}

func calculateLength(str string) int {
	length := 0
	for _, c := range []rune(str) {
		if len(string(c)) > 2 {
			length += 2
		} else {
			length++
		}
	}
	return length
}

func getMaxLength(lines [][]string) []int {
	var maxLength []int
	hasSeparatorLine := getHasSeparatorLine(lines)

	for lineIdx, line := range lines {
		if hasSeparatorLine && lineIdx == 1 {
			continue
		}

		for columnIdx, column := range line {
			length := calculateLength(column)
			if len(maxLength) <= columnIdx {
				maxLength = append(maxLength, length)
			} else if maxLength[columnIdx] < length {
				maxLength[columnIdx] = length
			}
		}
	}

	return maxLength
}

func fixColumnSize(text string) string {
	// split
	lines := getSeparateTable(text)

	// check has separator line
	hasSeparatorLine := getHasSeparatorLine(lines)

	// find max lenght per columns
	maxLength := getMaxLength(lines)

	// fix length
	for lineIdx, line := range lines {
		for columnIdx, length := range maxLength {
			if len(line) <= columnIdx {
				lines[lineIdx] = append(lines[lineIdx], "")
				line = lines[lineIdx]
			}

			column := line[columnIdx]

			if hasSeparatorLine && lineIdx == 1 {
				lines[lineIdx][columnIdx] = strings.Repeat("-", length)
			} else {
				lines[lineIdx][columnIdx] = fmt.Sprintf(
					"%s%s",
					column,
					strings.Repeat(" ", length-calculateLength(column)),
				)
			}
		}
	}

	// concat
	result := ""
	for _, line := range lines {
		result += "|" + strings.Join(line, "|") + "|\n"
	}

	return result
}

func findTables(text string) [][]int {
	lines := splitLine(text)
	results := [][]int{}
	startLine := -1

	var (
		lineIdx int
		line    string
	)

	for lineIdx, line = range lines {
		if checkTable(line) {
			if startLine == -1 {
				startLine = lineIdx
			}
		} else {
			if startLine != -1 && startLine+1 != lineIdx {
				results = append(results, []int{startLine, lineIdx - 1})
			}

			startLine = -1
		}
	}

	if startLine != -1 && startLine+1 < lineIdx {
		results = append(results, []int{startLine, lineIdx})
	}

	return results
}

func extractTables(text string) []string {
	lines := splitLine(text)
	tablePositions := findTables(text)

	var results []string
	for _, position := range tablePositions {
		start := position[0]
		end := position[1] + 1
		results = append(
			results,
			strings.Join(lines[start:end], "\n"),
		)
	}

	return results
}

func Execute(text string) string {
	var fixedTable []string
	for _, table := range extractTables(text) {
		fixedTable = append(fixedTable, fixColumnSize(table))
	}

	lines := splitLine(text)
	for idx, position := range findTables(text) {
		for i, fixedLine := range splitLine(fixedTable[idx]) {
			lines[position[0]+i] = fixedLine
		}
	}

	return strings.Join(lines, "\n")
}
