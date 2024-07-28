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

	colWidths := make([]int, len(rows[0]))

	for _, row := range rows {
		for i, col := range row {
			colLength := len(col)

			if i >= len(colWidths) {
				colWidths = append(colWidths, colLength)
			} else if colLength > colWidths[i] {
				colWidths[i] = colLength
			}
		}
	}

	table := ""
	colSeparator := picocolors.Dim(symbols.BAR)

	for i, row := range rows {
		var tableRow []string
		var tableRowSeparator []string

		for j, col := range row {
			colWith := colWidths[j]
			remainingColWidth := colWith - utils.StrLength(col)
			spacing := strings.Repeat(" ", max(remainingColWidth, 0))

			var tableCol string

			switch options.Align[j] {
			case TableAlignCenter:
				tableCol = fmt.Sprint(
					strings.Repeat(" ", (remainingColWidth+1)/2),
					col,
					strings.Repeat(" ", (remainingColWidth)/2),
				)

			case TableAlignRight:
				tableCol = spacing + col

			default:
				tableCol = col + spacing
			}

			tableRow = append(tableRow, " "+tableCol+" ")
			tableRowSeparator = append(tableRowSeparator, strings.Repeat(symbols.BAR_H, colWith+2))
		}

		if i == 0 {
			table += picocolors.Dim(fmt.Sprint(symbols.CONNECT_TOP_LEFT, strings.Join(tableRowSeparator, symbols.CONNECT_TOP), symbols.CONNECT_TOP_RIGHT, "\n"))
		}

		table += fmt.Sprint(colSeparator, strings.Join(tableRow, colSeparator), colSeparator, "\n")

		if i+1 < len(rows) {
			table += picocolors.Dim(fmt.Sprint(symbols.CONNECT_LEFT, strings.Join(tableRowSeparator, symbols.CONNECT_CENTER), symbols.CONNECT_RIGHT, "\n"))
		} else {
			table += picocolors.Dim(fmt.Sprint(symbols.CONNECT_BOTTOM_LEFT, strings.Join(tableRowSeparator, symbols.CONNECT_BOTTOM), symbols.CONNECT_BOTTOM_RIGHT, "\n"))
		}
	}

	options.Output.Write([]byte(table))
}
