package markdowntableformatter

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/ToQoz/gopwt"
	"github.com/ToQoz/gopwt/assert"
)

func TestMain(m *testing.M) {
	flag.Parse()
	gopwt.Empower()
	os.Exit(m.Run())
}

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

	assert.Require(t, len(result) == 6)

	for i, r := range result {
		assert.Require(t, r == want[i])
	}
}

func TestSplitColumn(t *testing.T) {
	var result []string

	result = splitColumn("# header")
	assert.Require(t, len(result) == 1)
	assert.Require(t, result[0] == "# header")

	result = splitColumn("|head1|head2|")
	assert.Require(t, len(result) == 2)
	assert.Require(t, result[0] == "head1")
	assert.Require(t, result[1] == "head2")

	result = splitColumn("|head1               |           head2|")
	assert.Require(t, len(result) == 2)
	assert.Require(t, result[0] == "head1")
	assert.Require(t, result[1] == "head2")
}

func TestCheckTable(t *testing.T) {
	assert.Require(t, checkTable("# header") == false)
	assert.Require(t, checkTable("[github](https://github.com/)") == false)
	assert.Require(t, checkTable("| hogehoge") == false)
	assert.Require(t, checkTable("|head1|head2|") == true)
	assert.Require(t, checkTable("|-----|") == true)
}

func TestGetSeparateTable(t *testing.T) {
	result := getSeparateTable(`
	|head1|head2|
	|-----|-----------|
	|fiiiiiiiit|ho|
	|a|bbb               |
	`)

	want := [][]string{
		[]string{"head1", "head2"},
		[]string{"-----", "-----------"},
		[]string{"fiiiiiiiit", "ho"},
		[]string{"a", "bbb"},
	}

	assert.Require(t, len(result) == 4)

	for lineIdx, line := range result {
		for columnIdx, column := range line {
			assert.Require(t, column == want[lineIdx][columnIdx])
		}
	}
}

func TestGetHasSeparatorLine(t *testing.T) {
	var target [][]string

	target = [][]string{
		[]string{"head1", "head2"},
		[]string{"-----", "-----------"},
		[]string{"fiiiiiiiit", "ho"},
		[]string{"a", "bbb"},
	}
	assert.Require(t, getHasSeparatorLine(target) == true)

	target = [][]string{
		[]string{"head1", "head2"},
		[]string{"fiiiiiiiit", "ho"},
		[]string{"a", "bbb"},
	}
	assert.Require(t, getHasSeparatorLine(target) == false)
}

func TestGetMaxLength(t *testing.T) {
	var target [][]string
	var maxLength []int
	var want []int

	target = [][]string{
		[]string{"head1", "head2"},
		[]string{"----------------------", "------"},
		[]string{"fiiiiiiiit", "ho"},
		[]string{"a", "bbb"},
	}
	maxLength = getMaxLength(target)

	want = []int{10, 5} // fiiiiiiiit => 10, head2 => 5

	assert.Require(t, len(maxLength) == 2)
	for i, length := range maxLength {
		assert.Require(t, length == want[i])
	}
}

func TestGetMaxLength2(t *testing.T) {
	var target [][]string
	var maxLength []int
	var want []int

	target = [][]string{
		[]string{"head1", "head2"},
		[]string{"----------------------", "------"},
		[]string{"こんにちは", "ho"},
		[]string{"a", "bbb"},
	}
	maxLength = getMaxLength(target)

	want = []int{10, 5} // こんにちは => 10, head2 => 5

	assert.Require(t, len(maxLength) == 2)
	for i, length := range maxLength {
		assert.Require(t, length == want[i])
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

	assert.Require(t, result == want)
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

	assert.Require(t, result == want)
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
		[]int{3, 6},
		[]int{8, 9},
	}

	assert.Require(t, len(resultsList) == 2)

	for i, results := range resultsList {
		assert.Require(t, len(results) == 2)

		for j, result := range results {
			assert.Require(t, wantsList[i][j] == result)
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
		[]int{2, 5},
	}

	assert.Require(t, len(resultsList) == 1)

	for i, results := range resultsList {
		assert.Require(t, len(results) == 2)

		for j, result := range results {
			assert.Require(t, wantsList[i][j] == result)
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
		[]int{3, 6},
	}

	assert.Require(t, len(resultsList) == 1)

	for i, results := range resultsList {
		assert.Require(t, len(results) == 2)

		for j, result := range results {
			assert.Require(t, wantsList[i][j] == result)
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

	assert.Require(t, len(result) == 2)

	for i, table := range result {
		resultLines := strings.Split(table, "\n")
		wantLines := strings.Split(strings.TrimSpace(want[i]), "\n")

		for j, resultLint := range resultLines {
			assert.Require(t, strings.TrimSpace(resultLint) == strings.TrimSpace(wantLines[j]))
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

	assert.Require(t, len(resultSplit) == len(wantSplit))

	for idx, resultLine := range resultSplit {
		a := strings.TrimSpace(resultLine)
		b := strings.TrimSpace(wantSplit[idx])
		assert.Require(t, a == b)
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

	assert.Require(t, len(resultSplit) == len(wantSplit))

	for idx, resultLine := range resultSplit {
		a := strings.TrimSpace(resultLine)
		b := strings.TrimSpace(wantSplit[idx])
		assert.Require(t, a == b)
	}
}
