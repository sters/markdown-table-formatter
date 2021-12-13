package markdowntableformatter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitLine(t *testing.T) {
	result := splitLine(`
	# header

	|head1|head2|
	|-----|-----------|
	|fiiiiiiiit|ho|
	|a|bbb               |
	`)

	want := []string{
		"# header",
		"",
		"|head1|head2|",
		"|-----|-----------|",
		"|fiiiiiiiit|ho|",
		"|a|bbb               |",
	}

	assert.Len(t, result, 6)

	for i, r := range result {
		assert.Equal(t, want[i], r)
	}
}

func TestSplitColumn(t *testing.T) {
	var result []string

	result = splitColumn("# header")
	assert.Len(t, result, 1)
	assert.Equal(t, "# header", result[0])

	result = splitColumn("|head1|head2|")
	assert.Len(t, result, 2)
	assert.Equal(t, "head1", result[0])
	assert.Equal(t, "head2", result[1])

	result = splitColumn("|head1               |           head2|")
	assert.Len(t, result, 2)
	assert.Equal(t, "head1", result[0])
	assert.Equal(t, "head2", result[1])
}

func TestCheckTable(t *testing.T) {
	assert.False(t, checkTable("# header"))
	assert.False(t, checkTable("[github](https://github.com/)"))
	assert.False(t, checkTable("| hogehoge"))
	assert.True(t, checkTable("|head1|head2|"))
	assert.True(t, checkTable("|-----|"))
}

func TestGetSeparateTable(t *testing.T) {
	result := getSeparateTable(`
	|head1|head2|
	|-----|-----------|
	|fiiiiiiiit|ho|
	|a|bbb               |
	`)

	want := [][]string{
		{"head1", "head2"},
		{"-----", "-----------"},
		{"fiiiiiiiit", "ho"},
		{"a", "bbb"},
	}

	assert.Len(t, result, 4)

	for lineIdx, line := range result {
		for columnIdx, column := range line {
			assert.Equal(t, want[lineIdx][columnIdx], column)
		}
	}
}

func TestGetHasSeparatorLine(t *testing.T) {
	var target [][]string

	target = [][]string{
		{"head1", "head2"},
		{"-----", "-----------"},
		{"fiiiiiiiit", "ho"},
		{"a", "bbb"},
	}
	assert.True(t, getHasSeparatorLine(target))

	target = [][]string{
		{"head1", "head2"},
		{"fiiiiiiiit", "ho"},
		{"a", "bbb"},
	}
	assert.False(t, getHasSeparatorLine(target))
}

func TestGetMaxLength(t *testing.T) {
	var target [][]string
	var maxLength []int
	var want []int

	target = [][]string{
		{"head1", "head2"},
		{"----------------------", "------"},
		{"fiiiiiiiit", "ho"},
		{"a", "bbb"},
	}
	maxLength = getMaxLength(target)

	want = []int{10, 5} // fiiiiiiiit => 10, head2 => 5

	assert.Len(t, maxLength, 2)
	for i, length := range maxLength {
		assert.Equal(t, want[i], length)
	}
}

func TestGetMaxLength2(t *testing.T) {
	var target [][]string
	var maxLength []int
	var want []int

	target = [][]string{
		{"head1", "head2"},
		{"----------------------", "------"},
		{"こんにちは", "ho"},
		{"a", "bbb"},
	}
	maxLength = getMaxLength(target)

	want = []int{10, 5} // こんにちは => 10, head2 => 5

	assert.Len(t, maxLength, 2)
	for i, length := range maxLength {
		assert.Equal(t, want[i], length)
	}
}

func TestFixColumnSize(t *testing.T) {
	target := ""
	target += "|head1|head2|" + "\n"
	target += "|-----|-----------|" + "\n"
	target += "|fiiiiiiiit|ho|aaaa|" + "\n"
	target += "|a|bbb               |" + "\n"
	result := fixColumnSize(target)

	want := ""
	want += "|head1     |head2|    |" + "\n"
	want += "|----------|-----|----|" + "\n"
	want += "|fiiiiiiiit|ho   |aaaa|" + "\n"
	want += "|a         |bbb  |    |" + "\n"

	assert.Equal(t, want, result)
}

func TestFixColumnSize2(t *testing.T) {
	target := ""
	target += "|head1|head2|" + "\n"
	target += "|-----|-----------|" + "\n"
	target += "|こんにちは|ho|aaaa|" + "\n"
	target += "|a|bbb               |" + "\n"
	result := fixColumnSize(target)

	want := ""
	want += "|head1     |head2|    |" + "\n"
	want += "|----------|-----|----|" + "\n"
	want += "|こんにちは|ho   |aaaa|" + "\n"
	want += "|a         |bbb  |    |" + "\n"

	assert.Equal(t, want, result)
}

func TestFindTables(t *testing.T) {
	resultsList := findTables(`
		# header
		[Github](https://github.com/)

		|head1|head2|
		|-----|-----------|
		|fiiiiiiiit|ho|aaaa|
		|a|bbb               |

		|head1|head2|
		|a|bbb               |

		// notable
		|head1|head2|
	`)

	wantsList := [][]int{
		{3, 6},
		{8, 9},
	}

	assert.Len(t, resultsList, 2)

	for i, results := range resultsList {
		assert.Len(t, results, 2)

		for j, result := range results {
			assert.Equal(t, wantsList[i][j], result)
		}
	}
}

func TestFindTables2(t *testing.T) {
	resultsList := findTables(`
	# Example file

	|hoge            |huga|
	|--------------|------|
	|its|confusing              |markdown|
	|table   |so|crazy|table|.|

	Try to use: [this](https://github.com/sters/markdown-table-formatter)
	`)

	wantsList := [][]int{
		{2, 5},
	}

	assert.Len(t, resultsList, 1)

	for i, results := range resultsList {
		assert.Len(t, results, 2)

		for j, result := range results {
			assert.Equal(t, wantsList[i][j], result)
		}
	}
}

func TestFindTables3(t *testing.T) {
	resultsList := findTables(`
		# header
		[Github](https://github.com/)

		|head1|head2|
		|-----|-----------|
		|fiiiiiiiit|ho|aaaa|
		|a|bbb               |
	`)

	wantsList := [][]int{
		{3, 6},
	}

	assert.Len(t, resultsList, 1)

	for i, results := range resultsList {
		assert.Len(t, results, 2)

		for j, result := range results {
			assert.Equal(t, wantsList[i][j], result)
		}
	}
}

func TestExtractTables(t *testing.T) {
	result := extractTables(`
		|head1|head2|
		|-----|-----------|
		|fiiiiiiiit|ho|aaaa|
		|a|bbb               |

		|head1|head2|
		|a|bbb               |

		// notable
		|head1|head2|
	`)

	want := []string{
		`
		|head1|head2|
		|-----|-----------|
		|fiiiiiiiit|ho|aaaa|
		|a|bbb               |
		`,
		`
		|head1|head2|
		|a|bbb               |
		`,
	}

	assert.Len(t, result, 2)

	for i, table := range result {
		resultLines := strings.Split(table, "\n")
		wantLines := strings.Split(strings.TrimSpace(want[i]), "\n")

		for j, resultLint := range resultLines {
			assert.Equal(t, strings.TrimSpace(wantLines[j]), strings.TrimSpace(resultLint))
		}
	}
}

func TestExecute(t *testing.T) {
	result := Execute(`
	# header
	[Github](https://github.com/)

	|head1|head2|
	|-----|-----------|
	|fiiiiiiiit|ho|aaaa|
	|a|bbb               |

	|head1|head2|
	|a|bbb               |

	|ヘッダー1|ヘッダーツー|
	|こんにちは|Hello, 世界               |

	// notable
	|head1|head2|
	`)

	want := `
	# header
	[Github](https://github.com/)

	|head1     |head2|    |
	|----------|-----|----|
	|fiiiiiiiit|ho   |aaaa|
	|a         |bbb  |    |

	|head1|head2|
	|a    |bbb  |

	|ヘッダー1 |ヘッダーツー|
	|こんにちは|Hello, 世界 |

	// notable
	|head1|head2|
	`

	resultSplit := strings.Split(strings.TrimSpace(result), "\n")
	wantSplit := strings.Split(strings.TrimSpace(want), "\n")

	assert.Len(t, resultSplit, len(wantSplit))

	for idx, resultLine := range resultSplit {
		a := strings.TrimSpace(resultLine)
		b := strings.TrimSpace(wantSplit[idx])
		assert.Equal(t, b, a)
	}
}

func TestExecute2(t *testing.T) {
	result := Execute(`
	# Example file

	|hoge            |huga|
	|--------------|------|
	|its|confusing              |markdown|
	|table   |so|crazy|table|.|

	Try to use: [this](https://github.com/sters/markdown-table-formatter)
	`)

	want := `
	# Example file

	|hoge |huga     |        |     | |
	|-----|---------|--------|-----|-|
	|its  |confusing|markdown|     | |
	|table|so       |crazy   |table|.|

	Try to use: [this](https://github.com/sters/markdown-table-formatter)
	`

	resultSplit := strings.Split(strings.TrimSpace(result), "\n")
	wantSplit := strings.Split(strings.TrimSpace(want), "\n")

	assert.Len(t, resultSplit, len(wantSplit))

	for idx, resultLine := range resultSplit {
		a := strings.TrimSpace(resultLine)
		b := strings.TrimSpace(wantSplit[idx])
		assert.Equal(t, b, a)
	}
}
