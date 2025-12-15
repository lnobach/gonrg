package outtable

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/lnobach/gonrg/obis"
	"github.com/lnobach/gonrg/version"
	"github.com/rodaine/table"
)

func PrintTable(res *obis.OBISListResult) {

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	raised := color.New(color.FgHiYellow, color.Bold)

	fmt.Println()
	fmt.Print("\u26A1\uFE0F") // lightning symbol
	raised.Print(version.GonrgName)
	fmt.Printf(" version %s\n", version.GonrgVersion)
	fmt.Print("Device ID: ")
	raised.Println(res.DeviceID)

	fmt.Println()

	tbl := table.New("Exact Key", "Simple Key", "Name", "Value", "Unit")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, part := range res.List {
		tbl.AddRow(part.ExactKey, part.SimplifiedKey, part.Name, part.PrettyValue(false), part.Unit)
	}

	tbl.Print()
}
