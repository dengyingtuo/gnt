package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync"
	"text/template"
	"time"
)

var debug bool
var configPath string
var inputPath string
var outputPath string

func init() {
	flag.BoolVar(&debug, "debug", false, "open debug output")
	flag.StringVar(&configPath, "config", "config.yml", "config file in YAML format")
	flag.StringVar(&inputPath, "input", "./xlsx", "input xlsx path")
	flag.StringVar(&outputPath, "output", "./output", "output file path")
}

func isFileExists(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

func main() {
	flag.Parse()
	if configPath == "" || inputPath == "" || outputPath == "" {
		flag.Usage()
		return
	}

	log.SetPrefix("gnt")
	log.SetFlags(log.Lshortfile | log.Ltime | log.Lmicroseconds)
	if !debug {
		log.SetOutput(ioutil.Discard)
	}

	bt := time.Now()
	if !isFileExists(configPath) {
		panic(fmt.Errorf("指定文件:%s不存在!", configPath))
	}
	if !isFileExists(inputPath) {
		panic(fmt.Errorf("指定文件:%s不存在!", inputPath))
	}
	if !isFileExists(outputPath) {
		panic(fmt.Errorf("指定文件:%s不存在!", outputPath))
	}

	cfg := readConfig(configPath)
	fp := path.Join(path.Dir(configPath), cfg.Template)
	// log.Println(fp)
	if !isFileExists(fp) {
		panic(fmt.Errorf("指定文件:%s不存在!", fp))
	}

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"quote": func(v interface{}) string {
			return fmt.Sprintf("\"%v\"")
		},
	}

	jobs := sync.WaitGroup{}
	for i, v := range cfg.List {
		xlsxData := readXlsxData(path.Join(inputPath, v.Input), cfg, i)
		tpl, err := template.New(cfg.Template).Funcs(funcMap).ParseFiles(path.Join(path.Dir(configPath), cfg.Template))
		if err != nil {
			panic(err)
		}
		fp := path.Join(outputPath, v.Output)
		outputFile, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			panic(err)
		}

		jobs.Add(1)
		go func() {
			err = tpl.Execute(outputFile, xlsxData)
			if err != nil {
				panic(err)
			}
			jobs.Done()
		}()
	}
	jobs.Wait()
	log.Print("处理完毕，耗时:", time.Now().Sub(bt))
}
