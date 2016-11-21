package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Item struct {
	Input  string
	Sheet  int
	PkCols []int
	Cols   []int
	Output string
}

type Config struct {
	Template string
	Sheet    int
	PkCols   []int
	PkSep    string
	Ext      string
	List     []Item
}

// idx: 配置项索引
func (cfg *Config) GetPkCols(idx int) []int {
	if len(cfg.List[idx].PkCols) > 0 {
		return cfg.List[idx].PkCols
	}
	return cfg.PkCols
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
	cols := cfg.List[idx].Cols
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
