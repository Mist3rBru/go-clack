package prompts_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Mist3rBru/go-clack/prompts"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	rows := [][]string{
		{"abcde", "123"},
		{"abc", "12345"},
	}
	expected := strings.Join([]string{
		"│ abcde │ 123   │",
		"│ abc   │ 12345 │",
		"",
	}, "\n")

	var writer MockWriter
	prompts.Table(rows, prompts.TableOptions{Output: &writer})

	assert.Equal(t, expected, writer.Data[0])
}

type TableCase struct {
	Rows   [][]string
	Align  []prompts.TableAlign
	Output string
}

func TestTableAlignment(t *testing.T) {
	cases := []TableCase{
		{
			Rows: [][]string{
				{"abc", "12"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignRight, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ abc │ 12 │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"abc", "12"},
				{"def", "3456"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignLeft, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ abc │   12 │",
				"│ def │ 3456 │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"abc", "1234"},
				{"def", "56"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignLeft, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ abc │ 1234 │",
				"│ def │   56 │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"abc", "1234"},
				{"def", "56"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignRight, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ abc │ 1234 │",
				"│ def │   56 │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"beep", "1234"},
				{"blip", "33450"},
				{"abc", "1006"},
				{"def", "45"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignLeft, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ beep │  1234 │",
				"│ blip │ 33450 │",
				"│ abc  │  1006 │",
				"│ def  │    45 │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"beep", "1234"},
				{"blip", "33450"},
				{"abc", "1006"},
				{"def", "45"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignLeft, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ beep │  1234 │",
				"│ blip │ 33450 │",
				"│ abc  │  1006 │",
				"│ def  │    45 │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"beep", "1234"},
				{"blip", "33450"},
				{"abc", "1006"},
				{"def", "45"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignLeft, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ beep │  1234 │",
				"│ blip │ 33450 │",
				"│ abc  │  1006 │",
				"│ def  │    45 │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"abc", "1234"},
				{"def", "56"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignCenter, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ abc │ 1234 │",
				"│ def │   56 │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"abc", "1234"},
				{"def", "56"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignRight, prompts.TableAlignLeft},
			Output: strings.Join([]string{
				"│ abc │ 1234 │",
				"│ def │ 56   │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"abc", "1234"},
				{"def", "56"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignRight, prompts.TableAlignLeft},
			Output: strings.Join([]string{
				"│ abc │ 1234 │",
				"│ def │ 56   │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"beep", "1234", "abc"},
				{"blip", "33450", "abc"},
				{"abc", "1006", "abcdef"},
				{"def", "45", "abcd"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignLeft, prompts.TableAlignCenter, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ beep │  1234 │    abc │",
				"│ blip │ 33450 │    abc │",
				"│ abc  │  1006 │ abcdef │",
				"│ def  │   45  │   abcd │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"beep", "1234", "abc"},
				{"blip", "1234567", "abc"},
				{"abc", "12345", "abcdef"},
				{"def", "12", "abcd"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignCenter, prompts.TableAlignCenter, prompts.TableAlignCenter},
			Output: strings.Join([]string{
				"│ beep │   1234  │   abc  │",
				"│ blip │ 1234567 │   abc  │",
				"│  abc │  12345  │ abcdef │",
				"│  def │    12   │  abcd  │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"beep", "1234", "abc"},
				{"blip", "1234567", "abc"},
				{"abc", "12345", "abcdef"},
				{"def", "12", "abcd"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignCenter, prompts.TableAlignRight, prompts.TableAlignCenter},
			Output: strings.Join([]string{
				"│ beep │    1234 │   abc  │",
				"│ blip │ 1234567 │   abc  │",
				"│  abc │   12345 │ abcdef │",
				"│  def │      12 │  abcd  │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"beep", "1234", "abc"},
				{"blip", "1234567", "abc"},
				{"abc", "12345", "abcdef"},
				{"def", "12", "abcd"},
			},
			Align: []prompts.TableAlign{prompts.TableAlignLeft, prompts.TableAlignCenter, prompts.TableAlignLeft},
			Output: strings.Join([]string{
				"│ beep │   1234  │ abc    │",
				"│ blip │ 1234567 │ abc    │",
				"│ abc  │  12345  │ abcdef │",
				"│ def  │    12   │ abcd   │",
				"",
			}, "\n"),
		},
		{
			Rows: [][]string{
				{"beep", "1234", "abc"},
				{"blip", "1234567", "abc"},
				{"abc", "12345", "abcdef"},
				{"def", "12", "abcd"},
			},
			Align: []prompts.TableAlign{"right", prompts.TableAlignRight, prompts.TableAlignRight},
			Output: strings.Join([]string{
				"│ beep │    1234 │    abc │",
				"│ blip │ 1234567 │    abc │",
				"│  abc │   12345 │ abcdef │",
				"│  def │      12 │   abcd │",
				"",
			}, "\n"),
		},
	}
	var writer MockWriter

	for i, _case := range cases {
		prompts.Table(_case.Rows, prompts.TableOptions{
			Align:  _case.Align,
			Output: &writer,
		})
		assert.Equal(t, _case.Output, writer.Data[i], fmt.Sprintf("CaseIndex: %d", i))
	}
}

func TestTableWithAnsiColors(t *testing.T) {
	rows := [][]string{
		{picocolors.Bold("abcde"), picocolors.Bold("123")},
		{picocolors.Red("abc"), picocolors.Green("12345")},
	}
	expected := strings.Join([]string{
		"│ abcde │ 123   │",
		"│ abc   │ 12345 │",
		"",
	}, "\n")

	var writer MockWriter
	prompts.Table(rows, prompts.TableOptions{Output: &writer})

	assert.Equal(t, expected, writer.Data[0])

}
