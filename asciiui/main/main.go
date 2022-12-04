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
			{[]asciiui.TableCell{{"Approved"}, {"Description"}}},
			{[]asciiui.TableCell{{"✅"}, {"Thing approved because it's great"}}},
			{[]asciiui.TableCell{{"⛔️"}, {"Thing rejected because it requires improvement"}}},
		},
	}

	fmt.Print(table.Render())
}
