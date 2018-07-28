package markdowntableformatter

import (
	"flag"
	"os"
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

	assert.Require(t, len(result) == 6)

	want := []string{
		"# header",
		"",
		"|head1|head2|",
		"|-----|-----------|",
		"|fiiiiiiiit|ho|",
		"|a|bbb               |",
	}
	for i, r := range result {
		assert.Require(t, r == want[i])
	}
}

func TestCheckTable(t *testing.T) {
	assert.Require(t, checkTable("# header") == false)
	assert.Require(t, checkTable("[github](https://github.com/)") == false)
	assert.Require(t, checkTable("| hogehoge") == false)
	assert.Require(t, checkTable("|head1|head2|") == true)
	assert.Require(t, checkTable("|-----|") == true)
	assert.Require(t, checkTable("head1|head2") == true)
}

func TestFixColumnSize(t *testing.T) {
	target := ""
	target += "|head1|head2|" + "\n"
	target += "|-----|-----------|" + "\n"
	target += "|fiiiiiiiit|ho|" + "\n"
	target += "|a|bbb               |" + "\n"
	result := fixColumnSize(target)

	want := ""
	want += "|head1     |head2|" + "\n"
	want += "|----------|-----|" + "\n"
	want += "|fiiiiiiiit|ho   |" + "\n"
	want += "|a         |bbb  |" + "\n"

	assert.Require(t, result == want)
}
