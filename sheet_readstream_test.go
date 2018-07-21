package xlsx_test

import (
	"github.com/plandem/xlsx"
	"github.com/plandem/xlsx/options"
	"github.com/plandem/xlsx/types"
	"github.com/stretchr/testify/require"
	"testing"
)

//load content using multi phases or normal mode - must be same for both
func testSheetReadFull(t *testing.T, sheet xlsx.Sheet) {
	cols, rows := sheet.Dimension()
	require.Equal(t, 14, cols)
	require.Equal(t, 28, rows)

	emptyRow := []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""}
	values := [][]string{
		emptyRow,
		{"    with leading space", "", "merged rows", "", "merged cols", "merged cols", "merged cols", "", "", "", "", "", "", ""},
		{"", "", "merged rows", "", "", "", "", "", "", "", "", "", "", ""},
		{"", "", "merged rows", "", "merged rows+cols", "merged rows+cols", "merged rows+cols", "", "", "", "", "", "", ""},
		{"with trailing space   ", "", "merged rows", "", "merged rows+cols", "merged rows+cols", "merged rows+cols", "", "", "", "", "", "", ""},
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		{"", "", "", "1", "2", "3", "4", "5", "", "", "", "", "", ""},
		{"", "", "", "6", "7", "8", "9", "10", "", "", "", "", "", ""},
		{"", "", "", "11", "12", "13", "14", "15", "", "", "", "", "", ""},
		{"", "", "", "16", "17", "18", "19", "20", "", "", "", "", "", ""},
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "last cell"},
	}

	for rowIdx, rv := range values {
		require.Equal(t, rv, sheet.Row(rowIdx).Values())
	}
}

//load content using single phase only - it will be without merged cells info - only first cell from range will be output
func testSheetReadLimited(t *testing.T, sheet xlsx.Sheet) {
	cols, rows := sheet.Dimension()
	require.Equal(t, 14, cols)
	require.Equal(t, 28, rows)

	emptyRow := []string{"", "", "", "", "", "", "", "", "", "", "", "", "", ""}
	values := [][]string{
		emptyRow,
		{"    with leading space", "", "merged rows", "", "merged cols", "", "", "", "", "", "", "", "", ""},
		emptyRow,
		{"", "", "", "", "merged rows+cols", "", "", "", "", "", "", "", "", ""},
		{"with trailing space   ", "", "", "", "", "", "", "", "", "", "", "", "", ""},
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		{"", "", "", "1", "2", "3", "4", "5", "", "", "", "", "", ""},
		{"", "", "", "6", "7", "8", "9", "10", "", "", "", "", "", ""},
		{"", "", "", "11", "12", "13", "14", "15", "", "", "", "", "", ""},
		{"", "", "", "16", "17", "18", "19", "20", "", "", "", "", "", ""},
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		emptyRow,
		{"", "", "", "", "", "", "", "", "", "", "", "", "", "last cell"},
	}

	for rowIdx, rv := range values {
		require.Equal(t, rv, sheet.Row(rowIdx).Values())
	}
}

func TestSpreadsheet(t *testing.T) {
	xl, err := xlsx.Open("./test_files/example_simple.xlsx")
	if err != nil {
		panic(err)
	}

	defer xl.Close()

	type beforeTestFn func(t *testing.T, xl *xlsx.Spreadsheet) xlsx.Sheet
	type testFn func(t *testing.T, sheet xlsx.Sheet)

	sheetTests := []struct {
		name     string
		open     beforeTestFn
		callback testFn
	}{
		{
			"SheetReadStream_SinglePhased", func(t *testing.T, xl *xlsx.Spreadsheet) xlsx.Sheet { return xl.SheetReader(0, false) }, testSheetReadLimited,
		},
		{
			"SheetReadStream_MultiPhased", func(t *testing.T, xl *xlsx.Spreadsheet) xlsx.Sheet { return xl.SheetReader(0, true) }, testSheetReadFull,
		},
		{
			"SheetReadWrite", func(t *testing.T, xl *xlsx.Spreadsheet) xlsx.Sheet { return xl.Sheet(0) }, testSheetReadFull,
		},
	}

	for _, info := range sheetTests {
		t.Run(info.name, func(tt *testing.T) {
			sheet := info.open(tt, xl)
			info.callback(tt, sheet)
		})
	}
}

func TestSheetReadStream_notImplemented(t *testing.T) {
	xl, err := xlsx.Open("./test_files/example_simple.xlsx")
	if err != nil {
		panic(err)
	}
	defer xl.Close()

	sheet := xl.SheetReader(0, true)
	defer sheet.Close()

	require.Panics(t, func() { sheet.Col(0) })
	require.Panics(t, func() { sheet.Cols() })
	require.Panics(t, func() { sheet.InsertCol(0) })
	require.Panics(t, func() { sheet.InsertRow(0) })
	require.Panics(t, func() { sheet.DeleteRow(0) })
	require.Panics(t, func() { sheet.DeleteCol(0) })
	require.Panics(t, func() { sheet.SetDimension(100, 100) })
	require.Panics(t, func() { sheet.SetActive() })
	require.Panics(t, func() { sheet.Set(options.NewSheetOptions(options.Sheet.Visibility(types.VisibilityTypeVisible))) })
	require.Panics(t, func() { sheet.SetName("aaa") })
}

func TestSheetReadStream_access(t *testing.T) {
	xl, err := xlsx.Open("./test_files/example_simple.xlsx")
	if err != nil {
		panic(err)
	}

	defer xl.Close()

	sheet := xl.SheetReader(0, true)
	defer sheet.Close()

	require.Equal(t, "8", sheet.CellByRef("F11").Value())
	require.Equal(t, "", sheet.CellByRef("F10").Value())
	require.Equal(t, "8", sheet.Cell(5, 10).Value())
	require.Equal(t, []string{"", "", "", "", "", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}, sheet.Range("D10:H13").Values())
}