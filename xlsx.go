package main

import (
	"fmt"
	"log"
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

func toValue(s, typ string) interface{} {
	switch typ {
	case "string":
		return s
	case "int":
		v, _ := strconv.ParseFloat(s, 64)
		return int(v)
	case "float":
		v, _ := strconv.ParseFloat(s, 64)
		return v
	}
	return nil
}

func (f *Field) Value() interface{} {
	return toValue(f.str, f.Type)
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

func makePk(dataRow []string, typeRow []string, pkCols []int, pkSep string) (pk, escapePk interface{}) {
	for _, v := range pkCols {
		if v < 0 || v >= len(dataRow) {
			panic(fmt.Errorf("主键列无效:%v, %d", pkCols, len(dataRow)))
		}
	}

	if len(pkCols) > 1 {
		var key string
		for i, col := range pkCols {
			colIdx := col - 1
			val := toValue(dataRow[colIdx], typeRow[colIdx])
			key += fmt.Sprint(val)
			if i < len(pkCols)-1 {
				key += pkSep
			}
		}
		return key, escape(key, "string")
	}

	colIdx := pkCols[0] - 1
	s := dataRow[colIdx]
	typ := typeRow[colIdx]
	pk = toValue(s, typ)
	return pk, escape(pk, typ)
}

type XlsxData struct {
	Name string
	Rows []*RowData
}

func isInSlice(list []int, _v int) bool {
	for _, v := range list {
		if v == _v {
			return true
		}
	}
	return false
}

func getExcludeCols(colNum int, cfg *Config, idx int) []int {
	cols := cfg.GetCols(idx)
	if len(cols) == 0 {
		return nil
	}
	exclude := []int{}
	if cols[0] > 0 {
		for i := 1; i <= colNum; i++ {
			if !isInSlice(cols, i) {
				exclude = append(exclude, i)
			}
		}
	} else if cols[0] < 0 {
		for i := 1; i <= colNum; i++ {
			if isInSlice(cols, -i) {
				exclude = append(exclude, i)
			}
		}
	}
	return exclude
}

func readRow(row *xlsx.Row, colNum int) []string {
	empty := true
	ret := []string{}
	for idx := 0; idx < colNum; idx++ {
		var v string
		if idx < len(row.Cells) {
			v, _ = row.Cells[idx].String()
		}
		v = strings.TrimSpace(v)
		if v != "" {
			empty = false
		}
		ret = append(ret, v)
	}

	if empty {
		return nil
	}

	// 填充空白单元格
	for i := len(ret); i < colNum; i++ {
		ret = append(ret, "")
	}
	return ret
}

func checkHeader(rowsData [][]string) error {
	if len(rowsData) < 3 {
		return fmt.Errorf("缺少header行!")
	}

	descRow := rowsData[0]
	nameRow := rowsData[1]
	typRow := rowsData[2]
	if len(descRow) != len(nameRow) || len(nameRow) != len(typRow) {
		return fmt.Errorf("描述行，名称行和类型行的列数应该相同！")
	}

	for i, _ := range descRow {
		if nameRow[i] == "" {
			return fmt.Errorf("第%d列名称为空", i+1)
		}
		if typRow[i] == "" {
			return fmt.Errorf("第%d列类型为空", i+1)
		}
	}

	nameMap := map[string]int{}
	for i, v := range nameRow {
		if nameMap[v] > 0 {
			return fmt.Errorf("%v 第%d,%d列名字冲突!", nameRow, nameMap[v], i+1)
		}
		nameMap[v] = i + 1
	}
	return nil
}

// 读取所有列
func readFull(fp string, cfg *Config, cfIdx int) [][]string {
	file, err := xlsx.OpenFile(fp)
	if err != nil {
		panic(err)
	}

	si := cfg.GetSheet(cfIdx)
	if si <= 0 {
		panic(fmt.Errorf("没有指定数据表sheet:", cfIdx))
	}
	sheet := file.Sheets[si-1]
	if len(sheet.Rows) < 4 {
		panic(fmt.Errorf("无效数据表:", cfIdx))
	}

	colNum := func() int {
		var idx int
		for idx = 0; idx < len(sheet.Rows[1].Cells); idx++ {
			// desc, _ := sheet.Rows[0].Cells[idx].String()
			name, _ := sheet.Rows[1].Cells[idx].String()
			typ, _ := sheet.Rows[2].Cells[idx].String()
			log.Println(idx, name, typ)
			if name == "" || typ == "" {
				break
			}
		}
		return idx
	}()

	log.Printf("%s, sheet %d 列数: %d\n", cfg.List[cfIdx].Input, cfg.List[cfIdx].Sheet, colNum)

	ret := [][]string{}
	for _, row := range sheet.Rows {
		list := readRow(row, colNum)
		if list == nil {
			break
		}
		ret = append(ret, list)
	}

	return ret
}

// fp: file path
func readXlsxData(fp string, cfg *Config, cfIdx int) *XlsxData {
	rowsData := readFull(fp, cfg, cfIdx)
	if err := checkHeader(rowsData); err != nil {
		cfItem := cfg.List[cfIdx]
		panic(fmt.Errorf("input: %s, sheet: %d, error: %v", cfItem.Input, cfItem.Sheet, err))
	}

	name := strings.Split(cfg.List[cfIdx].Input, ".")[0]
	xlsxData := &XlsxData{Name: name}
	colNum := len(rowsData[0])
	pkCols := cfg.GetPkCols(cfIdx)
	excludeCols := getExcludeCols(colNum, cfg, cfIdx)
	log.Println(colNum, pkCols, excludeCols)

	descRow := rowsData[0]
	nameRow := rowsData[1]
	typeRow := rowsData[2]
	for _, rowData := range rowsData[3:] {
		fieldList := []*Field{}
		log.Println("===>", rowData)
		for colIdx, val := range rowData {
			if !isInSlice(excludeCols, colIdx+1) {
				desc := descRow[colIdx]
				name := nameRow[colIdx]
				typ := typeRow[colIdx]
				field := &Field{Desc: desc, Name: name, Type: typ, str: val}
				log.Printf("field:%+v\n", field)
				fieldList = append(fieldList, field)
			}
		}
		row := &RowData{Data: fieldList}
		row.Pk, row.EscapePk = makePk(rowData, typeRow, pkCols, cfg.PkSep)
		xlsxData.Rows = append(xlsxData.Rows, row)
	}
	return xlsxData
}
