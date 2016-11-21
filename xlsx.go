package main

import (
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

type Field struct {
	Desc string
	Name string
	Type string
	str  string
}

func (f *Field) Value() interface{} {
	switch f.Type {
	case "string":
		return f.str
	case "int":
		v, _ := strconv.ParseFloat(f.str, 64)
		return int(v)
	case "float":
		v, _ := strconv.ParseFloat(f.str, 64)
		return v
	}
	return nil
}

func (f *Field) EscapeValue() interface{} {
	if f.Type == "string" {
		return fmt.Sprintf("\"%s\"", f.str)
	}
	return f.Value()
}

type RowData struct {
	Pk       interface{}
	EscapePk interface{}
	Data     []*Field
}

func escape(pk interface{}, typ string) interface{} {
	if typ == "string" {
		return fmt.Sprintf("\"%v\"", pk)
	}
	return pk
}

type XlsxData struct {
	Name string
	Rows []*RowData
}

func makePk(rowData []*Field, pkCols []int, pkSep string) (pk, escapePk interface{}) {
	for _, v := range pkCols {
		if v < 0 || v >= len(rowData) {
			panic(fmt.Errorf("主键列无效:%v", pkCols))
		}
	}

	if len(pkCols) > 1 {
		var key string
		for i, v := range pkCols {
			key += fmt.Sprint(rowData[v-1].Value())
			if i < len(pkCols)-1 {
				key += pkSep
			}
		}
		return key, escape(key, "string")
	}

	pkField := rowData[pkCols[0]-1]
	return pkField.Value(), escape(pkField.Value(), pkField.Type)
}

func readRow(row *xlsx.Row) []string {
	ret := []string{}
	for _, cell := range row.Cells {
		s, err := cell.String()
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, strings.TrimSpace(s))
	}
	return ret
}

func readHeader(sheet *xlsx.Sheet, ecols []int) (descRow, nameRow, typRow []string, err error) {
	descRow = readRow(sheet.Rows[0])
	nameRow = readRow(sheet.Rows[1])
	typRow = readRow(sheet.Rows[2])
	if len(descRow) != len(nameRow) || len(nameRow) != len(typRow) {
		return nil, nil, nil, fmt.Errorf("描述行，名称行和类型行的列数应该相同！")
	}

	for i, _ := range descRow {
		if nameRow[i] == "" {
			return nil, nil, nil, fmt.Errorf("第%d列名称为空", i+1)
		}
		if typRow[i] == "" {
			return nil, nil, nil, fmt.Errorf("第%d列类型为空", i+1)
		}
	}

	var retDescRow, retNameRow, retTypRow []string
	for i, _ := range descRow {
		if !isInSlice(ecols, i+1) {
			retDescRow = append(retDescRow, descRow[i])
			retNameRow = append(retNameRow, nameRow[i])
			retTypRow = append(retTypRow, typRow[i])
		}
	}
	return retDescRow, retNameRow, retTypRow, nil
}

func isInSlice(list []int, _v int) bool {
	for _, v := range list {
		if v == _v {
			return true
		}
	}
	return false
}

func getExcludeCols(cfg *Config, sheet *xlsx.Sheet, idx int) []int {
	cols := cfg.GetCols(idx)
	if len(cols) == 0 {
		return nil
	}
	exclude := []int{}
	if cols[0] > 0 {
		for i := 1; i <= len(sheet.Cols); i++ {
			if !isInSlice(cols, i) {
				exclude = append(exclude, i)
			}
		}
	} else if cols[0] < 0 {
		for i := 1; i < len(sheet.Cols); i++ {
			if isInSlice(cols, -i) {
				exclude = append(exclude, i)
			}
		}
	}
	return exclude
}

// fp: file path
func readXlsxData(fp string, cfg *Config, cfIdx int) *XlsxData {
	file, err := xlsx.OpenFile(fp)
	if err != nil {
		log.Fatal(err)
	}

	ret := &XlsxData{}
	_, ret.Name = path.Split(fp)
	si := cfg.GetSheet(cfIdx)
	if si <= 0 {
		panic(fmt.Errorf("没有指定数据表sheet"))
	}
	sheet := file.Sheets[si-1]
	colNum := len(sheet.Rows[0].Cells)
	ecols := getExcludeCols(cfg, sheet, cfIdx)
	// log.Println(ecols)
	pkCols := cfg.GetPkCols(cfIdx)
	if len(pkCols) == 0 {
		panic("没有指定主键列")
	}
	// log.Println(cfIdx, pkCols)
	descRow, nameRow, typRow, err := readHeader(sheet, ecols)
	if err != nil {
		panic(err)
	}

	for _, row := range sheet.Rows[3:] {
		var rowData []*Field
		isEmpty := true
		for colIdx := 0; colIdx < colNum; colIdx++ {
			if !isInSlice(ecols, colIdx+1) {
				var val string
				if colIdx < len(row.Cells) {
					cell := row.Cells[colIdx]
					val, _ = cell.String()
				}
				field := &Field{str: val}
				rowData = append(rowData, field)
				if val != "" {
					isEmpty = false
				}
			}
		}

		if isEmpty {
			break
		} else {
			for i, field := range rowData {
				field.Desc = descRow[i]
				field.Name = nameRow[i]
				field.Type = typRow[i]
				// log.Println(*field)
			}
			rd := &RowData{Data: rowData}
			rd.Pk, rd.EscapePk = makePk(rowData, pkCols, cfg.PkSep)
			ret.Rows = append(ret.Rows, rd)
		}
	}
	return ret
}
