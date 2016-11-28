package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"gopkg.in/yaml.v2"
)

type Item struct {
	Input    string
	Sheet    int
	PkCols   []string
	Cols     []string
	ColsConv map[string]string
	Output   string
}

func (item *Item) GetConv(letterCol string) (cmd string, args []string) {
	if item.ColsConv == nil {
		return
	}

	var conv string
	if conv = item.ColsConv[letterCol]; conv == "" {
		return
	}

	list := strings.Split(conv, " ")
	if len(list) == 0 {
		return
	}

	cmd = list[0]
	args = list[1:]
	return
}

type Config struct {
	Template string
	Sheet    int
	PkCols   []string
	PkSep    string
	Ext      string
	List     []Item
}

func toLetterColumn(icol int) string {
	var sign byte
	if icol < 0 {
		sign = '-'
		icol = int(math.Abs(float64(icol)))
	}
	icol -= 1
	ret := ""
	for {
		n := icol / 26
		m := icol % 26
		ret = string('A'+byte(m)) + ret
		if n == 0 {
			break
		}
		icol = n - 1
	}

	if sign == '-' {
		ret = fmt.Sprintf("%c%s", sign, ret)
	}
	return ret
}

// 字母表示的列转为列序号, 例如A:1
func toIntColumn(name string) (int, error) {
	name = strings.ToUpper(name)
	lcol := name
	if lcol[0] == '-' {
		lcol = lcol[1:]
	}
	ret := 0
	for i := range lcol {
		v := lcol[len(lcol)-i-1]
		if v < 'A' || v > 'Z' {
			return 0, fmt.Errorf("无效列名:%s", name)
		}
		ret += int((byte(v) - 'A' + 1)) * int(math.Pow(26, float64(i)))
	}

	if name[0] == '-' {
		ret = -ret
	}
	return ret, nil
}

func toIntColumns(names []string) ([]int, error) {
	ret := []int{}
	for _, name := range names {
		v, err := toIntColumn(name)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}
	return ret, nil
}

// idx: 配置项索引
func (cfg *Config) GetPkCols(idx int) []int {
	if len(cfg.List[idx].PkCols) > 0 {
		cols, err := toIntColumns(cfg.List[idx].PkCols)
		if err != nil {
			panic(fmt.Errorf("主键列无效:%s %s", cfg.List[idx].Input, err))
		}
		return cols
	}
	cols, err := toIntColumns(cfg.PkCols)
	if err != nil {
		panic(fmt.Errorf("主键列无效: %s", err))
	}
	return cols
}

// idx: 配置项索引
func (cfg *Config) GetSheet(idx int) int {
	if cfg.List[idx].Sheet > 0 {
		return cfg.List[idx].Sheet
	}
	return cfg.Sheet
}

// idx: 配置项索引
// nil表示所有列(不可能处理没有数据的表)
func (cfg *Config) GetCols(idx int) []int {
	cols, err := toIntColumns(cfg.List[idx].Cols)
	if err != nil {
		panic(fmt.Errorf("包含列无效: %s, %s", cfg.List[idx].Input, err))
	}

	if len(cols) == 0 {
		return nil
	}

	// 校验
	if cols[0] >= 0 {
		// 包含
		for _, ci := range cols {
			if ci <= 0 {
				panic(fmt.Errorf("包含冲突: %s %v %d", cfg.Template, cols, ci))
			}
		}
	} else {
		// 不包含
		for _, ci := range cols {
			if ci >= 0 {
				panic(fmt.Errorf("包含冲突: %s %v %d", cfg.Template, cols, ci))
			}
		}
	}
	return cols
}

func readConfig(fp string) *Config {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Fatal(err)
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
