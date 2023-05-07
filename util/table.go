package util

import (
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func renderTable(table *tablewriter.Table, headers []string) {

	var colColors []tablewriter.Colors
	var headerColors []tablewriter.Colors
	colCount := len(headers)
	headerCount := len(headers)

	table.SetHeader(headers)

	for i := 0; i < colCount; i++ {
		colColors = append(colColors, tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor})
	}
	table.SetColumnColor(colColors...)

	for i := 0; i < headerCount; i++ {
		headerColors = append(headerColors, tablewriter.Colors{tablewriter.Bold, tablewriter.BgMagentaColor})
	}
	table.SetHeaderColor(headerColors...)

	table.SetRowLine(true)
	table.SetCenterSeparator(color.HiBlackString("_"))
	table.SetColumnSeparator(color.HiBlackString("|"))
	table.SetRowSeparator(color.HiBlackString("_"))
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.Render()
}

func MakeColored(str string) string {
	newStr := strings.Split(str, ",")

	if len(newStr) > 1 {
		return makeColoredString(newStr[0])
	} else {
		return makeColoredString(str)
	}
}
func makeColoredString(str string) string {
	if str == "Started" || str == "True" || str == "true" || str == "Running" || str == "Auto" {
		green := color.New(color.BgGreen)
		coloredStr := green.Sprint(str)
		return coloredStr
	} else if str == "Stopped" || str == "false" || str == "False" || strings.Contains(str, "Stopped") || str == "Manual" {
		red := color.New(color.BgRed)
		coloredStr := red.Sprint(str)
		return coloredStr
	} else {
		return str
	}
}

func MakeTableFromStruct(s interface{}) {

	header := createHeaders(s)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	row := createRow(s)
	table.Append(row)
	renderTable(table, header)
}

func MakeTableWithRowHeader(row []string, header []string) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)
	table.Append(row)
	renderTable(table, header)

}

func MakeTable(slice interface{}) {

	sliceInterface := convertStructSliceToInterface(slice)
	h := createHeaders(sliceInterface[0])
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(h)
	for _, s := range sliceInterface {
		row := (createRow(s))
		table.Append(row)
	}

	renderTable(table, h)
}

func createHeaders(s interface{}) []string {
	var fields []string
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		fields = append(fields, t.Field(i).Name)
	}
	return fields
}

func convertStructSliceToInterface(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		log.Fatal("The parameter is not a slice.")
	}
	var interfaceStruct []interface{}

	for i := 0; i < s.Len(); i++ {
		interfaceStruct = append(interfaceStruct, s.Index(i).Interface())
	}
	return interfaceStruct
}

func createRow(s interface{}) []string {
	var row []string
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			field := MakeColored(field.String())
			row = append(row, field)
		case reflect.Int, reflect.Int64:
			row = append(row, strconv.FormatInt(field.Int(), 10))
		case reflect.Bool:
			field := strconv.FormatBool(field.Bool())
			field = MakeColored(field)
			row = append(row, field)
		default:
			//data = append(data, "PSComputerName couldn't resolved")
		}
	}
	return row
}
