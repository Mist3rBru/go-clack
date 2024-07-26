package prompts

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Mist3rBru/go-clack/core/utils"
	"github.com/Mist3rBru/go-clack/prompts/symbols"
	"github.com/Mist3rBru/go-clack/third_party/picocolors"
)

type TableAlign string

const (
	TableAlignLeft   TableAlign = "left"
	TableAlignRight  TableAlign = "right"
	TableAlignCenter TableAlign = "center"
)

type TableOptions struct {
	Align  []TableAlign
	Output io.Writer
}

func Table(rows [][]string, options TableOptions) {
	if len(options.Align) == 0 {
		options.Align = make([]TableAlign, len(rows[0]))
	}
	if options.Output == nil {
		options.Output = os.Stdout
	}

	var sizes []int

	for _, row := range rows {
		for i, col := range row {
			colLength := len(col)

			if i >= len(sizes) || sizes[i] == 0 || colLength > sizes[i] {
				if i < len(sizes) {
					sizes[i] = colLength
				} else {
					sizes = append(sizes, colLength)
				}
			}
		}
	}

	var table string
	separator := picocolors.Dim(symbols.BAR)

	for _, row := range rows {
		var tableRow []string

		for i, col := range row {
			remainingWidth := sizes[i] - utils.StrLength(col)
			spacing := strings.Repeat(" ", max(remainingWidth, 0))

			var tableCol string

			switch options.Align[i] {
			case TableAlignCenter:
				tableCol = fmt.Sprint(
					strings.Repeat(" ", (remainingWidth+1)/2),
					col,
					strings.Repeat(" ", (remainingWidth)/2),
				)

			case TableAlignRight:
				tableCol = spacing + col

			default:
				tableCol = col + spacing
			}

			tableRow = append(tableRow, tableCol)
		}

		table += fmt.Sprint(separator, " ", strings.Join(tableRow, " "+separator+" "), " ", separator, "\n")
	}

	options.Output.Write([]byte(table))
}
