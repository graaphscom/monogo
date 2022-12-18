package asciiui

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	borderEdge           = "+"
	horizontalBorderFill = "-"
	cellPadding          = " "
	verticalBorderFill   = "|"
)

func (table Table) Render() (string, error) {
	if len(table.Rows) < 1 {
		return "", nil
	}

	if err := table.validate(); err != nil {
		return "", err
	}

	columnsWidths := table.findMaxTextLengthByColumns()
	cellPaddingWidth := len(cellPadding)
	horizontalBorder := borderEdge

	for _, columnWidth := range columnsWidths {
		horizontalBorder += strings.Repeat(horizontalBorderFill, columnWidth+(cellPaddingWidth*2)) + borderEdge
	}
	horizontalBorder += "\n"

	result := horizontalBorder

	for _, row := range table.Rows {
		rowResult := verticalBorderFill
		for idx, cell := range row.Cells {
			rowResult += cellPadding +
				cell.Text +
				strings.Repeat(" ", columnsWidths[idx]-utf8.RuneCountInString(cell.Text)) +
				cellPadding +
				verticalBorderFill
		}
		result += rowResult + "\n" + horizontalBorder
	}

	return result, nil
}

type Table struct {
	Rows []TableRow
}

type TableRow struct {
	Cells []TableCell
}

type TableCell struct {
	Text string
}

func (table Table) validate() error {
	if len(table.Rows) < 1 {
		return nil
	}

	firstRowColumnsCount := len(table.Rows[0].Cells)
	for idx, row := range table.Rows {
		if len(row.Cells) != firstRowColumnsCount {
			return fmt.Errorf("each row of the table has to have the same columns count (bad row index: %d)", idx)
		}
	}

	return nil
}

func (table Table) findMaxTextLengthByColumns() []int {
	if len(table.Rows) < 1 {
		return []int{}
	}
	result := make([]int, len(table.Rows[0].Cells))

	for _, row := range table.Rows {
		for idx, cell := range row.Cells {
			if result[idx] < utf8.RuneCountInString(cell.Text) {
				result[idx] = utf8.RuneCountInString(cell.Text)
			}
		}
	}

	return result
}
