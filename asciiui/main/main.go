package main

import (
	"fmt"

	"github.com/graaphscom/monogo/asciiui"
)

func main() {
	//window, err := asciiui.NewWindow()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	table := asciiui.Table{
		Rows: []asciiui.TableRow{
			{Cells: []asciiui.TableCell{{Text: "Approved"}, {Text: "Description"}}},
			{Cells: []asciiui.TableCell{{Text: "✅"}, {Text: "Thing approved because it's great"}}},
			{Cells: []asciiui.TableCell{{Text: "⛔️"}, {Text: "Thing rejected because it requires improvement"}}},
		},
	}

	fmt.Print(table.Render())
}
